package mock

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/risasim/projectM5/project/src/server/app"
	"github.com/risasim/projectM5/project/src/server/auth"
	"github.com/risasim/projectM5/project/src/server/db"
	"github.com/risasim/projectM5/project/src/server/db/model"
	"github.com/risasim/projectM5/project/src/server/state"
)

// TestApp should mock the app, with predifined parts
type TestApp struct {
	App      *app.App
	MockRepo db.UserRepositoryInterface
	// Token is JWT token that is generated, for easier testing
	Token            string
	NonExistentToken string
}

func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()

	gin.SetMode(gin.TestMode)

	hashedPassword, err := db.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	mockRepo := NewMockUserRepository()
	_ = mockRepo.InsertUser(
		model.PostUser{
			Username:   "testuser",
			Password:   hashedPassword,
			DeathSound: "default.mp3",
			PiSN:       "pi-0001",
		}, "", true)

	hashedApiKey, err := db.HashPassword("kazoo")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_ = mockRepo.InsertUser(
		model.PostUser{
			Username:   "testuser2",
			Password:   hashedPassword,
			DeathSound: "default.mp3",
			PiSN:       "pi-0002",
		}, hashedApiKey, false)

	hashedApiKey2, err := db.HashPassword("corndog")
	_ = mockRepo.InsertUser(
		model.PostUser{
			Username:   "testuser3",
			Password:   hashedPassword,
			DeathSound: "default.mp3",
			PiSN:       "pi-0003",
		}, hashedApiKey2, false)

	//Mocking the creation of the real app but wihout the cb
	app := &app.App{}
	app.UserRepo = mockRepo
	app.SetupLogin()

	gameManager := state.NewGameManager()
	// run the broacasters in its own go routines
	go gameManager.BroadcastLeaderBoardHandler()
	go gameManager.BroadcastPisHandler()

	app.GameManager = gameManager
	app.CreateRoutes()

	token, _ := auth.GenerateTestJWT("testuser", true, []byte("jwt_secret"), 60)
	NonExistentToken, _ := auth.GenerateTestJWT("NonExistentUser", false, []byte("jwt_secret"), 60)
	return &TestApp{
		App:              app,
		MockRepo:         mockRepo,
		Token:            token,
		NonExistentToken: NonExistentToken,
	}
}
