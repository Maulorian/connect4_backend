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
	fmt.Println(node1.state)

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
	node.GenerateChildren()
	for _, child := range node.children {
		if child.state.NbMoves == 0 {
			t.Error("children should have 1 move")
		}
	}
	if len(node.children) != game.Cols {
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
	parent.simulations = 1200
	var node = NewNode(state, parent)
	node.wins = 120
	node.simulations = 244
	fmt.Println(node.GetUCT())
}