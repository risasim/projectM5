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

type PiCredentials struct {
	ApiKey string `json:"apiKey"`
	PiSn   string `json:"piSn"`
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

func CheckAdmin(ctx *gin.Context) {
	adminVal, exists := ctx.Get("isAdmin")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing admin claim"})
		return
	}
	isAdmin, ok := adminVal.(bool)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid admin claim type"})
		return
	}
	if !isAdmin {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "fuck you"})
		return
	}
	ctx.Next()
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
	if user.IsAdmin {
		ctx.JSON(http.StatusOK, gin.H{"token": tokenString, "role": "admin"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"token": tokenString, "role": "user"})
	}
}

func (h *LoginHandler) PiLogin(ctx *gin.Context) {
	var credentials PiCredentials
	err := ctx.ShouldBindJSON(&credentials)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	user, err := h.userRepository.GetPiUser(credentials.PiSn)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if user == nil {
		ctx.JSON(401, gin.H{"error": "invalid device or API key"})
		return
	}

	if !user.ApiKey.Valid {
		ctx.JSON(401, gin.H{"error": "fuck you"})
		return
	}

	if !db.VerifyPassword(credentials.ApiKey, user.ApiKey.String) {
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

// ParseToken does parse the JWT and returns the claims store in it
func (h *LoginHandler) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}
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
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	tokenString := parts[1]

	claims, err := h.ParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Store th claim in the context that is passed to user
	c.Set("claims", claims)
	c.Set("username", claims["username"])
	c.Set("isAdmin", claims["admin"])

	c.Next()
}

func (h *LoginHandler) WSQueryAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	var tokenString string
	if authHeader != "" {
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenString = parts[1]
	} else {
		tokenString = c.Query("token")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization missing"})
			return
		}
	}

	claims, err := h.ParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.Set("claims", claims)
	c.Set("username", claims["username"])
	c.Set("isAdmin", claims["admin"])
	c.Next()
}
