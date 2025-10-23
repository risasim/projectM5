package state

import (
	"database/sql"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/risasim/projectM5/project/src/server/communication"
	"github.com/risasim/projectM5/project/src/server/db"
)

type EndPointHandlerInterface interface {
	uploadSound(g *gin.Context)
}

type EndPointHandler struct {
	db          *sql.DB
	GameManager *GameManager
}

func NewEndPointHandler(db *sql.DB) *EndPointHandler {
	return &EndPointHandler{db: db}
}

func (e EndPointHandler) UploadSound(c *gin.Context) {
	username := c.Query("username")
	repo := db.NewUsersRepository(e.db)
	user, err := repo.GetUser(username)
	if err != nil || user == nil {
		c.JSON(404, gin.H{"error": "User not found in database"})
		return
	}
	if e.GameManager.GameStatus == Active {
		c.JSON(400, gin.H{"error": "Cannot change sound during a game"})
		return
	}
	file, err := c.FormFile("sound")
	if err != nil {
		c.JSON(400, gin.H{"error": "Missing sound file"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot open file"})
		return
	}
	defer src.Close()
	fileBytes := make([]byte, file.Size)
	_, err = src.Read(fileBytes)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot read file"})
		return
	}
	b64Sound := base64.StdEncoding.EncodeToString(fileBytes)
	_, err = e.db.Exec("UPDATE users SET deathSound=$1 WHERE username=$2", b64Sound, username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database deathSound update failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully updated death sound"})
}

func (e EndPointHandler) GetGameStatus(c *gin.Context) {
	c.JSON(200, gin.H{"status": e.GameManager.GameStatus.String()})
}

func (e EndPointHandler) StartGame(c *gin.Context) {
	gameTypeString := c.Query("GameType")
	gameType, _ := communication.ParseGameType(gameTypeString)
	if e.GameManager.GameStatus == Active {
		c.JSON(400, gin.H{"error": "A game is already active"})
		return
	}
	e.GameManager.StartNewGame(gameType)
	c.JSON(200, gin.H{"message": "New game started"})
}

func (e EndPointHandler) StopGame(c *gin.Context) {
	if e.GameManager.GameStatus == idle {
		c.JSON(400, gin.H{"error": "No game is active"})
		return
	}
	e.GameManager.EndGame()
	c.JSON(200, gin.H{"message": "Game Stopped"})
}

func (e EndPointHandler) DeleteUser(c *gin.Context) {
	username := c.Query("username")
	_, err := e.db.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user from Database"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}

func (e EndPointHandler) JoinGame(c *gin.Context) {
	username := c.Query("username")
	repo := db.NewUsersRepository(e.db)
	user, err := repo.GetUser(username)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found in database"})
	}
	player := Player{
		Username:   username,
		ID:         int(user.ID),
		EncodingID: 0,
		DeathSound: user.DeathSound,
	}
	err = e.GameManager.AddPlayer(player)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully joined game"})
}
