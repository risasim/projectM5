package state

import (
	"database/sql"
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
	file, err := c.FormFile("sound")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file", "details": err.Error()})
		return
	}
	_, err = e.db.Exec("UPDATE users SET deathSound=$1 WHERE username=$2", file.Filename, username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save filename to database", "details": err.Error()})
	}
	err = c.SaveUploadedFile(file, "src/server/soundEffects"+file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Successfully updated death sound"})
}

func (e EndPointHandler) GetGameStatus(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "Game_Status": e.GameManager.GameStatus.String()})
}

type StartGameRequest struct {
	GameType communication.GameType `json:"game_type"`
}

func (e EndPointHandler) StartGame(c *gin.Context) {
	var req StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	if e.GameManager.GameStatus == Active {
		c.JSON(400, gin.H{"error": "A game is already active"})
		return
	}
	err := e.GameManager.StartNewGame(req.GameType)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start game", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "New game started"})
}

func (e EndPointHandler) StopGame(c *gin.Context) {
	if e.GameManager.GameStatus == idle {
		c.JSON(400, gin.H{"error": "No game is active"})
		return
	}
	err := e.GameManager.EndGame()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to end game", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Game Stopped"})
}

func (e EndPointHandler) DeleteUser(c *gin.Context) {
	username := c.Query("username")
	_, err := e.db.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user from Database", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "User deleted"})
}

func (e EndPointHandler) JoinGame(c *gin.Context) {
	username := c.Query("username")
	repo := db.NewUsersRepository(e.db)
	user, err := repo.GetUser(username)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found in database", "details": err.Error()})
	}
	player := Player{
		Username:   username,
		PiSN:       user.PiSN,
		DeathSound: user.DeathSound,
	}
	err = e.GameManager.AddPlayer(player)
	if err != nil {
		c.JSON(500, gin.H{"error": "couldn't add player", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Successfully joined game"})
}
