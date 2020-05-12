package main

import (
	"connect4_backend/ai"
	"connect4_backend/game"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
	fmt.Println()
	//c.Header("Access-Control-Allow-Origin", "*")
	////c.Header("Access-Control-Allow-Origin", "https://secure-island-74494.herokuapp.com/")
	//c.Header("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	//c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	c.JSON(200, gin.H{
		"move": move,
	})
}
func main() {
	route := gin.Default()
	route.POST("/getmove", getMove)
	route.OPTIONS("/sendcors", sendCors)
	_ = route.Run()

	//MONGO CODE
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//err = client.Connect(ctx)
	//facebook_ads := client.Database("FacebookAds").Collection("facebook_ads")
	//cur, err := facebook_ads.Find(ctx, bson.D{{"ad_id", 254742718983521}})
	//if err != nil { log.Fatal(err) }
	//defer cur.Close(ctx)
	//for cur.Next(ctx) {
	//	var result bson.M
	//	err := cur.Decode(&result)
	//	if err != nil { log.Fatal(err) }
	//	// do something with result....
	//	fmt.Print(result["app_store_ids"])
	//}
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
}

func sendCors(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	//c.Header("Access-Control-Allow-Origin", "https://secure-island-74494.herokuapp.com/")
	context.Header("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	context.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	context.JSON(200, gin.H{})
}
