package player

import (
	"database/sql"
	"log"
	"proyecto-gin-hexagonal/internal/core"
	"time"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	DB *sql.DB // The database connection
}

func NewPlayerHandler(db *sql.DB) *PlayerHandler {
	return &PlayerHandler{DB: db}
}

func (ph *PlayerHandler) CreatePlayer(ctx *gin.Context) {

	var player core.Player
	if err := ctx.BindJSON(&player); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	player.CreationTime = time.Now().UTC()

	// 2. Ejecutar la consulta de inserción en MySQL
	query := `INSERT INTO players (name, age, creation_time) VALUES (?, ?, ?)`
	result, err := ph.DB.ExecContext(ctx.Request.Context(), query, player.Name, player.Age, player.CreationTime)
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

}
