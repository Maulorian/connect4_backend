package main

import (
	"connect4_backend/ai"
	"connect4_backend/db"
	"connect4_backend/game"
	"fmt"
	"sync"
	"testing"
)

func TestProcess(t *testing.T) {
	waitGroup := sync.WaitGroup{}
	var s = game.NewState()
	fmt.Println(s)

	var node = ai.NewNode(*s, nil)
	go Process(node, &waitGroup)
	var dbNode = db.GetNode(node.State.GetID())
	bestChild := dbNode.ChildWithBestWinRate(nil)

	var move = bestChild.Move
	fmt.Println(bestChild)
	fmt.Println(move)
	waitGroup.Wait()

}
