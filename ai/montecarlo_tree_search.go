package ai

import (
	"connect4_backend/game"
	"fmt"
	"math/rand"
	"time"
)

const CalculationTime = 1 //seconds
const ExplorationParameter = 10

func Analysis(node *Node) {
	//fmt.Println("Analysis()")
	rand.Seed(time.Now().UnixNano())
	set := make(map[int]bool)
	var spent int64
	var delta int64
	var lastLoopCallTime = time.Now()
	for resourcesLeft(spent) {
		var leaf = Traverse(node, set)
		outcome := leaf.State.Playout()
		Backpropagate(leaf, outcome)

		delta = time.Now().Sub(lastLoopCallTime).Nanoseconds()
		spent += delta
		lastLoopCallTime = time.Now()
		//break
	}
	fmt.Println("total simulation :", node.Simulations)
	for _, child := range node.Children {
		fmt.Println("wins:", child.Wins, "simulations:", child.Simulations, "winRate:", child.WinRate(), "utc:", child.GetUCT(), "move:", child.State.Move, "done by:", child.State.PreviousPlayer)
	}
}

func Backpropagate(node *Node, o int) {
	//fmt.Println("Backpropagate()")
	//fmt.Println(node)
	node.Simulations++

	if o == game.Draw {
		node.Wins += 0.5
	} else if o == node.State.PreviousPlayer {
		node.Wins++
	}
	if node.isRoot() {
		return
	}
	Backpropagate(node.parent, o)

}
func Traverse(node *Node, set map[int]bool) *Node {
	//fmt.Println("Traverse()")
	var root = node
	//fmt.Println(root)
	for !root.isLeaf() {
		var bestChoice = root.ChildWithBestUTC()
		if bestChoice == nil {
			break
		}
		root = bestChoice
	}
	child, err := root.GenerateChildren(set)

	if err != nil {
		//fmt.Println(err.Error())
		return root
	}
	return child
}

func resourcesLeft(spent int64) bool {
	return spent < (CalculationTime * time.Second).Nanoseconds()
}
