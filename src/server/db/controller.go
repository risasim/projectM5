package db

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/risasim/projectM5/project/src/server/auth"
	"github.com/risasim/projectM5/project/src/server/db/model"
)

type UserControllerInterface interface {
	InsertUser(g *gin.Context)
	GetUsers(g *gin.Context)
}

type UserController struct {
	db *sql.DB
}

func (u UserController) InsertUser(g *gin.Context) {
	// TODO implement middleware for authentication
	db := u.db
	var user model.PostUser
	if err := g.ShouldBindJSON(&user); err == nil {
		hash, _ := auth.HashPassword(user.Password)
		user.Password = hash
		usersRepo := NewUsersRepository(db)
		apiKey := uuid.New().String()
		insert := usersRepo.InsertUser(user, apiKey, false)
		if insert {
			g.JSON(200, gin.H{"status": "sucess", "msg": "Inserted new User"})
		} else {
			g.JSON(500, gin.H{"status": "fail", "msg": "Something went wrong"})
		}
	} else {
		g.JSON(400, gin.H{"status": "fail", "msg": "Failed to dework "})
	}

}

func (u UserController) GetUsers(g *gin.Context) {
	db := u.db
	usersRepo := NewUsersRepository(db)
	getUsers := usersRepo.SelectUsers()
	if getUsers != nil {
		g.JSON(200, gin.H{"status": "success", "data": getUsers, "msg": "Sending users"})
	} else {
		g.JSON(500, gin.H{"status": "failure", "data": nil, "msg": "Someting went wrong in the db"})
	}
}

func NewUserController(db *sql.DB) UserControllerInterface {
	return &UserController{db: db}
}
