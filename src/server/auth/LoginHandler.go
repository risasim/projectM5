package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/risasim/projectM5/project/src/server/db"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginHandler struct {
	userRepository db.UserRepositoryInterface
	secretKey      []byte
	timeoutStr     string
}

func NewLoginHandler(repo db.UserRepositoryInterface, secretKey []byte, timeoutStr string) *LoginHandler {
	return &LoginHandler{
		userRepository: repo,
		secretKey:      secretKey,
		timeoutStr:     timeoutStr,
	}
}

func (h *LoginHandler) Login(ctx *gin.Context) {
	var credentials Credentials
	err := ctx.ShouldBindJSON(&credentials)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	user, err := h.userRepository.GetUser(credentials.Username)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !db.VerifyPassword(credentials.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	tokenString, err := h.createToken(user.Username, user.IsAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "opsie"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
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

func (h *LoginHandler) VerifyToken(tokenString string) error {
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

func (h *LoginHandler) AuthenticationMiddleware(c *gin.Context) {
	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Split the header to get the token part
	tokenString := strings.Split(authHeader, "Bearer ")[1]

	err := h.VerifyToken(tokenString)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Token is valid, proceed with the request
	c.Next()
}
