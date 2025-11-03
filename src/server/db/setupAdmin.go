package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

// DefaultDeathSound is the filename (not full path) used when a user has no custom sound
var DefaultDeathSound string

func init() {
	soundDir := "soundEffects"

	// Ensure directory exists — both locally and in Docker
	if err := os.MkdirAll(soundDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create soundEffects directory:", err)
	}

	DefaultDeathSound = "default.mp3"

	defaultPath := filepath.Join(soundDir, DefaultDeathSound)
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		log.Printf("⚠️ Warning: Default sound file %s not found; make sure it exists in your container/image", defaultPath)
	}
}

// SeedUsers creates default users if they don't exist
func SeedUsers(db *sql.DB) {
	users := []struct {
		username   string
		password   string
		isAdmin    bool
		deathSound string
		piSN       string
		apiKey     string
	}{
		{
			username:   "admin",
			password:   "adminpass",
			isAdmin:    true,
			deathSound: DefaultDeathSound,
			piSN:       "69",
			apiKey:     "",
		},
		{
			username:   "berk",
			password:   "hamburgers",
			isAdmin:    false,
			deathSound: DefaultDeathSound,
			piSN:       "ae616eb0e54290a6",
			apiKey:     "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			username:   "orbay",
			password:   "wachtword",
			isAdmin:    false,
			deathSound: DefaultDeathSound,
			piSN:       "ae616eb0e54290a69",
			apiKey:     "123e4567-e89b-12d3-a416-426614174000",
		},
		{
			username:   "Chris Kyle",
			password:   "sniper",
			isAdmin:    false,
			deathSound: DefaultDeathSound,
			piSN:       "ae616eb0e54290a70",
			apiKey:     "123e4567-e89b-12d3-a426-426614174000",
		},
		{
			username:   "Peter Griffin",
			password:   "griffi",
			isAdmin:    false,
			deathSound: DefaultDeathSound,
			piSN:       "ae616eb0e54290a74",
			apiKey:     "123e4567-e89b-12d3-a426-4266141740069",
		},
	}

	for _, u := range users {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", u.username).Scan(&exists)
		if err != nil {
			log.Fatal("Failed to check for user:", err)
		}

		if !exists {
			passwordHash, err := HashPassword(u.password)
			if err != nil {
				log.Fatal("Failed to hash password for user", u.username, ":", err)
			}

			var apiKeyHash sql.NullString
			if u.apiKey != "" {
				hashed, err := HashPassword(u.apiKey)
				if err != nil {
					log.Fatal(err)
				}
				apiKeyHash = sql.NullString{String: hashed, Valid: true}
			} else {
				apiKeyHash = sql.NullString{Valid: false}
			}

			_, err = db.Exec(`
				INSERT INTO users (username, password, isAdmin, deathSound, pi_SN, api_key)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				u.username, passwordHash, u.isAdmin, u.deathSound, u.piSN, apiKeyHash)
			if err != nil {
				log.Fatal("Failed to create user", u.username, ":", err)
			}

			log.Printf("✅ User created: username=%s, password=%s\n", u.username, u.password)
		}
	}
}
