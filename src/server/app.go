package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/risasim/projectM5/project/src/server/auth"
	"github.com/risasim/projectM5/project/src/server/db"
	"github.com/risasim/projectM5/project/src/server/state"
	"log"
	"os"
	"time"
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
	Routes       *gin.Engine
	GameManager  *state.GameManager
	upgrader     websocket.Upgrader
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
	repo := db.NewUsersRepository(a.DB)
	a.loginHandler = auth.NewLoginHandler(repo, []byte(config.JWTSecret), "60")
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

	protected := routes.Group("/api")
	protected.Use(a.loginHandler.AuthenticationMiddleware)

	protected.GET("/users", userController.GetUsers)
	protected.POST("/addUser", userController.InsertUser)
	//For web
	//protected.POST("/music")
	//protected.GET("/gameStatus")
	//protected.POST("/startGame")
	//Ending the game ?
	//protected.POST("/gameStatus")
	//protected.DELETE("/user")
	//protected.POSt("joinGame")
	// For pi
	//protected.GET("/music")

	//routes.GET("/wsLeaderboard")
	//routes.GET("/wsPis")
	a.Routes = routes
	a.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

func (a *App) Run(gm *state.GameManager) {
	a.GameManager = gm
	a.Routes.Run(":8080")
}
