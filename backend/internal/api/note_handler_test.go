package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/repository/noterepo"
	"noteapp/internal/usecase/noteuc"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// setupTest initializes the necessary components for the tests.
func setupTest() (*chi.Mux, *noteuc.NoteUsecase) {
	repo := noterepo.NewInMemoryNoteRepository()
	nc := noteuc.NewNoteUsecase(repo)
	handler := NewNoteHandler(nc)

	router := chi.NewRouter()
	router.Post("/notes", handler.CreateNote)
	router.Get("/notes/{id}", handler.GetNoteByID)
	router.Delete("/notes/{id}", handler.DeleteNote)
	router.Post("/notes/{id}/contents", handler.AddContent)
	router.Put("/notes/{id}/contents/{contentId}", handler.UpdateContent)
	router.Delete("/notes/{id}/contents/{contentId}", handler.DeleteContent)
	router.Post("/users/{userID}/notes/{noteID}/keyword", handler.TagNote)
	router.Get("/users/{userID}/notes", handler.FindNotesByKeyword)
	router.Delete("/users/{userID}/notes/{noteID}/keyword/{keyword}", handler.UntagNote)
	router.Post("/users/{ownerID}/notes/{noteID}/shares", handler.ShareNote)
	router.Delete("/users/{ownerID}/notes/{noteID}/shares", handler.RevokeAccess)
	router.Get("/users/{userID}/accessible-notes", handler.GetAccessibleNotesForUser)

	return router, nc
}

func TestNoteHandler_FindNotesByKeyword(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	note1, _ := nc.CreateNote("", "Note 1", "owner-1")
	note2, _ := nc.CreateNote("", "Note 2", "owner-1")
	note3, _ := nc.CreateNote("", "Note 3", "owner-1")

	nc.TagNote(note1, "user-1", "testing")
	nc.TagNote(note2, "user-1", "testing")
	nc.TagNote(note3, "user-2", "testing")

	req := httptest.NewRequest(http.MethodGet, "/users/user-1/notes?keyword=testing", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}
	var notes []*noteuc.NoteDTO
	if err := json.NewDecoder(rr.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
}

func TestNoteHandler_FindNotesByKeyword_NoMatch(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	note1, _ := nc.CreateNote("", "Note 1", "owner-1")
	note2, _ := nc.CreateNote("", "Note 2", "owner-1")
	note3, _ := nc.CreateNote("", "Note 3", "owner-1")

	nc.TagNote(note1, "user-1", "testing")
	nc.TagNote(note2, "user-1", "testing")
	nc.TagNote(note3, "user-2", "testing")

	req := httptest.NewRequest(http.MethodGet, "/users/user-1/notes?keyword=go", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}
	var notes []*noteuc.NoteDTO
	if err := json.NewDecoder(rr.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("expected 0 notes, got %d", len(notes))
	}
}

func TestNoteHandler_FindNotesByKeyword_EmptyKeyword(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	note1, _ := nc.CreateNote("", "Note 1", "owner-1")

	nc.TagNote(note1, "user-1", "testing")

	req := httptest.NewRequest(http.MethodGet, "/users/user-1/notes?keyword=", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}
	var notes []*noteuc.NoteDTO
	if err := json.NewDecoder(rr.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("expected 0 notes, got %d", len(notes))
	}
}

func TestNoteHandler_AddContent_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := AddContentRequest{
		Type: "text",
		Data: "Test content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notes/"+noteID+"/contents", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	var response struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response.ID == "" {
		t.Error("expected a non-empty content ID")
	}
}

func TestNoteHandler_AddContent_NotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	requestBody := AddContentRequest{
		Type: "text",
		Data: "Test content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notes/non-existent-id/contents", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_AddContent_InvalidJSON(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	invalidBody := []byte(`{"type": "text", "data":`) // Malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/notes/some-id/contents", bytes.NewBuffer(invalidBody))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_AddContent_UnsupportedContentType(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := AddContentRequest{
		Type: "unsupported",
		Data: "Test content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notes/"+noteID+"/contents", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_DeleteNote_InvalidID(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	req := httptest.NewRequest(http.MethodDelete, "/notes/", nil) // Empty ID
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_DeleteNote_NotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	req := httptest.NewRequest(http.MethodDelete, "/notes/non-existent-id", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_DeleteNote_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr.Code)
	}

	// Verify the note is actually deleted
	_, err = nc.GetNoteByID(noteID)
	if err != noteuc.ErrNoteNotFound {
		t.Errorf("expected ErrNoteNotFound after deletion, but got %v", err)
	}
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
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
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

	var note noteuc.NoteDTO
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
		Title: "Test Title",
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
	requestBody := CreateNoteRequest{Title: ""}
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

func TestNoteHandler_UpdateContent_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	contentID, err := nc.CreateAndAddContent(noteID, "", "Initial Content", "text")
	if err != nil {
		t.Fatalf("setup: failed to add content: %v", err)
	}

	requestBody := UpdateContentRequest{
		Data: "Updated Content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/"+noteID+"/contents/"+contentID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	if len(note.Contents) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(note.Contents))
	}
	if note.Contents[0].Data != "Updated Content" {
		t.Errorf("expected content to be 'Updated Content', got '%s'", note.Contents[0].Data)
	}
}

func TestNoteHandler_UpdateContent_NoteNotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	requestBody := UpdateContentRequest{
		Data: "Updated Content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/non-existent-id/contents/some-content-id", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_UpdateContent_ContentNotFound(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := UpdateContentRequest{
		Data: "Updated Content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/"+noteID+"/contents/non-existent-content-id", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_UpdateContent_InvalidJSON(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	invalidBody := []byte(`{"data":`) // Malformed JSON
	req := httptest.NewRequest(http.MethodPut, "/notes/some-id/contents/some-content-id", bytes.NewBuffer(invalidBody))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_DeleteContent_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	contentID, err := nc.CreateAndAddContent(noteID, "", "Initial Content", "text")
	if err != nil {
		t.Fatalf("setup: failed to add content: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID+"/contents/"+contentID, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr.Code)
	}

	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	if len(note.Contents) != 0 {
		t.Errorf("expected 0 content blocks, got %d", len(note.Contents))
	}
}

func TestNoteHandler_DeleteContent_NoteNotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	req := httptest.NewRequest(http.MethodDelete, "/notes/non-existent-id/contents/some-content-id", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_DeleteContent_ContentNotFound(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID+"/contents/non-existent-content-id", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_TagNote_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user-1"
	keyword := "test-keyword"
	requestBody := TagNoteRequest{Keyword: keyword}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+userID+"/notes/"+noteID+"/keyword", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}
	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	if len(note.Keywords[userID]) != 1 {
		t.Fatalf("expected 1 keyword, got %d", len(note.Keywords[userID]))
	}
	if note.Keywords[userID][0] != keyword {
		t.Errorf("expected keyword to be '%s', got '%s'", keyword, note.Keywords[userID][0])
	}
}

func TestNoteHandler_TagNote_NoteNotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	userID := "user-1"
	keyword := "test-keyword"
	requestBody := TagNoteRequest{Keyword: keyword}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+userID+"/notes/non-existent-id/keyword", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_TagNote_EmptyKeyword(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user-1"
	requestBody := TagNoteRequest{Keyword: ""}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+userID+"/notes/"+noteID+"/keyword", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_UntagNote_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID1 := "user-1"
	userID2 := "user-2"
	keyword := "test-keyword"
	nc.TagNote(noteID, userID1, keyword)
	nc.TagNote(noteID, userID2, keyword)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+userID1+"/notes/"+noteID+"/keyword/"+keyword, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr.Code)
	}
	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	if len(note.Keywords[userID1]) != 0 {
		t.Errorf("expected 0 keywords for user1, got %d", len(note.Keywords[userID1]))
	}
	if len(note.Keywords[userID2]) != 1 {
		t.Errorf("expected 1 keyword for user2, got %d", len(note.Keywords[userID2]))
	}
}

func TestNoteHandler_UntagNote_NoteNotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	req := httptest.NewRequest(http.MethodDelete, "/users/user-1/notes/non-existent-id/keyword/test-keyword", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_UntagNote_UserNotFound(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user-1"
	keyword := "test-keyword"
	nc.TagNote(noteID, userID, keyword)

	req := httptest.NewRequest(http.MethodDelete, "/users/non-existent-user/notes/"+noteID+"/keyword/"+keyword, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_UntagNote_KeywordNotFound(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user-1"
	keyword := "test-keyword"
	nc.TagNote(noteID, userID, keyword)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+userID+"/notes/"+noteID+"/keyword/non-existent-keyword", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_ShareNote_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user-2"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"user_id":    collaboratorID,
		"permission": "read",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+ownerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}
	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	perm, exists := note.Collaborators[collaboratorID]
	if !exists {
		t.Fatalf("expected collaborator '%s' to be in shares", collaboratorID)
	}
	if perm != "read" {
		t.Errorf("expected permission to be 'read'; got '%s'", perm)
	}
}

func TestNoteHandler_ShareNote_Unauthorized(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	nonOwnerID := "non-owner"
	collaboratorID := "user-2"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"user_id":    collaboratorID,
		"permission": "read",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+nonOwnerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d; got %d", http.StatusForbidden, rr.Code)
	}
}

func TestNoteHandler_ShareNote_NoteNotFound(t *testing.T) {
	// Arrange
	router, _ := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user-2"
	nonExistentNoteID := "non-existent-note-id"

	requestBody := map[string]interface{}{
		"user_id":    collaboratorID,
		"permission": "read",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+ownerID+"/notes/"+nonExistentNoteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_ShareNote_InvalidPermission(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user-2"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"user_id":    collaboratorID,
		"permission": "invalid-permission",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/users/"+ownerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_ShareNote_UpdatePermission(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user-2"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	// First, share with "read" permission
	initialRequestBody := map[string]interface{}{
		"user_id":    collaboratorID,
		"permission": "read",
	}
	initialBody, _ := json.Marshal(initialRequestBody)
	initialReq := httptest.NewRequest(http.MethodPost, "/users/"+ownerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(initialBody))
	initialRR := httptest.NewRecorder()
	router.ServeHTTP(initialRR, initialReq)
	if initialRR.Code != http.StatusCreated {
		t.Fatalf("initial share failed: expected status %d; got %d", http.StatusCreated, initialRR.Code)
	}

	// Now, update the permission to "read-write"
	updateRequestBody := map[string]interface{}{
		"user_id":    collaboratorID,
		"permission": "read-write",
	}
	updateBody, _ := json.Marshal(updateRequestBody)
	updateReq := httptest.NewRequest(http.MethodPost, "/users/"+ownerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(updateBody))
	updateRR := httptest.NewRecorder()

	// Act
	router.ServeHTTP(updateRR, updateReq)

	// Assert
	if updateRR.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, updateRR.Code)
	}

	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	perm, exists := note.Collaborators[collaboratorID]
	if !exists {
		t.Fatalf("expected collaborator '%s' to be in shares", collaboratorID)
	}
	if perm != "read-write" {
		t.Errorf("expected permission to be 'read-write'; got '%s'", perm)
	}
}

func TestNoteHandler_GetAccessibleNotesForUser(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	userID := "user-1"
	ownerID := "owner-1"

	// Create notes
	_, err := nc.CreateNote("", "Owned Note", userID)
	if err != nil {
		t.Fatalf("failed to create owned note: %v", err)
	}
	sharedNoteID, err := nc.CreateNote("", "Shared Note", ownerID)
	if err != nil {
		t.Fatalf("failed to create shared note: %v", err)
	}
	_, err = nc.CreateNote("", "Unrelated Note", ownerID)
	if err != nil {
		t.Fatalf("failed to create unrelated note: %v", err)
	}

	// Share the note
	err = nc.ShareNote(sharedNoteID, ownerID, userID, "read")
	if err != nil {
		t.Fatalf("failed to share note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID+"/accessible-notes", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	var notes []noteuc.NoteDTO
	if err := json.NewDecoder(rr.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(notes) != 2 {
		t.Fatalf("expected 2 accessible notes, but got %d", len(notes))
	}
}

func TestNoteHandler_RevokeAccess_Success(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	collaboratorID1 := "user-1"
	collaboratorID2 := "user-2"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID1, "read")
	nc.ShareNote(noteID, ownerID, collaboratorID2, "read")
	nc.TagNote(noteID, collaboratorID1, "test-keyword-1")
	nc.TagNote(noteID, collaboratorID2, "test-keyword-2")

	body, _ := json.Marshal(map[string]string{
		"user_id": collaboratorID1,
	})
	req := httptest.NewRequest(http.MethodDelete, "/users/"+ownerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr.Code)
	}
	note, _ := nc.GetNoteByID(noteID)
	if _, ok := note.Collaborators[collaboratorID1]; ok {
		t.Errorf("Expected collaborator 1 to be removed, but they still exist")
	}
	if _, ok := note.Collaborators[collaboratorID2]; !ok {
		t.Errorf("Expected collaborator 2 to remain, but they were removed")
	}
	if _, ok := note.Keywords[collaboratorID1]; ok {
		t.Errorf("Expected collaborator 1's keywords to be removed, but they still exist")
	}
	if _, ok := note.Keywords[collaboratorID2]; !ok {
		t.Errorf("Expected collaborator 2's keywords to remain, but they were removed")
	}
}

func TestNoteHandler_RevokeAccess_NotOwner(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user-1"
	nonOwnerID := "user-2"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID, "read")

	body, _ := json.Marshal(map[string]string{"user_id": collaboratorID})
	req := httptest.NewRequest(http.MethodDelete, "/users/"+nonOwnerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d; got %d", http.StatusForbidden, rr.Code)
	}
}

func TestNoteHandler_RevokeAccess_CollaboratorNotFound(t *testing.T) {
	// Arrange
	router, nc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user-1"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID, "read")

	body, _ := json.Marshal(map[string]string{"user_id": "non-existent-user"})
	req := httptest.NewRequest(http.MethodDelete, "/users/"+ownerID+"/notes/"+noteID+"/shares", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}
