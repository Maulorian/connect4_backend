package ai

import (
	"connect4_backend/db"
	"connect4_backend/game"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math"
	"math/rand"
)

type Node struct {
	Simulations int
	Wins        float32
	State       game.State
	Children    []*Node
	parent      *Node
}

func (node Node) WinRate() float32 {
	return node.Wins / float32(node.Simulations)
}
func NewNode(state game.State, parent *Node) *Node {
	return &Node{State: state, parent: parent}
}

func (node Node) isLeaf() bool {
	return node.Simulations == 0
	//return len(node.children) == 0
}

func (node Node) isRoot() bool {
	return node.parent == nil
}

func (node Node) isTerminal() bool {
	return node.State.Outcome != game.None
}

func (node *Node) GenerateChildren(set map[int]bool) (*Node, error) {
	if node.isTerminal() {
		return nil, errors.New("node is terminal")
	}
	var rows = node.State.GetFreeRows()
	//fmt.Println(rows)
	for col := 0; col < game.Cols; col++ {
		row := rows[col]

		var deepCopy = node.State
		if row == -1 {
			continue
		}
		deepCopy.PlayMove(game.Coordinate{
			Col: col,
			Row: row,
		})
		stateId := deepCopy.GetID()
		if set[stateId] {
			continue
		}
		set[stateId] = true
		node.Children = append(node.Children, NewNode(deepCopy, node))
	}
	if len(node.Children) == 0 {
		return nil, errors.New("node has no children")
	}
	//fmt.Println(node.children)
	return node.GetRandomChild(), nil

}

func (node Node) GetStatsFromDb() {
	var dbNode = db.GetNode(node.State.GetID())
	node.Simulations = dbNode.Simulations
	node.Wins = dbNode.Wins
}

func (node Node) ChildWithBestUTC() *Node {
	//fmt.Println("ChildWithBestUTC()")
	var bestChoice *Node
	var maxValue = math.SmallestNonzeroFloat64
	// debugger

	//forsl
	for _, child := range node.Children {
		var uct = child.GetUCT()
		if uct > maxValue {
			maxValue = uct
			bestChoice = child
		}
	}
	return bestChoice
}

func (node Node) GetUCT() float64 {
	//fmt.Println("GetUCT")

	if node.parent == nil {
		fmt.Println("parent is nil")
	}
	if node.Simulations == 0 {
		return math.MaxInt64
	}
	var term1 = node.Wins / float32(node.Simulations)
	var log = math.Log(float64(node.parent.Simulations))
	var division = log / float64(node.Simulations)
	var power = math.Pow(division, 0.5)

	var term2 = power

	// return term1 + (EXPLORATION_PARAMETER * term2)
	return float64(term1) + ExplorationParameter*term2
}
func (node *Node) UpdateChildrenFromDatabase() {
	fmt.Println("UpdateChildrenFromDatabase()")
	for _, child := range node.Children {
		stateId := child.State.GetID()
		dbChild := db.GetNode(stateId)
		if dbChild == nil {
			fmt.Println(stateId, "had no entry in db")
			continue
		}
		child.Simulations = dbChild.Simulations
		child.Wins = dbChild.Wins
	}
}
func (node Node) ChildWithBestWinRate() Node {
	fmt.Println("ChildWithBestWinRate()")
	var childWithBestWinRate Node
	var maxValue = float32(math.SmallestNonzeroFloat32)

	for _, child := range node.Children {
		var winRate = child.WinRate()
		if winRate > maxValue {
			maxValue = winRate
			childWithBestWinRate = *child
		}
	}

	return childWithBestWinRate
}

func (node *Node) GetRandomChild() *Node {

	if len(node.Children) == 1 {
		return node.Children[0]
	}
	var randomIndex = rand.Intn(len(node.Children) - 1)
	return node.Children[randomIndex]
}

func UpdateNodes(nodes []Node) {
	fmt.Println("UpdateNodes():", len(nodes))
	var operations []mongo.WriteModel
	for _, node := range nodes {
		var dbNode = node.ConvertToDatabase()
		//fmt.Println("dbNode: ", dbNode)
		upsert := mongo.NewUpdateOneModel()
		update := bson.M{
			"$inc":      bson.D{{"simulations", dbNode.Simulations}, {"wins", dbNode.Wins}},
			"$set":      bson.D{{"move", dbNode.Move}},
			"$addToSet": bson.D{{"children", bson.D{{"$each", dbNode.Children}}}},
		}

		upsert.SetFilter(bson.M{"state_id": dbNode.StateID})
		upsert.SetUpdate(update)
		upsert.SetUpsert(true)
		operations = append(operations, upsert)
	}
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(false)
	client := db.GetMongoClient()
	client.Connect(context.TODO())
	collection := client.Database("connect4").Collection("nodes")
	result, err := collection.BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UpsertedCount:", result.UpsertedCount, "nodes")
	fmt.Println("ModifiedCount:", result.ModifiedCount, "nodes")
}

func (node Node) ConvertToDatabase() db.Node {
	var children = make([]int, 0)
	for _, child := range node.Children {
		children = append(children, child.State.GetID())
	}
	var dbNode = db.Node{
		StateID:     node.State.GetID(),
		Simulations: node.Simulations,
		Wins:        node.Wins,
		Move:        node.State.Move,
		Children:    children,
	}
	return dbNode
}

func (node Node) FlattenChildren() []Node {
	var list []Node
	for _, child := range node.Children {
		list = append(list, *child)
	}
	list = append(list, node)
	return list
}

//func flattenChildren(node Node, list *[]Node, set map[int]bool) {
//	stateId := node.State.GetID()
//	if !set[stateId] {
//		*list = append(*list, node)
//		set[stateId] = true
//	}
//	for _, child := range node.Children {
//		flattenChildren(*child, list, set)
//	}
//}
