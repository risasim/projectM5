package tests

import (
	"bytes"
	"encoding/json"
	"github.com/risasim/projectM5/project/src/server/tests/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PiRequest struct {
	ApiKey string `json:"apiKey"`
	PiSn   string `json:"piSn"`
}

func TestAuthLogin_Success(t *testing.T) {
	ta := mock.SetupTestApp(t)

	// Prepare request body
	body, _ := json.Marshal(LoginRequest{
		Username: "testuser",
		Password: "secret123",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP 200 for valid login")
	assert.Contains(t, w.Body.String(), "token", "Response should contain a token")
}

func TestAuthLogin_InvalidCredentials(t *testing.T) {
	ta := mock.SetupTestApp(t)

	body, _ := json.Marshal(LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 for invalid credentials")
}

func TestPiAuth_ValidCredentials(t *testing.T) {
	ta := mock.SetupTestApp(t)

	body, _ := json.Marshal(PiRequest{
		PiSn:   "pi-0002",
		ApiKey: "kazoo",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/piAuth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 for valid token")
	assert.Contains(t, w.Body.String(), "token", "Response should contain a token")
}

func TestAuthLogin_MissingApiKey(t *testing.T) {
	ta := mock.SetupTestApp(t)

	body, _ := json.Marshal(PiRequest{
		PiSn:   "pi-0002",
		ApiKey: "", // missing
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthLogin_MissingSnKey(t *testing.T) {
	ta := mock.SetupTestApp(t)

	body, _ := json.Marshal(PiRequest{
		PiSn:   "",
		ApiKey: "kazoo",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthLogin_NonExistantPiSN(t *testing.T) {
	ta := mock.SetupTestApp(t)

	body, _ := json.Marshal(PiRequest{
		PiSn:   "34",
		ApiKey: "kazoo",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
