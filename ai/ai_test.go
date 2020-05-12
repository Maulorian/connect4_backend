package ai

import (
	"connect4_backend/game"
	"fmt"
	"testing"
)

func TestTraverse(t *testing.T) {
	var s = game.NewState()
	var node = NewNode(*s, nil)
	Traverse(node)
}

func TestGetBestMove(t *testing.T) {

	var s = game.NewState()
	fmt.Println(s)

	var node = NewNode(*s, nil)

	move := GetBestMove(node)
	fmt.Println(move)
}
func TestChildWithBestUTC(t *testing.T) {
	var s = game.NewState()
	var node = NewNode(*s, nil)
	Traverse(node)
	child := node.ChildWithBestUTC()
	fmt.Println(child)
}

func BenchmarkNode_GenerateChildren(b *testing.B) {
	var s = game.NewState()
	s.Playout()
	var node = NewNode(*s, nil)

	for i := 0; i < b.N; i++ {
		_, err := node.GenerateChildren()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
func BenchmarkGetBestMove(b *testing.B) {
	var s = game.NewState()
	var node = NewNode(*s, nil)
	for i := 0; i < b.N; i++ {

		move := GetBestMove(node)
		fmt.Println(move)
	}
}
