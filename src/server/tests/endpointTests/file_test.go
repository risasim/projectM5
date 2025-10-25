package endpointTests

import (
	"bytes"
	"github.com/risasim/projectM5/project/src/server/tests/mock"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadSoundWithoutPreset(t *testing.T) {
	ta := mock.SetupTestApp(t)

	filePath := "testSound.mp3"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("testSound.mp3 not found in current directory")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("failed to open %s: %v", filePath, err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("sound", filepath.Base(filePath))
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		t.Fatalf("failed to copy file content: %v", err)
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/api/uploadSound", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+ta.Token)

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	expected := `"status":"success"`
	if !bytes.Contains(w.Body.Bytes(), []byte(expected)) {
		t.Errorf("expected body to contain %s, got %s", expected, w.Body.String())
	}
}

func TestGetSound(t *testing.T) {
	ta := mock.SetupTestApp(t)
	req, _ := http.NewRequest("GET", "/api/sound", nil)
	req.Header.Set("Authorization", "Bearer "+ta.Token)

	w := httptest.NewRecorder()
	ta.App.Routes.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "audio/mpeg", w.Header().Get("Content-Type"), w.Body.String())
}
