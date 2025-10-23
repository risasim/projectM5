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

func TestAuthLogin_Success(t *testing.T) {
	ta := mock.SetupTestApp(t)

	// Prepare request body
	body, _ := json.Marshal(LoginRequest{
		Username: "testuser",
		Password: "secret123",
	})

	req, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
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

	req, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 for invalid credentials")
}
