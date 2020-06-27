package ai

import (
	"connect4_backend/game"
	"fmt"
	"testing"
)

func TestTraverse(t *testing.T) {
	var s = game.NewState()
	var node = NewNode(*s, nil)
	Traverse(node, nil)
}
func TestNode_GenerateChildren(t *testing.T) {
	var s = game.NewState()
	//
	//s.PlayMove(game.Coordinate{
	//	Col: 3,
	//	Row: 5,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 0,
	//	Row: 5,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 2,
	//	Row: 5,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 4,
	//	Row: 5,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 2,
	//	Row: 4,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 2,
	//	Row: 3,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 3,
	//	Row: 4,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 4,
	//	Row: 4,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 0,
	//	Row: 4,
	//})
	//s.PlayMove(game.Coordinate{
	//	Col: 4,
	//	Row: 3,
	//})
	var node = NewNode(*s, nil)
	set := make(map[int]bool)
	node.GenerateChildren(set)

	for _, c := range node.Children {
		c.PrettyPrint()
	}
}
func TestChildWithBestUTC(t *testing.T) {
	var s = game.NewState()
	var node = NewNode(*s, nil)
	Traverse(node, nil)
	child := node.ChildWithBestUTC()
	fmt.Println(child)
}

func BenchmarkNode_GenerateChildren(b *testing.B) {
	var s = game.NewState()
	s.Playout()
	var node = NewNode(*s, nil)

	for i := 0; i < b.N; i++ {
		_, err := node.GenerateChildren(nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
