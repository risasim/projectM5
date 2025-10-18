package db

import (
	"database/sql"
	"log"
)

func SeedAdmin(db *sql.DB) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE isAdmin = TRUE)").Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check for admin user:", err)
	}

	if !exists {
		passwordHash, err := HashPassword("adminpass")
		if err != nil {
			log.Fatal("Failed to hash admin password:", err)
		}

		_, err = db.Exec(`INSERT INTO users (username, password, isAdmin) VALUES ($1, $2, $3)`,
			"admin", passwordHash, true)
		if err != nil {
			log.Fatal("Failed to create admin user:", err)
		}

		log.Println("Admin user created: username=admin, password=adminpass")
	}
}
