package ai

import (
	"connect4_backend/game"
	"fmt"
	"testing"
)

func TestNewNode(t *testing.T) {
	var state = game.State{
		Grid:           [7][6]int{},
		Outcome:        0,
		CurrentPlayer:  0,
		PreviousPlayer: 0,
		Move:           game.Coordinate{},
		NbMoves:        0,
	}
	var node1 = NewNode(state, nil)
	state.PlayMove(game.Coordinate{
		Col: 0,
		Row: 0,
	})
	state.PlayMove(game.Coordinate{
		Col: 1,
		Row: 0,
	})
	fmt.Println(node1.State)

}
func TestNode_Expand(t *testing.T) {
	var state = game.State{
		Grid:           [7][6]int{},
		Outcome:        0,
		CurrentPlayer:  0,
		PreviousPlayer: 0,
		Move:           game.Coordinate{},
		NbMoves:        0,
	}
	var node = NewNode(state, nil)
	node.GenerateChildren(nil)
	for _, child := range node.Children {
		if child.State.NbMoves == 0 {
			t.Error("children should have 1 move")
		}
	}
	if len(node.Children) != game.Cols {
		t.Error("should have 7 children")
	}
}
func TestNode_GetUCT(t *testing.T) {
	var state = game.State{
		Grid:           [7][6]int{},
		Outcome:        0,
		CurrentPlayer:  0,
		PreviousPlayer: 0,
		Move:           game.Coordinate{},
		NbMoves:        0,
	}
	var parent = NewNode(state, nil)
	parent.Simulations = 672313
	var node = NewNode(state, parent)
	node.Wins = 187591.5
	node.Simulations = 279740
	fmt.Println(node.GetUCT())
}
func TestCopy(t *testing.T) {
	var state = game.NewState()
	var node = NewNode(*state, nil)
	var deepCopy = node.State
	deepCopy.PlayMove(game.Coordinate{
		Col: 0,
		Row: 5,
	})
	state.PlayMove(game.Coordinate{
		Col: 1,
		Row: 5,
	})
	node.State.PlayMove(game.Coordinate{
		Col: 2,
		Row: 5,
	})
	fmt.Println(deepCopy)
	fmt.Println(state)
	fmt.Println(node.State)
}
func TestNode_WinRate(t *testing.T) {
	var state = game.NewState()
	var node = NewNode(*state, nil)
	node.Simulations = 2
	node.Wins = 0.5
	fmt.Println(node.WinRate())
	if node.WinRate() != 0.25 {
		t.Error("winrate wrong")
	}
}
func TestNode_FlattenChildren(t *testing.T) {
	var state = game.NewState()
	var node = NewNode(*state, nil)
	node.GenerateChildren(nil)
	nodes := node.FlattenChildren()
	fmt.Println(nodes)
}
