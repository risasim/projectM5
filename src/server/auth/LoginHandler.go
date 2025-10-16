package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/risasim/projectM5/project/src/server/db"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginHandler struct {
	userRepository db.UsersRepository
	secretKey      []byte
	timeoutStr     string
}

func newLoginHandler(repo db.UsersRepository, secretKey []byte, timeoutStr string) *LoginHandler {
	return &LoginHandler{
		userRepository: repo,
		secretKey:      secretKey,
		timeoutStr:     timeoutStr,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.GetUser(credentials.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !db.VerifyPassword(credentials.Password, user.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := h.createToken(user.Username, user.IsAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *LoginHandler) createToken(username string, isAdmin bool) (string, error) {
	timeoutMinutes, err := strconv.Atoi(h.timeoutStr)
	if err != nil {
		timeoutMinutes = 60
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"admin":    isAdmin,
			"exp":      time.Now().Add(time.Minute * time.Duration(timeoutMinutes)).Unix(),
		})

	return token.SignedString(h.secretKey)
}

func (h *LoginHandler) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid token")
	}
	return nil
}
