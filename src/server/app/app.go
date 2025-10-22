package app

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/risasim/projectM5/project/src/server/auth"
	"github.com/risasim/projectM5/project/src/server/communication"
	"github.com/risasim/projectM5/project/src/server/db"
	"github.com/risasim/projectM5/project/src/server/state"
)

// config does load all the configuration details from the .env file
type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

// LoadConfig does load the data from the .env file
func loadConfig() *config {
	return &config{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", "mydb"),
		DBSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "jwt_secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// App holds the db and the gameManager in one structure
type App struct {
	DB           *sql.DB
	UserRepo     db.UserRepositoryInterface
	Routes       *gin.Engine
	GameManager  *state.GameManager
	loginHandler *auth.LoginHandler
}

// CreateConnection opens the connection with the db via the .env values
func (a *App) CreateConnection() {
	var config *config = loadConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, config.DBSSLMode)

	// Retry connection with exponential backoff
	var db *sql.DB
	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Failed to open database connection (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Test the connection
		if err = db.Ping(); err != nil {
			log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
			db.Close()
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		log.Println("Database connection established successfully")
		break
	}

	if err != nil {
		log.Fatal("Failed to connect to database after all retries:", err)
	}

	a.DB = db
}

// SetupLogin sets up the login handler
func (a *App) SetupLogin() {
	var config *config = loadConfig()
	if a.UserRepo == nil {
		a.UserRepo = db.NewUsersRepository(a.DB)
	}
	a.loginHandler = auth.NewLoginHandler(a.UserRepo, []byte(config.JWTSecret), "60")
}

// Migrate does runs the migrations in the db/migrations folder
func (a *App) Migrate() {
	driver, err := postgres.WithInstance(a.DB, &postgres.Config{})
	if err != nil {
		log.Println(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Steps(2); err != nil {
		log.Println(err)
	}
}

func (a *App) InitDatabase() {
	a.CreateConnection()
	a.Migrate()
	db.SeedAdmin(a.DB)
}

func (a *App) CreateRoutes() {
	routes := gin.Default()
	routes.POST("/auth", a.loginHandler.Login)
	userController := db.NewUserController(a.DB)
	routes.POST("/piAuth", a.loginHandler.PiLogin)

	protected := routes.Group("/api")
	protected.Use(a.loginHandler.AuthenticationMiddleware)

	protected.GET("/users", userController.GetUsers)
	protected.POST("/addUser", userController.InsertUser)
	//For web
	protected.POST("/uploadSound", func(c *gin.Context) {
		username := c.Query("username")
		repo := db.NewUsersRepository(a.DB)
		user, err := repo.GetUser(username)
		if err != nil || user == nil {
			c.JSON(404, gin.H{"error": "User not found in database"})
			return
		}
		if a.GameManager.GameStatus == state.Active {
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
		_, err = a.DB.Exec("UPDATE users SET deathSound=$1 WHERE username=$2", b64Sound, username)
		if err != nil {
			c.JSON(500, gin.H{"error": "Database deathSound update failed"})
			return
		}
		err = a.GameManager.SendNewMusicToPi(username, b64Sound, file.Filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "sendNewMusicToPi failed"})
			return
		}

		c.JSON(200, gin.H{"message": "Successfully updated death sound"})
	})
	protected.GET("/gameStatus", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": a.GameManager.GameStatus.String()})
	})
	protected.POST("/startGame", func(c *gin.Context) {
		gameTypeString := c.Query("GameType")
		gameType, _ := communication.ParseGameType(gameTypeString)
		if a.GameManager.GameStatus == state.Active {
			c.JSON(400, gin.H{"error": "A game is already active"})
			return
		}
		a.GameManager.StartNewGame(gameType)
		c.JSON(200, gin.H{"message": "New game started"})
	})
	protected.POST("/stopGame", func(c *gin.Context) {
		a.GameManager.EndGame()
		c.JSON(200, gin.H{"message": "Game Stopped"})
	})
	//protected.POST("/gameStatus") ???
	protected.DELETE("/user", func(c *gin.Context) {
		username := c.Query("username")
		_, err := a.DB.Exec("DELETE FROM users WHERE username=$1", username)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete user from Database"})
			return
		}
		c.JSON(200, gin.H{"message": "User deleted"})
	})
	protected.POST("joinGame", func(c *gin.Context) {
		username := c.Query("username")
		repo := db.NewUsersRepository(a.DB)
		user, err := repo.GetUser(username)
		if err != nil {
			c.JSON(404, gin.H{"error": "User not found in database"})
		}
		player := state.Player{
			Username:   username,
			ID:         int(user.ID),
			EncodingID: 0,
		}
		err = a.GameManager.AddPlayer(player)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Successfully joined game"})
	})
	// For pi
	//protected.GET("/music") connected this to the POST from website

	protected.GET("/wsLeaderboard", a.GameManager.WsLeaderBoardHandler)
	protected.GET("/wsPis", a.GameManager.WsPisHandler)

	a.Routes = routes
}

func (a *App) Run(gm *state.GameManager) {
	a.GameManager = gm
	a.Routes.Run(":8080")
}
