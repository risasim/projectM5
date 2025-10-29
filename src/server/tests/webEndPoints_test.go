package tests

import (
	"bytes"
	"encoding/json"
	"github.com/risasim/projectM5/project/src/server/auth"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/risasim/projectM5/project/src/server/communication"
	"github.com/risasim/projectM5/project/src/server/state"
	"github.com/risasim/projectM5/project/src/server/tests/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetGameStatus(t *testing.T) {
	ta := mock.SetupTestApp(t)
	tests := []struct {
		status     state.GameStatus
		wantStatus string
	}{
		{state.Idle, state.Idle.String()},
		{state.Created, state.Created.String()},
		{state.Started, state.Started.String()},
	}
	for _, test := range tests {
		ta.App.GameManager.GameStatus = test.status
		req := httptest.NewRequest(http.MethodGet, "/api/gameStatus", nil)
		req.Header.Set("Authorization", "Bearer "+ta.Token)
		w := httptest.NewRecorder()
		ta.App.Routes.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, test.wantStatus, response["Game_Status"])
	}
}

func TestCreateGame_Statuses(t *testing.T) {
	ta := mock.SetupTestApp(t)
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	for _, test := range tests {
		ta.App.GameManager.GameStatus = test.status
		payload := state.StartGameRequest{GameType: communication.Freefall}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/createGame", bytes.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ta.Token)
		w := httptest.NewRecorder()
		ta.App.Routes.ServeHTTP(w, req)
		if test.status != state.Idle {
			println(w.Body.String())
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestCreateGame_InvalidJSON(t *testing.T) {
	ta := mock.SetupTestApp(t)
	//Colon Missing
	body := []byte(`"game_type" "Freefall"`)
	req := httptest.NewRequest(http.MethodPost, "/api/createGame", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ta.Token)
	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request", response["error"])
	assert.NotEmpty(t, response["details"])
}

func TestStartGame(t *testing.T) {
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	ta := mock.SetupTestApp(t)
	for _, test := range tests {
		// Set the game status before making the request
		ta.App.GameManager.GameStatus = test.status

		req := httptest.NewRequest(http.MethodPost, "/api/startGame", nil)
		req.Header.Set("Authorization", "Bearer "+ta.Token)
		w := httptest.NewRecorder()
		ta.App.Routes.ServeHTTP(w, req)
		if test.status != state.Created {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestStartGameNoAdmin(t *testing.T) {
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	ta := mock.SetupTestApp(t)
	badToken, _ := auth.GenerateTestJWT("testuser", false, []byte("jwt_secret"), 60)
	for _, test := range tests {
		// Set the game status before making the request
		ta.App.GameManager.GameStatus = test.status

		req := httptest.NewRequest(http.MethodPost, "/api/startGame", nil)
		req.Header.Set("Authorization", "Bearer "+badToken)
		w := httptest.NewRecorder()
		ta.App.Routes.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	}
}

func TestStopGame(t *testing.T) {

	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	ta := mock.SetupTestApp(t)
	for _, test := range tests {
		ta.App.GameManager.GameStatus = test.status
		req := httptest.NewRequest(http.MethodPost, "/api/stopGame", nil)
		req.Header.Set("Authorization", "Bearer "+ta.Token)
		w := httptest.NewRecorder()
		ta.App.Routes.ServeHTTP(w, req)
		if test.status != state.Started {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestJoinGame(t *testing.T) {
	ta := mock.SetupTestApp(t)
	tests := []struct {
		status state.GameStatus
	}{
		{state.Idle},
		{state.Created},
		{state.Started},
	}
	for _, test := range tests {
		ta.App.GameManager.GameStatus = test.status
		req := httptest.NewRequest(http.MethodPost, "/api/joinGame", nil)
		req.Header.Set("Authorization", "Bearer "+ta.Token)
		w := httptest.NewRecorder()
		ta.App.Routes.ServeHTTP(w, req)
		if test.status != state.Created {
			assert.Equal(t, http.StatusBadRequest, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	ta := mock.SetupTestApp(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/user", nil)
	req.Header.Set("Authorization", "Bearer "+ta.Token)
	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func createFileUploadRequest(filename string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("sound", filename)
	_, _ = part.Write([]byte{})
	_ = writer.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/uploadSound", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUploadSoundMp3AndNonMp3(t *testing.T) {
	ta := mock.SetupTestApp(t)

	mp3req := createFileUploadRequest("sound.mp3")
	mp3req.Header.Set("Authorization", "Bearer "+ta.Token)
	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, mp3req)
	assert.Equal(t, http.StatusOK, w.Code)

	nonMp3req := createFileUploadRequest("sound.txt")
	nonMp3req.Header.Set("Authorization", "Bearer "+ta.Token)
	w1 := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w1, nonMp3req)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
}

func TestGetSoundFile(t *testing.T) {
	ta := mock.SetupTestApp(t)
	existingReq := httptest.NewRequest(http.MethodGet, "/api/sound", nil)
	existingReq.Header.Add("Authorization", "Bearer "+ta.Token)
	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, existingReq)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadSoundUserNotFound(t *testing.T) {
	ta := mock.SetupTestApp(t)
	mp3req := createFileUploadRequest("sound.mp3")
	mp3req.Header.Set("Authorization", "Bearer "+ta.NonExistentToken)
	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, mp3req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
