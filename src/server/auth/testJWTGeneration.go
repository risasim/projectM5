package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateTestJWT creates a signed JWT for testing purposes.
func GenerateTestJWT(username string, isAdmin bool, secret []byte, durationMinutes int) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"isAdmin":  isAdmin,
		"exp":      time.Now().Add(time.Duration(durationMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
