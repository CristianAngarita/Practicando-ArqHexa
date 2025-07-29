package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Player struct {
	Name         string    `json:"name" binding:"required"`
	Age          int       `json:"age" binding:"required"`
	CreationTime time.Time `json:"-"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ginEngine := gin.Default()
	ginEngine.POST("/players", func(ctx *gin.Context) {
		var player Player
		if err := ctx.BindJSON(&player); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		player.CreationTime = time.Now().UTC()

		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(c, options.Client.ApplyURI(os.Getenv("MONGO_URI")))
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(c, nil)
		if err != nil {
			log.Fatal(err)
		}

		collection := client.Database("go-l").Collection("players")
		insertResult, err := collection.InsertOne(c, player)
		if err != nil {
			log.Fatal(err)
		}

		ctx.JSON(200, gin.H{"player_id": insertResult.InsertedID})
	})

	log.Fatalln(ginEngine.Run(":8001"))
}
