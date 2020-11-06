package ai

import (
	"connect4_backend/database"
	"connect4_backend/game"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (this Node) WinRate() float32 {
	return this.Wins / float32(this.Simulations)
}
func NewNode(state game.State, parent *Node) *Node {
	return &Node{State: state, parent: parent}
}

func (this Node) isLeaf() bool {
	return this.Simulations == 0
	//return len(this.children) == 0
}

func (this Node) isRoot() bool {
	return this.parent == nil
}

func (this Node) isTerminal() bool {
	return this.State.Outcome != game.None
}

func (this *Node) GenerateChildren(set map[int]bool) (*Node, error) {
	if this.isTerminal() {
		return nil, errors.New("this is terminal")
	}
	var rows = this.State.GetFreeColumns()
	//fmt.Println(rows)
	for col := 0; col < game.Cols; col++ {
		row := rows[col]

		var deepCopy = this.State
		if row == -1 {
			continue
		}
		move := game.Coordinate{
			Col: col,
			Row: row,
		}
		deepCopy.PlayMove(move)

		node := NewNode(deepCopy, this)
		left := deepCopy.GetLeftStateID()
		right := deepCopy.GetRightStateID()
		if set[left] || set[right] {

			continue
		}
		set[left] = true
		set[right] = true
		//fmt.Println("Adding:")
		//node.PrettyPrint()
		this.Children = append(this.Children, node)
	}
	if len(this.Children) == 0 {
		return nil, errors.New("this has no children")
	}
	//fmt.Println(this.children)
	return this.GetRandomChild(), nil

}

func (this Node) GetStatsFromDb() {
	var dbNode = database.GetNode(this.State.GetLeftStateID(), database.Connection())
	this.Simulations = dbNode.Simulations
	this.Wins = dbNode.Wins
}

func (this Node) ChildWithBestUTC() *Node {
	//fmt.Println("ChildWithBestUTC()")
	var bestChoice *Node
	var maxValue = math.SmallestNonzeroFloat64
	// debugger

	//forsl
	for _, child := range this.Children {
		var uct = child.GetUCT()
		if uct > maxValue {
			maxValue = uct
			bestChoice = child
		}
	}
	return bestChoice
}

func (this Node) GetUCT() float64 {
	//fmt.Println("GetUCT")

	if this.parent == nil {
		fmt.Println("parent is nil")
	}
	if this.Simulations == 0 {
		return math.MaxInt64
	}
	var term1 = this.Wins / float32(this.Simulations)
	var log = math.Log(float64(this.parent.Simulations))
	var division = log / float64(this.Simulations)
	var power = math.Pow(division, 0.5)

	var term2 = power

	// return term1 + (EXPLORATION_PARAMETER * term2)
	return float64(term1) + ExplorationParameter*term2
}
func (this *Node) UpdateChildrenFromDatabase() {
	fmt.Println("UpdateChildrenFromDatabase()")
	client := database.Connection()
	for _, child := range this.Children {
		stateId := child.State.GetLeftStateID()
		dbChild := database.GetNode(stateId, client)
		if dbChild == nil {
			fmt.Println(stateId, "had no entry in database")
			continue
		}
		child.Simulations = dbChild.Simulations
		child.Wins = dbChild.Wins
		child.PrettyPrint()
	}
}
func (this Node) ChildWithBestWinRate() Node {
	fmt.Println("ChildWithBestWinRate()")
	var childWithBestWinRate Node
	var maxValue = float32(math.SmallestNonzeroFloat32)

	for _, child := range this.Children {
		var winRate = child.WinRate()
		if winRate > maxValue {
			maxValue = winRate
			childWithBestWinRate = *child
		}
	}

	return childWithBestWinRate
}

func (this *Node) GetRandomChild() *Node {

	if len(this.Children) == 1 {
		return this.Children[0]
	}
	var randomIndex = rand.Intn(len(this.Children) - 1)
	return this.Children[randomIndex]
}

func UpdateNodes(nodes []Node) {
	fmt.Println("Updating ", len(nodes), "nodes in database")
	var operations []mongo.WriteModel
	for _, node := range nodes {
		var dbNode = node.ConvertToDatabase()
		//fmt.Println("dbNode: ", dbNode)
		upsert := mongo.NewUpdateOneModel()
		update := bson.M{
			//"$inc":      bson.D{{"simulations", dbNode.Simulations}, {"wins", dbNode.Wins}},
			"$set":      bson.D{{"move", dbNode.Move}, {"simulations", dbNode.Simulations}, {"wins", dbNode.Wins}},
			"$addToSet": bson.D{{"children", bson.D{{"$each", dbNode.Children}}}},
		}

		upsert.SetFilter(bson.M{"state_id": dbNode.StateID})
		upsert.SetUpdate(update)
		upsert.SetUpsert(true)
		operations = append(operations, upsert)
	}
	client := database.Connection()
	collection := client.Database("connect4").Collection("nodes")
	result, err := collection.BulkWrite(context.TODO(), operations)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UpsertedCount:", result.UpsertedCount, "nodes")
	fmt.Println("ModifiedCount:", result.ModifiedCount, "nodes")
}

func (this Node) ConvertToDatabase() database.Node {
	var children = make([]int, 0)
	for _, child := range this.Children {
		children = append(children, child.State.GetLeftStateID())
	}
	var dbNode = database.Node{
		StateID:     this.State.GetLeftStateID(),
		Simulations: this.Simulations,
		Wins:        this.Wins,
		Move:        this.State.Move,
		Children:    children,
	}
	return dbNode
}

func (this Node) FlattenChildren() []Node {
	var list []Node
	for _, child := range this.Children {
		list = append(list, *child)
	}
	list = append(list, this)
	return list
}

func (this Node) PrettyPrint() {
	fmt.Println("state_id", this.State.GetLeftStateID(), "wins:", this.Wins, "simulations:", this.Simulations, "winRate:", this.WinRate(), "utc:", this.GetUCT(), "move:", this.State.Move, "done by:", this.State.PreviousPlayer)

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
