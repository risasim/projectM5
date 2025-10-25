package state

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/risasim/projectM5/project/src/server/communication"
	"github.com/risasim/projectM5/project/src/server/db"
)

type EndPointHandlerInterface interface {
	uploadSound(g *gin.Context)
}

type EndPointHandler struct {
	Repo        db.UserRepositoryInterface
	GameManager *GameManager
}

func NewEndPointHandler(repo db.UserRepositoryInterface, gameManager *GameManager) *EndPointHandler {
	return &EndPointHandler{Repo: repo, GameManager: gameManager}
}

func (e EndPointHandler) UploadSound(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(400, gin.H{"error": "no username in context"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "username in context is not a string"})
		return
	}

	user, err := e.Repo.GetUser(usernameStr)
	if err != nil || user == nil {
		c.JSON(404, gin.H{"error": "User not found in database"})
		return
	}
	file, err := c.FormFile("sound")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file", "details": err.Error()})
		return
	}

	if filepath.Ext(file.Filename) != ".mp3" {
		c.JSON(400, gin.H{"error": "only mp3 files are allowed"})
		return
	}

	saveDir := "soundEffects"
	os.MkdirAll(saveDir, os.ModePerm)

	savePath := filepath.Join(saveDir, file.Filename)
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file", "details": err.Error()})
		return
	}

	err = e.Repo.UpdateDeathSound(usernameStr, file.Filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save filename to database", "details": err.Error()})
	}

	c.JSON(200, gin.H{"status": "success", "message": "Successfully updated death sound"})
}

// GetSound does retrieve the sound based on the data from db and then tries path depending if it is test or deployment
func (e EndPointHandler) GetSound(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(400, gin.H{"error": "no username in context"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "username in context is not a string"})
		return
	}

	user, err := e.Repo.GetUser(usernameStr)
	if err != nil || user == nil {
		c.JSON(404, gin.H{"error": "User or sound not found"})
		return
	}

	saveDir := "soundEffects"
	filePath := filepath.Join(saveDir, user.DeathSound)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		alternativePaths := []string{
			filepath.Join("..", "soundEffects", user.DeathSound),
			filepath.Join("..", "..", "soundEffects", user.DeathSound),
		}

		found := false
		for _, altPath := range alternativePaths {
			if _, err := os.Stat(altPath); err == nil {
				filePath = altPath
				found = true
				break
			}
		}

		if !found {
			c.JSON(404, gin.H{"error": "Sound file not found"})
			return
		}
	}

	c.File(filePath)
}

func (e EndPointHandler) GetGameStatus(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "Game_Status": e.GameManager.GameStatus.String()})
}

type StartGameRequest struct {
	GameType communication.GameType `json:"game_type"`
}

func (e EndPointHandler) CreateGame(c *gin.Context) {
	var req StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	if e.GameManager.GameStatus == Created {
		c.JSON(400, gin.H{"error": "A game is already created"})
		return
	}
	if e.GameManager.GameStatus == Started {
		c.JSON(400, gin.H{"error": "A game is already running"})
		return
	}
	err := e.GameManager.CreateNewGame(req.GameType)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start game", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "New game created"})
}

func (e EndPointHandler) StartGame(c *gin.Context) {
	if e.GameManager.GameStatus == Idle {
		c.JSON(400, gin.H{"error": "There is no game to start"})
		return
	}
	if e.GameManager.GameStatus == Started {
		c.JSON(400, gin.H{"error": "A game is already running"})
		return
	}
	err := e.GameManager.StartGame()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start game", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Game started"})
}

func (e EndPointHandler) StopGame(c *gin.Context) {
	if e.GameManager.GameStatus == Idle {
		c.JSON(400, gin.H{"error": "No game has been created"})
		return
	}
	if e.GameManager.GameStatus == Created {
		c.JSON(400, gin.H{"error": "No game is running"})
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
	username, exists := c.Get("username")
	if !exists {
		c.JSON(400, gin.H{"error": "no username in context"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "username in context is not a string"})
		return
	}
	err := e.Repo.DeleteUser(usernameStr)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user from Database", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "User deleted"})
}

func (e EndPointHandler) JoinGame(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(400, gin.H{"error": "no username in context"})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(400, gin.H{"error": "username in context is not a string"})
		return
	}

	user, err := e.Repo.GetUser(usernameStr)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found in database", "details": err.Error()})
	}
	player := Player{
		Username:   usernameStr,
		PiSN:       user.PiSN,
		DeathSound: user.DeathSound,
	}
	if e.GameManager.GameStatus == Idle {
		c.JSON(400, gin.H{"error": "No game to join"})
		return
	}
	if e.GameManager.GameStatus == Started {
		c.JSON(400, gin.H{"error": "A game is already running"})
		return
	}
	err = e.GameManager.AddPlayer(player)
	if err != nil {
		c.JSON(500, gin.H{"error": "couldn't add player", "details": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Successfully joined game"})
}
