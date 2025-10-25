package db

import (
	"database/sql"
	"fmt"
	"github.com/risasim/projectM5/project/src/server/db/model"
	"log"
)

type UserRepositoryInterface interface {
	SelectUsers() []model.GetUserResponse
	InsertUser(user model.PostUser, apiKey string, isAdmin bool) bool
	GetUser(username string) (*model.GetUserAuth, error)
	GetPiUser(piSN string) (*model.GetUserAuth, error)
	UpdateDeathSound(username string, path string) error
	DeleteUser(username string) error
}

// UsersRepository does execute the sql calls on the db
type UsersRepository struct {
	db *sql.DB
}

func (u UsersRepository) UpdateDeathSound(username, path string) error {
	result, err := u.db.Exec("UPDATE users SET deathSound=$1 WHERE username=$2", path, username)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no user found with username %s", username)
	}
	return nil
}

func (u UsersRepository) DeleteUser(username string) error {
	result, err := u.db.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no user found with username %s", username)
	}
	return nil
}

func (u UsersRepository) GetPiUser(piSN string) (*model.GetUserAuth, error) {
	var user model.GetUserAuth
	err := u.db.QueryRow("SELECT * FROM users WHERE pi_SN = $1 ", piSN).Scan(
		&user.ID,
		&user.IsAdmin,
		&user.Username,
		&user.Password,
		&user.DeathSound,
		&user.PiSN,
		&user.ApiKey,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no user found
		}
		log.Println("GetUser error:", err)
		return nil, err
	}
	return &user, nil
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
	stmt, err := u.db.Prepare("INSERT INTO users (isAdmin,username, password,deathSound, api_key,pi_SN) VALUES ($1, $2, $3, $4, $5,$6)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err2 := stmt.Exec(isAdmin, user.Username, user.Password, user.DeathSound, apiKey, user.PiSN)
	if err2 != nil {
		log.Println(err2)
		return false
	}
	return true
}

func (u UsersRepository) GetUser(username string) (*model.GetUserAuth, error) {
	var user model.GetUserAuth
	err := u.db.QueryRow("SELECT * FROM users WHERE username = $1 ", username).Scan(
		&user.ID,
		&user.IsAdmin,
		&user.Username,
		&user.Password,
		&user.DeathSound,
		&user.PiSN,
		&user.ApiKey,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no user found
		}
		log.Println("GetUser error:", err)
		return nil, err
	}
	return &user, nil
}
