package mock

import (
	"database/sql"
	"github.com/risasim/projectM5/project/src/server/db"
	"github.com/risasim/projectM5/project/src/server/db/model"
	"sync"
)

// MockUserRepository is a mock of user repository used for testing
type MockUserRepository struct {
	mu        sync.Mutex
	users     map[string]model.GetUserAuth
	usersByPi map[string]*model.GetUserAuth
	autoid    uint
}

func (m *MockUserRepository) GetPiUser(piSN string) (*model.GetUserAuth, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, exists := m.usersByPi[piSN]; exists {
		return u, nil
	}
	return nil, nil
}

// SelectUsers does mock the real
func (m MockUserRepository) SelectUsers() []model.GetUserResponse {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []model.GetUserResponse
	for _, u := range m.users {
		result = append(result, model.GetUserResponse{
			ID:       u.ID,
			Username: u.Username,
			IsAdmin:  u.IsAdmin,
		})
	}
	return result
}

// InsertUser does mock the behaviour of its real db counterpart
func (m MockUserRepository) InsertUser(user model.PostUser, apiKey string, isAdmin bool) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Username]; exists {
		return false
	}

	newUser := model.GetUserAuth{
		ID:         m.autoid,
		IsAdmin:    isAdmin,
		Username:   user.Username,
		Password:   user.Password,
		DeathSound: user.DeathSound,
		PiSN:       sqlNullString(user.PiSN),
		ApiKey:     sqlNullString(apiKey),
	}
	m.users[user.Username] = newUser

	if newUser.PiSN.Valid {
		m.usersByPi[newUser.PiSN.String] = &newUser
	}

	m.autoid++
	return true
}

// GetUser retrieves a user by username.
func (m *MockUserRepository) GetUser(username string) (*model.GetUserAuth, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u, exists := m.users[username]; exists {
		return &u, nil
	}
	return nil, nil // mimic sql.ErrNoRows => nil,nil
}

func NewMockUserRepository() db.UserRepositoryInterface {
	return &MockUserRepository{
		users:     make(map[string]model.GetUserAuth),
		usersByPi: make(map[string]*model.GetUserAuth),
		autoid:    1,
	}
}

// sqlNullString is helper to build sql.NullString
func sqlNullString(s string) (ns sql.NullString) {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
