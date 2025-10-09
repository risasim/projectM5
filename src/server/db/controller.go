package db

import (
	"database/sql"
	"github.com/gin-gonic/gin"
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
	//db := u.db
	//var user := model.PostUser
	//if err := g.ShouldBindJSON(&user); err == nil {
	//	usersRepo := NewUsersRepository(db)
	//	//TODO implement the apikey generation and the password hashing
	//	insert := usersRepo.InsertUser(user,)
	//}
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

func NewUserControllers(db *sql.DB) UserControllerInterface {
	return &UserController{db: db}
}
