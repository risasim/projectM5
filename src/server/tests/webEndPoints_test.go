package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/risasim/projectM5/project/src/server/communication"
	"github.com/risasim/projectM5/project/src/server/state"
	"github.com/risasim/projectM5/project/src/server/tests/mock"
	"github.com/stretchr/testify/assert"
)

func setUpRouter(handler *state.EndPointHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("uploadSound", handler.UploadSound)
	router.GET("/sound", handler.GetSound)
	router.GET("/gameStatus", handler.GetGameStatus)
	router.POST("/createGame", handler.CreateGame)
	router.POST("startGame", handler.StartGame)
	router.POST("/stopGame", handler.StopGame)
	router.POST("/joinGame", handler.JoinGame)
	router.DELETE("/deleteUser", handler.DeleteUser)
	return router
}

func TestGetGameStatus(t *testing.T) {
	tests := []struct {
		status     state.GameStatus
		wantStatus string
	}{
		{state.Idle, state.Idle.String()},
		{state.Created, state.Created.String()},
		{state.Started, state.Started.String()},
	}
	for _, test := range tests {

		handler := state.EndPointHandler{Repo: nil, GameManager: &state.GameManager{GameStatus: test.status}}
		router := setUpRouter(&handler)
		req := httptest.NewRequest(http.MethodGet, "/gameStatus", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, test.wantStatus, response["Game_Status"])
	}
}

func TestCreateGame_Statuses(t *testing.T) {
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	for _, test := range tests {
		handler := state.EndPointHandler{Repo: nil, GameManager: &state.GameManager{GameStatus: test.status}}
		router := setUpRouter(&handler)
		payload := state.StartGameRequest{GameType: communication.Freefall}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/createGame", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if test.status != state.Idle {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestCreateGame_InvalidJSON(t *testing.T) {
	handler := state.EndPointHandler{Repo: nil, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := setUpRouter(&handler)
	//Colon Missing
	body := []byte(`"game_type" "Freefall"`)
	req := httptest.NewRequest(http.MethodPost, "/createGame", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request", response["error"])
	assert.NotEmpty(t, response["details"])
}

// TODO Complete after fixing start game
func TestStartGame(t *testing.T) {
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	for _, test := range tests {
		handler := state.EndPointHandler{Repo: nil, GameManager: &state.GameManager{GameStatus: test.status}}
		router := setUpRouter(&handler)
		req := httptest.NewRequest(http.MethodPost, "/startGame", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if test.status != state.Created {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

// TODO Complete after FIXING end game
func TestStopGame(t *testing.T) {
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	for _, test := range tests {
		handler := state.EndPointHandler{Repo: nil, GameManager: &state.GameManager{GameStatus: test.status}}
		router := setUpRouter(&handler)
		req := httptest.NewRequest(http.MethodPost, "/stopGame", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if test.status != state.Started {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

// TODO fix no session to join so get error
func TestJoinGame(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	for _, test := range tests {
		handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: test.status}}
		router := gin.Default()
		router.POST("/joinGame", func(c *gin.Context) {
			c.Set("username", "testuser")
			handler.JoinGame(c)
		})
		req := httptest.NewRequest(http.MethodPost, "/joinGame", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if test.status != state.Created {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := gin.Default()
	router.POST("/deleteUser", func(c *gin.Context) {
		c.Set("username", "testuser")
		handler.DeleteUser(c)
	})
	req := httptest.NewRequest(http.MethodPost, "/deleteUser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func createFileUploadRequest(filename string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("sound", filename)
	_, _ = part.Write([]byte{})
	_ = writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/uploadSound", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUploadSoundMp3AndNonMp3(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := gin.Default()
	router.POST("/uploadSound", func(c *gin.Context) {
		c.Set("username", "testuser")
		handler.UploadSound(c)
	})

	mp3req := createFileUploadRequest("sound.mp3")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, mp3req)
	assert.Equal(t, http.StatusOK, w.Code)

	nonMp3req := createFileUploadRequest("sound.txt")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, nonMp3req)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
}

func TestDownloadSoundUserNotFound(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := gin.Default()
	router.POST("/uploadSound", func(c *gin.Context) {
		c.Set("username", "NonExistentUser")
		handler.UploadSound(c)
	})

	mp3req := createFileUploadRequest("sound.mp3")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, mp3req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUploadSoundUserNotAString(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := gin.Default()
	router.POST("/uploadSound", func(c *gin.Context) {
		c.Set("username", 123)
		handler.UploadSound(c)
	})
	mp3req := createFileUploadRequest("sound.mp3")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, mp3req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadSoundUserNoUsername(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := gin.Default()
	router.POST("/uploadSound", func(c *gin.Context) {
		handler.UploadSound(c)
	})
	mp3req := createFileUploadRequest("sound.mp3")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, mp3req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSoundFile(t *testing.T) {
	testApp := mock.SetupTestApp(t)
	handler := state.EndPointHandler{Repo: testApp.MockRepo, GameManager: &state.GameManager{GameStatus: state.Idle}}
	router := gin.Default()
	router.GET("/getSound", func(c *gin.Context) {
		c.Set("username", "testuser")
		handler.GetSound(c)
	})
	existingReq := httptest.NewRequest(http.MethodGet, "/getSound", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, existingReq)
	assert.Equal(t, http.StatusOK, w.Code)
}
