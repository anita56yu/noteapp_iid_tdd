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
func setupTest() (*chi.Mux, *usecase.NoteUsecase) {
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := usecase.NewNoteUsecase(repo)
	handler := NewNoteHandler(noteUsecase)

	router := chi.NewRouter()
	router.Post("/notes", handler.CreateNote)
	router.Get("/notes/{id}", handler.GetNoteByID)
	return router, noteUsecase
}

func TestNoteHandler_GetNoteByID_InvalidID(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/notes/", nil) // Empty ID
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_GetNoteByID_NotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	req := httptest.NewRequest(http.MethodGet, "/notes/non-existent-id", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_GetNoteByID_Success(t *testing.T) {
	// Arrange
	router, noteUsecase := setupTest()
	noteID, err := noteUsecase.CreateNote("", "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/notes/"+noteID, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	var note usecase.NoteDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &note); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if note.ID != noteID {
		t.Errorf("expected note ID '%s'; got '%s'", noteID, note.ID)
	}
	if note.Title != "Test Title" {
		t.Errorf("expected note title 'Test Title'; got '%s'", note.Title)
	}
}

func TestNoteHandler_CreateNote_Success(t *testing.T) {
	// Arrange
	router, _ := setupTest()
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
	var responseBody struct {
		ID string `json:"id"`
	}
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
	router, _ := setupTest()
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
	router, _ := setupTest()
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
