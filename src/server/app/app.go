package app

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// config does load all the configuration details from the .env file
type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// LoadConfig does load the data from the .env file
func loadConfig() *config {
	return &config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "mydb"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
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
	DB *sql.DB
}

// CreateConnection opens the connection with the db via the .env values
func (a *App) CreateConnection() {
	var config *config = loadConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", config.DBUser, config.DBPassword, config.DBHost, config.DBName, config.DBSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
}

// Migrate does runs the migrations in the db/migrations folder
func (a *App) Migrate() {
	driver, err := postgres.WithInstance(a.DB, &postgres.Config{})
	if err != nil {
		log.Println(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations/",
		"PhoShoData", driver)
	if err != nil {
		log.Println(err)
	}
	if err := m.Steps(2); err != nil {
		log.Println(err)
	}
}

func (a *App) CreateRoutes() {
	//routes := gin.Default()
	//userController :=
}
