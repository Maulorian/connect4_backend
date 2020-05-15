package main

import (
	"connect4_backend/ai"
	"connect4_backend/game"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"time"
)

//func hello(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Welcome home!")
//}
func getMove(c *gin.Context) {
	fmt.Println("getMove()")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var s game.State
	err = json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(s)
	var node = ai.NewNode(s, nil)

	move := ai.GetBestMove(node)
	fmt.Println(move)
	c.JSON(200, gin.H{
		"move": move,
	})
}
func CORSMiddleware(c *gin.Context) {
	fmt.Println("Sending cors headers")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}

	c.Next()
}
func main() {
	route := gin.Default()
	route.Use(CORSMiddleware)
	route.POST("/getmove", getMove)
	route.POST("/savegame", saveGame)
	_ = route.Run()
}

func saveGame(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var g Game
	err = json.Unmarshal(body, &g)
	if err != nil {
		fmt.Println(err)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://adp:bCi1M4NPkFgEfRzX@yeda-lan6r.gcp.mongodb.net"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	games := client.Database("connect4").Collection("games")
	var toInsert = make(map[string]interface{})
	toInsert["nb_moves"] = len(g.Moves)
	toInsert["played_at"] = time.Now()
	toInsert["ip"] = c.ClientIP()
	toInsert["moves"] = g.Moves
	toInsert["outcome"] = g.Outcome
	toInsert["started_at"] = time.Unix(g.StartedAt, 0)
	toInsert["ended_at"] = time.Unix(g.EndedAt, 0)
	toInsert["duration"] = g.Duration
	id, _ := games.InsertOne(ctx, toInsert)
	fmt.Println("added game:", id)
}

type Game struct {
	Moves     []game.Coordinate `bson:"moves" json:"moves"`
	Outcome   int               `bson:"outcome" json:"outcome"`
	StartedAt int64             `bson:"started_at" json:"started_at"`
	EndedAt   int64             `bson:"ended_at" json:"ended_at"`
	Duration  int               `bson:"duration" json:"duration"`
}
