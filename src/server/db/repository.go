package db

import (
	"database/sql"
	"gitlab.utwente.nl/computer-systems-project/2025-2026/students-projects/cs25-29/project/src/server/db/model"
	"log"
)

type UserRepositoryInterface interface {
	SelectUsers() []model.GetUserResponse
	InsertUser(user model.PostUser, apiKey string, isAdmin bool) bool
	GetUser() model.GetUserAuth
}

// UsersRepository does execute the sql calls on the db
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository is a constructor for the UsersRepository
func NewUsersRepository(db *sql.DB) UserRepositoryInterface {
	return &UsersRepository{db: db}
}

// SelectUsers returns all the users that are in the db
func (u UsersRepository) SelectUsers() []model.GetUserResponse {
	var users []model.GetUserResponse
	rows, err := u.db.Query("SELECT id,isAdmin,username FROM users ")
	if err != nil {
		log.Println(err)
		return nil
	}
	for rows.Next() {
		var (
			id       uint
			isAdmin  bool
			username string
		)
		err = rows.Scan(&id, &isAdmin, &username)
		if err != nil {
			log.Println(err)
		} else {
			user := model.GetUserResponse{ID: id, IsAdmin: isAdmin, Username: username}
			users = append(users, user)
		}
	}
	return users
}

// InsertUser does insert the user into the postgressql database, it need the apiKey and piSN to be gene
func (u UsersRepository) InsertUser(user model.PostUser, apiKey string, isAdmin bool) bool {
	stmt, err := u.db.Prepare("INSERT INTO users (isAdmin,username, password, api_key,pi_SN) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(isAdmin, user.Username, user.Password, apiKey, user.PiSN)
	if err2 != nil {
		log.Println(err2)
		return false
	}
	return true
}

func (u UsersRepository) GetUser() model.GetUserAuth {
	panic("implement me")
}
