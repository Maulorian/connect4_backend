package database

import (
	"connect4_backend/game"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Node struct {
	StateID     int
	Simulations int
	Wins        float32
	Children    []int
	Move        game.Coordinate
}

func (node Node) WinRate() float32 {
	return node.Wins / float32(node.Simulations)
}

func GetNode(stateId int, client *mongo.Client) *Node {
	fmt.Println("Getting node from database with state_id=", stateId)
	nodes := client.Database("connect4").Collection("nodes")
	filter := bson.D{{"state_id", stateId}}
	var result Node
	err := nodes.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return &result
}
