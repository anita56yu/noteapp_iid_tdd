package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/repository"
	"noteapp/internal/usecase"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// setupTest initializes the necessary components for the tests.
func setupTest() http.Handler {
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := usecase.NewNoteUsecase(repo)
	handler := NewNoteHandler(noteUsecase)

	router := chi.NewRouter()
	router.Post("/notes", handler.CreateNote)
	return router
}

func TestNoteHandler_CreateNote_Success(t *testing.T) {
	// Arrange
	router := setupTest()
	requestBody := CreateNoteRequest{
		Title:   "Test Title",
		Content: "Test Content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}
	location := rr.Header().Get("Location")
	if !strings.HasPrefix(location, "/notes/") {
		t.Errorf("expected Location header to start with '/notes/'; got '%s'", location)
	}
	var responseBody struct{ ID string `json:"id"` }
	if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if responseBody.ID == "" {
		t.Error("expected a non-empty ID in the response body")
	}
	idFromLocation := strings.TrimPrefix(location, "/notes/")
	if responseBody.ID != idFromLocation {
		t.Errorf("ID in body ('%s') does not match ID in Location header ('%s')", responseBody.ID, idFromLocation)
	}
}

func TestNoteHandler_CreateNote_InvalidJSON(t *testing.T) {
	// Arrange
	router := setupTest()
	invalidBody := []byte(`{"title": "Test", "content":`) // Malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(invalidBody))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_CreateNote_EmptyTitle(t *testing.T) {
	// Arrange
	router := setupTest()
	requestBody := CreateNoteRequest{Title: "", Content: "Test Content"}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
	expectedErr := "title cannot be empty"
	if !strings.Contains(rr.Body.String(), expectedErr) {
		t.Errorf("expected error message '%s' in body; got '%s'", expectedErr, rr.Body.String())
	}
}
