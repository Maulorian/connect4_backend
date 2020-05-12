package ai

import (
	"connect4_backend/game"
	"fmt"
	"math/rand"
	"time"
)

const CalculationTime = 3 //seconds
const ExplorationParameter = 1.414213562373095

func GetBestMove(node *Node) game.Coordinate {
	rand.Seed(time.Now().UnixNano())

	var spent int64
	var delta int64
	var lastLoopCallTime = time.Now()
	var dic = make(map[int]int)
	for resourcesLeft(spent) {
		//fmt.Println(node.simulations)
		var leaf = Traverse(node)

		outcome := leaf.state.Playout()
		dic[outcome] += 1

		//fmt.Println(outcome)
		Backpropagate(leaf, outcome)

		delta = time.Now().Sub(lastLoopCallTime).Nanoseconds()
		spent += delta
		lastLoopCallTime = time.Now()
		//fmt.Println(spent)
	}
	fmt.Println(dic)

	fmt.Println("total simulation :", node.simulations)
	for _, child := range node.children {
		fmt.Println(child.GetUCT(), child.winRate, child.state.Move, child.wins, child.simulations)
	}

	return node.ChildWithBestWinRate().state.Move
}

func Backpropagate(node *Node, o int) {
	node.simulations++

	if o == game.Draw {
		node.wins += 0.5
	} else if game.GetPlayer(o) == node.state.PreviousPlayer {
		node.wins++
	}
	node.winRate = node.wins / float32(node.simulations)
	if node.isRoot(){return}
	Backpropagate(node.parent, o)

}
func Traverse(node *Node) *Node {
	//fmt.Println("Traverse()")
	var root = node
	//fmt.Println(root)
	for !root.isLeaf() {
		var bestChoice = root.ChildWithBestUTC()
		//if bestChoice == nil{
		//	break
		//}
		root = bestChoice
		if len(root.children) == 0 {
			break
		}
	}

	child, err := root.GenerateChildren()

	if err != nil {
		//fmt.Println(err.Error())
		return root
	}
	return child
}

func resourcesLeft(spent int64) bool {
	return spent < (CalculationTime * time.Second).Nanoseconds()
}
