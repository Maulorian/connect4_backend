package db

import (
	"connect4_backend/game"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

type Node struct {
	StateID     int
	Simulations int
	Wins        float32
	Children    []int
	Move        game.Coordinate
}

func (node Node) ChildWithBestWinRate() Node {
	fmt.Println("ChildWithBestWinRate()")
	var childWithBestWinRate Node
	var maxValue = float32(math.SmallestNonzeroFloat32)

	for _, childID := range node.Children {
		var dbChild = GetNode(childID)
		if dbChild == nil {
			continue
		}
		var winRate = dbChild.WinRate()
		if winRate > maxValue {
			maxValue = winRate
			childWithBestWinRate = *dbChild
		}
	}

	return childWithBestWinRate
}
func (node Node) WinRate() float32 {
	return node.Wins / float32(node.Simulations)
}

func GetNode(id int) *Node {
	fmt.Println("GetNode()")
	client := GetMongoClient()
	nodes := client.Database("connect4").Collection("nodes")
	filter := bson.D{{"state_id", id}}
	var result Node
	err := nodes.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return &result
}
