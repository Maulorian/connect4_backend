package main

import (
	"connect4_backend/ai"
	"connect4_backend/db"
	"connect4_backend/game"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"runtime"
	"sync"
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
	analysisWG := sync.WaitGroup{}
	updatesWG := sync.WaitGroup{}
	analysisWG.Add(1)
	updatesWG.Add(1)
	var node = ai.NewNode(s, nil)
	//var mutex = sync.Mutex{}
	//go Process(node, &analysisWG, &updatesWG, &mutex)
	go func() {
		fmt.Println("Analysis")

		//mutex.Lock()
		ai.Analysis(node)
		fmt.Println("Analysis done")
		analysisWG.Done()
		nodes := node.FlattenChildren()
		//mutex.Unlock()

		ai.UpdateNodes(nodes)

		fmt.Println("updates from", node.State.GetID(), "done")
		updatesWG.Done()
	}()

	fmt.Println("waiting for analysis to finish")
	analysisWG.Wait()
	fmt.Println("analysis finished: node has", len(node.Children), "children")

	var stateId = node.State.GetID()
	fmt.Println("stateid:", stateId)
	var dbNode = db.GetNode(stateId)
	fmt.Println("dbNode:", dbNode)

	if dbNode == nil || len(dbNode.Children) == 0 {
		fmt.Println("dbNode is either nil or has no children")
		bestChild := node.ChildWithBestWinRate()
		c.JSON(200, gin.H{
			"move": bestChild.State.Move,
		})
		fmt.Println("getMove() done")

		return
	}
	fmt.Println("waiting for updates to finish")
	updatesWG.Wait()
	fmt.Println("updates finished")

	node.UpdateChildrenFromDatabase()
	bestChild := node.ChildWithBestWinRate()
	var move = bestChild.State.Move
	fmt.Println(bestChild)
	fmt.Println(move)

	c.JSON(200, gin.H{
		"move": move,
	})
	//analysisWG.Wait()
	fmt.Println("getMove() done")

}
func Process(node *ai.Node, analysisWG *sync.WaitGroup, updatesWG *sync.WaitGroup, mutex *sync.Mutex) {
	fmt.Println("Analysis")
	analysisWG.Add(1)
	updatesWG.Add(1)
	mutex.Lock()
	ai.Analysis(node)
	fmt.Println("Analysis done")
	analysisWG.Done()
	nodes := node.FlattenChildren()
	mutex.Unlock()

	ai.UpdateNodes(nodes)

	fmt.Println("updates from", node.State.GetID(), "done")
	updatesWG.Done()
}

func CORSMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		fmt.Println("Sending cors headers")
		c.AbortWithStatus(200)
		return
	}

	c.Next()
}
func main() {
	fmt.Println(runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

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
	client := db.GetMongoClient()
	games := client.Database("connect4").Collection("games")
	var toInsert = make(map[string]interface{})
	toInsert["nb_moves"] = len(g.Moves)
	toInsert["ip"] = c.ClientIP()
	toInsert["moves"] = g.Moves
	toInsert["outcome"] = g.Outcome
	toInsert["started_at"] = time.Unix(g.StartedAt, 0)
	toInsert["ended_at"] = time.Unix(g.EndedAt, 0)
	toInsert["duration"] = g.Duration
	id, _ := games.InsertOne(context.TODO(), toInsert)
	fmt.Println("added game:", id)
}

type Game struct {
	Moves     []game.Coordinate `bson:"moves" json:"moves"`
	Outcome   int               `bson:"outcome" json:"outcome"`
	StartedAt int64             `bson:"started_at" json:"started_at"`
	EndedAt   int64             `bson:"ended_at" json:"ended_at"`
	Duration  int               `bson:"duration" json:"duration"`
}
