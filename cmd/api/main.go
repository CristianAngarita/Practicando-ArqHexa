package cmd

import (
	"context"
	"database/sql"
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

	ginEngine := gin.Default()

	ginEngine.POST("/players", func(ctx *gin.Context) {
		var player Player
		if err := ctx.BindJSON(&player); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		player.CreationTime = time.Now().UTC()

		// 2. Ejecutar la consulta de inserción en MySQL
		query := `INSERT INTO players (name, age, creation_time) VALUES (?, ?, ?)`
		result, err := db.ExecContext(ctx.Request.Context(), query, player.Name, player.Age, player.CreationTime)
		if err != nil {
			log.Printf("Error al insertar jugador en MySQL: %v", err)
			ctx.JSON(500, gin.H{"error": "No se guardó jugador"})
			return
		}

		// Obtener el ID generado automáticamente por MySQL
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error al obtener el ID generado: %v", err)
			//  éxito parcial si la inserción fue exitosa
			ctx.JSON(200, gin.H{"message": "Player created, pero no se pudo recuperar la ID"})
			return
		}

		ctx.JSON(200, gin.H{"player_id": id})
	})

	log.Fatalln(ginEngine.Run(":8001"))

}
