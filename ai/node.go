package ai

import (
	"connect4_backend/game"
	"errors"
	"fmt"
	"math"
	"math/rand"
)

type Node struct {
	winRate     float32
	simulations int
	wins        float32
	state       game.State
	children    []*Node
	parent      *Node
}

func NewNode(state game.State, parent *Node) *Node {
	return &Node{state: state, parent: parent}
}



func (node Node) isLeaf() bool {
	//fmt.Println("isLeaf():", node)

	return node.simulations == 0
	//return len(node.children) == 0

}

func (node Node) isRoot() bool {
	return node.parent == nil
}

func (node Node) isTerminal() bool {
	return node.state.Outcome != game.None
}

func (node *Node) GenerateChildren() (*Node, error) {
	if node.isTerminal() {
		return nil, errors.New("node is terminal")
	}
	var rows = node.state.GetFreeRows()
	//fmt.Println(rows)
	for col := 0; col < game.Cols; col++ {
		row := rows[col]

		var deepCopy = node.state
		if row == -1 {
			continue
		}
		deepCopy.PlayMove(game.Coordinate{
			Col: col,
			Row: row,
		})

		node.children = append(node.children, NewNode(deepCopy, node))
	}
	//fmt.Println(node.children)
	return node.GetRandomChild(), nil
}

func (node Node) ChildWithBestWinRate() *Node {

	var childWithBestWinRate *Node
	//var maxValue = float32(math.SmallestNonzeroFloat32)
	var maxValue = math.MinInt8
	// debugger
	for _, child := range node.children {
		var winRate = child.simulations
		if winRate > maxValue {
			maxValue = winRate
			childWithBestWinRate = child
		}
	}

	return childWithBestWinRate
}

func (node Node) ChildWithBestUTC() *Node {
	//fmt.Println("ChildWithBestUTC()")
	var bestChoice *Node
	var maxValue = math.SmallestNonzeroFloat64
	// debugger

	//forsl
	for _, child := range node.children {
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

	if node.parent == nil{
		fmt.Println("parent is nil")
	}
	if node.simulations == 0 {
		return math.MaxInt64
	}
	var term1 = node.wins / float32(node.simulations)
	var log = math.Log(float64(node.parent.simulations))
	var division = log/float64(node.simulations)
	var power = math.Pow(division, 0.5)

	var term2 = power

	// return term1 + (EXPLORATION_PARAMETER * term2)
	return float64(term1) + ExplorationParameter * term2
}

func (node Node) GetRandomChild() *Node {

	if len(node.children) == 1 {
		return node.children[0]
	}
	var randomIndex = rand.Intn(len(node.children) - 1)
	return node.children[randomIndex]
}


