package cmd

import (
	"context"
	"database/sql"
	"log"
	"os"
	"proyecto-gin-hexagonal/cmd/api/handlers/player"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Obtener la URI de MySQL desde las variables de entorno
	mysqlURI := os.Getenv("MYSQL_URI")
	if mysqlURI == "" {
		log.Fatal("MYSQL_URI environment variable not set")
	}

	// 1. Conexión a MySQL
	db, err := sql.Open("mysql", mysqlURI)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a MySQL: %v", err)
	}
	defer db.Close()

	// Verificar la conexión
	ctxDB, cancelDB := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDB()
	if err = db.PingContext(ctxDB); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos MySQL: %v", err)
	}
	log.Println("Conexión a MySQL establecida con éxito!")

	playerHandler := player.NewPlayerHandler(db)
	ginEngine := gin.Default()

	ginEngine.POST("/players", func(ctx *gin.Context) {
		playerHandler.CreatePlayer(ctx)
	})
	log.Fatalln(ginEngine.Run(":8001"))

}
