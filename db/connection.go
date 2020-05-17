package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetMongoClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://adp:bCi1M4NPkFgEfRzX@yeda-lan6r.gcp.mongodb.net")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
