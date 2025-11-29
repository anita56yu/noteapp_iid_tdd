package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/middleware"
	"noteapp/internal/repository/contentrepo"
	"noteapp/internal/repository/noterepo"
	"noteapp/internal/repository/userrepo"
	"noteapp/internal/usecase/contentuc"
	"noteapp/internal/usecase/noteuc"
	"noteapp/internal/usecase/useruc"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// setupTest initializes the necessary components for the tests.
func setupTest() (*chi.Mux, *noteuc.NoteUsecase, *contentuc.ContentUsecase, *useruc.UserUsecase) {
	noteRepo := noterepo.NewInMemoryNoteRepository()
	contentRepo := contentrepo.NewInMemoryContentRepository()
	userRepo := userrepo.NewInMemoryUserRepository()
	nuc := noteuc.NewNoteUsecase(noteRepo)
	cuc := contentuc.NewContentUsecase(contentRepo)
	uuc := useruc.NewUserUsecase(userRepo, &useruc.UserMapper{})
	uuc.Register("owner-1", "owner-1", "password")
	uuc.Register("user2id", "user-2", "password")
	uuc.Register("user1id", "user-1", "password")
	handler := NewNoteHandler(nuc, cuc, uuc)
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(middleware.MockAuthMiddleware())
		r.Post("/notes", handler.CreateNote)
		r.Post("/notes/{noteID}/keyword", handler.TagNote)
		r.Get("/notes", handler.FindNotesByKeyword)
		r.Delete("/notes/{noteID}/keyword/{keyword}", handler.UntagNote)
		r.Post("/notes/{noteID}/shares", handler.ShareNote)
		r.Delete("/notes/{noteID}/shares", handler.RevokeAccess)
		r.Get("/notes/accessible-notes", handler.GetAccessibleNotesForUser)
	})
	router.Get("/notes/{id}", handler.GetNoteByID)
	router.Put("/notes/{id}", handler.UpdateNote)
	router.Delete("/notes/{id}", handler.DeleteNote)
	router.Post("/notes/{id}/contents", handler.AddContent)
	router.Put("/notes/{id}/contents/{contentId}", handler.UpdateContent)
	router.Delete("/notes/{id}/contents/{contentId}", handler.DeleteContent)

	router.Get("/ws/notes/{noteID}", handler.HandleWebSocket)
	return router, nuc, cuc, uuc
}

func TestNoteHandler_FindNotesByKeyword(t *testing.T) {
	// Arrange
	router, nc, _, _ := setupTest()
	note1, _ := nc.CreateNote("", "Note 1", "owner-1")
	note2, _ := nc.CreateNote("", "Note 2", "owner-1")
	note3, _ := nc.CreateNote("", "Note 3", "owner-1")

	nc.TagNote(note1, "user-1", "testing", 0)
	nc.TagNote(note2, "user-1", "testing", 0)
	nc.TagNote(note3, "user-2", "testing", 0)

	req := MockAuthenticatedRequestForTest(http.MethodGet, "/notes?keyword=testing", "user-1", nil)
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
	router, nc, _, _ := setupTest()
	note1, _ := nc.CreateNote("", "Note 1", "owner-1")
	note2, _ := nc.CreateNote("", "Note 2", "owner-1")
	note3, _ := nc.CreateNote("", "Note 3", "owner-1")

	nc.TagNote(note1, "user-1", "testing", 0)
	nc.TagNote(note2, "user-1", "testing", 0)
	nc.TagNote(note3, "user-2", "testing", 0)

	req := MockAuthenticatedRequestForTest(http.MethodGet, "/notes?keyword=go", "user-1", nil)
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
	router, nc, _, _ := setupTest()
	note1, _ := nc.CreateNote("", "Note 1", "owner-1")

	nc.TagNote(note1, "user-1", "testing", 0)

	req := MockAuthenticatedRequestForTest(http.MethodGet, "/notes?keyword=", "user-1", nil)
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
	router, nuc, cuc, _ := setupTest()
	noteID, err := nuc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := AddContentRequest{
		Type:        "text",
		Data:        "Test content",
		NoteVersion: intPtr(0),
		Index:       intPtr(-1),
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

	// Verify that the content ID was added to the note.
	note, err := nuc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note by ID: %v", err)
	}

	if len(note.ContentIDs) != 1 {
		t.Fatalf("expected note to have 1 content ID, but got %d", len(note.ContentIDs))
	}

	if note.ContentIDs[0] != response.ID {
		t.Errorf("expected content ID '%s' to be in note's content IDs, but it was not", response.ID)
	}

	// Verify that the content was created with the correct data.
	content, err := cuc.GetContentByID(response.ID)
	if err != nil {
		t.Fatalf("failed to get content by ID: %v", err)
	}
	if content.Data != requestBody.Data {
		t.Errorf("expected content data to be '%s', but got '%s'", requestBody.Data, content.Data)
	}
	if content.Type != string(contentuc.TextContentType) {
		t.Errorf("expected content type to be '%s', but got '%s'", contentuc.TextContentType, content.Type)
	}
}

func TestNoteHandler_AddContent_NotFound(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
	requestBody := AddContentRequest{
		Type:        "text",
		Data:        "Test content",
		NoteVersion: intPtr(0),
		Index:       intPtr(-1),
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
	router, _, _, _ := setupTest()
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := AddContentRequest{
		Type:        "unsupported",
		Data:        "Test content",
		Index:       intPtr(-1),
		NoteVersion: intPtr(0),
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

func TestNoteHandler_AddContent_OutOfBoundsIndex(t *testing.T) {
	// Arrange
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	requestBody := AddContentRequest{
		Type:        "text",
		Data:        "Test content",
		Index:       intPtr(5), // Out of bounds
		NoteVersion: intPtr(0),
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
	router, _, _, _ := setupTest()
	requestBody := DeleteNoteRequest{
		NoteVersion: intPtr(0),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/notes/", bytes.NewBuffer(body))
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
	router, _, _, _ := setupTest()
	requestBody := DeleteNoteRequest{
		NoteVersion: intPtr(0),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/notes/non-existent-id", bytes.NewBuffer(body))
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
	router, nc, _, uuc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	uuc.AddAccessibleNote("owner-1", noteID)

	requestBody := DeleteNoteRequest{
		NoteVersion: intPtr(0),
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID, bytes.NewBuffer(body))
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
	ownerDTO, _ := uuc.Login("owner-1", "password")
	if len(ownerDTO.AccessibleNoteIDs) != 0 {
		t.Errorf("expected owner to have 0 accessible notes after deletion, but got %d", len(ownerDTO.AccessibleNoteIDs))
	}
}

func TestNoteHandler_DeleteNoteWithCollaborator_Success(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	uuc.AddAccessibleNote("owner-1", noteID)
	uuc.AddAccessibleNote("user2id", noteID)
	err = nc.ShareNote(noteID, "owner-1", "user2id", "read", 0)
	if err != nil {
		t.Fatalf("setup: failed to share note: %v", err)
	}

	requestBody1 := DeleteNoteRequest{
		NoteVersion: intPtr(1),
	}
	body1, _ := json.Marshal(requestBody1)

	req1 := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID, bytes.NewBuffer(body1))
	rr1 := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr1, req1)

	// Assert
	if rr1.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr1.Code)
	}

	// Verify the note is actually deleted
	_, err = nc.GetNoteByID(noteID)
	if err != noteuc.ErrNoteNotFound {
		t.Errorf("expected ErrNoteNotFound after deletion, but got %v", err)
	}
	ownerDTO, _ := uuc.Login("owner-1", "password")
	if len(ownerDTO.AccessibleNoteIDs) != 0 {
		t.Errorf("expected owner to have 0 accessible notes after deletion, but got %d", len(ownerDTO.AccessibleNoteIDs))
	}
	collabDTO, _ := uuc.Login("user-2", "password")
	if len(collabDTO.AccessibleNoteIDs) != 0 {
		t.Errorf("expected collaborator to have 0 accessible notes after deletion, but got %d", len(collabDTO.AccessibleNoteIDs))
	}
}

func TestNoteHandler_GetNoteByID_InvalidID(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
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
	router, _, _, _ := setupTest()
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
	router, nuc, cuc, _ := setupTest()
	noteID, err := nuc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	// Add content to the note
	contentData1 := "First content"
	contentID1, err := cuc.CreateContent(noteID, "", contentData1, contentuc.TextContentType)
	if err != nil {
		t.Fatalf("setup: failed to create content 1: %v", err)
	}
	if err := nuc.AddContent(noteID, contentID1, -1, 0); err != nil {
		t.Fatalf("setup: failed to add content 1 to note: %v", err)
	}

	contentData2 := "Second content"
	contentID2, err := cuc.CreateContent(noteID, "", contentData2, contentuc.TextContentType)
	if err != nil {
		t.Fatalf("setup: failed to create content 2: %v", err)
	}
	if err := nuc.AddContent(noteID, contentID2, -1, 1); err != nil {
		t.Fatalf("setup: failed to add content 2 to note: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/notes/"+noteID, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	var response GetNoteByIDResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if response.ID != noteID {
		t.Errorf("expected note ID '%s'; got '%s'", noteID, response.ID)
	}
	if response.Title != "Test Title" {
		t.Errorf("expected note title 'Test Title'; got '%s'", response.Title)
	}

	if len(response.Contents) != 2 {
		t.Fatalf("expected 2 contents, got %d", len(response.Contents))
	}

	// Verify content 1
	if response.Contents[0].ID != contentID1 {
		t.Errorf("expected first content ID '%s'; got '%s'", contentID1, response.Contents[0].ID)
	}
	if response.Contents[0].Data != contentData1 {
		t.Errorf("expected first content data '%s'; got '%s'", contentData1, response.Contents[0].Data)
	}

	// Verify content 2
	if response.Contents[1].ID != contentID2 {
		t.Errorf("expected second content ID '%s'; got '%s'", contentID2, response.Contents[1].ID)
	}
	if response.Contents[1].Data != contentData2 {
		t.Errorf("expected second content data '%s'; got '%s'", contentData2, response.Contents[1].Data)
	}
}

func TestNoteHandler_CreateNote_Success(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
	ownerID := "owner-1"
	requestBody := CreateNoteRequest{
		Title: "Test Title",
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes", ownerID, bytes.NewBuffer(body))
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

func TestNoteHandler_CreateNote_UserNotFound(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
	nonExistentUserID := "non-existent-user"
	requestBody := CreateNoteRequest{
		Title: "Test Title",
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes", nonExistentUserID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_CreateNote_InvalidJSON(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
	ownerID := "owner-1"
	invalidBody := []byte(`{"title": "Test", "content":`) // Malformed JSON
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes", ownerID, bytes.NewBuffer(invalidBody))
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
	router, _, _, _ := setupTest()
	requestBody := CreateNoteRequest{Title: ""}
	ownerID := "owner-1"
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes", ownerID, bytes.NewBuffer(body))
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

func intPtr(i int) *int {
	return &i
}

func TestNoteHandler_UpdateContent_Success(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTest()
	noteID, err := nuc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	contentID, err := cuc.CreateContent(noteID, "", "Initial Content", contentuc.TextContentType)
	if err != nil {
		t.Fatalf("setup: failed to create content: %v", err)
	}
	if err := nuc.AddContent(noteID, contentID, -1, 0); err != nil {
		t.Fatalf("setup: failed to add content to note: %v", err)
	}

	requestBody := UpdateContentRequest{
		Data:           "Updated Content",
		ContentVersion: intPtr(0),
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

	// Verify that the content was updated.
	content, err := cuc.GetContentByID(contentID)
	if err != nil {
		t.Fatalf("failed to get content by ID: %v", err)
	}
	if content.Data != requestBody.Data {
		t.Errorf("expected content data to be '%s', but got '%s'", requestBody.Data, content.Data)
	}
	if content.Version != 1 {
		t.Errorf("expected content version to be 1, but got %d", content.Version)
	}
}

func TestNoteHandler_UpdateContent_NoteNotFound(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
	requestBody := UpdateContentRequest{
		Data:           "Updated Content",
		ContentVersion: intPtr(0),
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := UpdateContentRequest{
		Data:           "Updated Content",
		ContentVersion: intPtr(0),
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
	router, _, _, _ := setupTest()
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

func TestNoteHandler_UpdateContent_MissingVersion(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTest()
	noteID, _ := nuc.CreateNote("", "Test Title", "owner-1")
	contentID, _ := cuc.CreateContent(noteID, "", "Initial Content", contentuc.TextContentType)
	nuc.AddContent(noteID, contentID, -1, 0)

	// Request body without version
	requestBody := map[string]string{
		"data": "Updated Content",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/"+noteID+"/contents/"+contentID, bytes.NewBuffer(body))
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
	router, nuc, cuc, _ := setupTest()
	noteID, err := nuc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	contentID, err := cuc.CreateContent(noteID, "", "Initial Content", contentuc.TextContentType)
	if err != nil {
		t.Fatalf("setup: failed to create content: %v", err)
	}
	if err := nuc.AddContent(noteID, contentID, -1, 0); err != nil {
		t.Fatalf("setup: failed to add content to note: %v", err)
	}

	deleteReq := DeleteContentRequest{
		ContentVersion: intPtr(0),
		NoteVersion:    intPtr(1),
	}
	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID+"/contents/"+contentID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr.Code)
	}

	// Verify that the content ID was removed from the note.
	note, err := nuc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	if len(note.ContentIDs) != 0 {
		t.Errorf("expected 0 content blocks, got %d", len(note.ContentIDs))
	}

	// Verify that the content was deleted from the repository.
	_, err = cuc.GetContentByID(contentID)
	if err != contentuc.ErrContentNotFound {
		t.Errorf("expected ErrContentNotFound, but got %v", err)
	}
}

func TestNoteHandler_DeleteContent_NoteNotFound(t *testing.T) {
	// Arrange
	router, _, _, _ := setupTest()
	deleteReq := DeleteContentRequest{
		ContentVersion: intPtr(0),
		NoteVersion:    intPtr(0),
	}
	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest(http.MethodDelete, "/notes/non-existent-id/contents/some-content-id", bytes.NewBuffer(body))
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	deleteReq := DeleteContentRequest{
		ContentVersion: intPtr(0),
		NoteVersion:    intPtr(0),
	}
	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID+"/contents/non-existent-content-id", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_DeleteContent_MissingVersion(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTest()
	noteID, _ := nuc.CreateNote("", "Test Title", "owner-1")
	contentID, _ := cuc.CreateContent(noteID, "", "Initial Content", contentuc.TextContentType)
	nuc.AddContent(noteID, contentID, -1, 0)

	// Empty body
	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID+"/contents/"+contentID, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_TagNote_Success(t *testing.T) {
	// Arrange
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user1id"
	keyword := "test-keyword"
	requestBody := TagNoteRequest{Keyword: keyword, NoteVersion: intPtr(0)}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/keyword", userID, bytes.NewBuffer(body))
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
	router, _, _, _ := setupTest()
	userID := "user-1"
	keyword := "test-keyword"
	requestBody := TagNoteRequest{Keyword: keyword, NoteVersion: intPtr(0)}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/non-existent-id/keyword", userID, bytes.NewBuffer(body))
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "userid1"
	requestBody := TagNoteRequest{Keyword: "", NoteVersion: intPtr(0)}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/keyword", userID, bytes.NewBuffer(body))
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID1 := "user-1"
	userID2 := "user-2"
	keyword := "test-keyword"
	nc.TagNote(noteID, userID1, keyword, 0)
	nc.TagNote(noteID, userID2, keyword, 1)
	requestBody := UntagNoteRequest{NoteVersion: intPtr(2)}
	body, _ := json.Marshal(requestBody)

	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/keyword/"+keyword, userID1, bytes.NewBuffer(body))
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
	router, _, _, _ := setupTest()
	requestBody := UntagNoteRequest{NoteVersion: intPtr(0)}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/non-existent-id/keyword/test-keyword", "user-1", bytes.NewBuffer(body))
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user-1"
	keyword := "test-keyword"
	nc.TagNote(noteID, userID, keyword, 0)
	requestBody := UntagNoteRequest{NoteVersion: intPtr(1)}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/keyword/"+keyword, "non-existent-user", bytes.NewBuffer(body))
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
	router, nc, _, _ := setupTest()
	noteID, err := nc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	userID := "user-1"
	keyword := "test-keyword"
	nc.TagNote(noteID, userID, keyword, 0)
	requestBody := UntagNoteRequest{NoteVersion: intPtr(1)}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/keyword/non-existent-keyword", userID, bytes.NewBuffer(body))
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
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorUsername := "user-2"
	collaboratorID := "user2id"

	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"username":     collaboratorUsername,
		"permission":   "read",
		"note_version": 0,
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(body))
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
	userDTO, err := uuc.Login(collaboratorUsername, "password")
	if err != nil {
		t.Fatalf("failed to login as collaborator: %v", err)
	}
	if len(userDTO.AccessibleNoteIDs) != 1 || userDTO.AccessibleNoteIDs[0] != noteID {
		t.Errorf("expected shared notes to contain note ID '%s'; got %v", noteID, userDTO.AccessibleNoteIDs)
	}
}

func TestNoteHandler_ShareNote_Unauthorized(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	nonOwnerID := "non-owner"
	collaboratorUsername := "user-2"

	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"username":     collaboratorUsername,
		"permission":   "read",
		"note_version": 0,
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/shares", nonOwnerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d; got %d", http.StatusForbidden, rr.Code)
	}
	userDTO, err := uuc.Login(collaboratorUsername, "password")
	if err != nil {
		t.Fatalf("failed to login as collaborator: %v", err)
	}
	if len(userDTO.AccessibleNoteIDs) != 0 {
		t.Errorf("expected no accessible notes; got %v", userDTO.AccessibleNoteIDs)
	}
}

func TestNoteHandler_ShareNote_NoteNotFound(t *testing.T) {
	// Arrange
	router, _, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorUsername := "user-2"
	nonExistentNoteID := "non-existent-note-id"

	requestBody := map[string]interface{}{
		"username":     collaboratorUsername,
		"permission":   "read",
		"note_version": 0,
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+nonExistentNoteID+"/shares", ownerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}

	userDTO, err := uuc.Login(collaboratorUsername, "password")
	if err != nil {
		t.Fatalf("failed to login as collaborator: %v", err)
	}
	if len(userDTO.AccessibleNoteIDs) != 0 {
		t.Errorf("expected no accessible notes; got %v", userDTO.AccessibleNoteIDs)
	}
}

func TestNoteHandler_ShareNote_InvalidPermission(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorUsername := "user-2"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"username":     collaboratorUsername,
		"permission":   "invalid-permission",
		"note_version": 0,
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
	userDTO, err := uuc.Login(collaboratorUsername, "password")
	if err != nil {
		t.Fatalf("failed to login as collaborator: %v", err)
	}
	if len(userDTO.AccessibleNoteIDs) != 0 {
		t.Errorf("expected no accessible notes; got %v", userDTO.AccessibleNoteIDs)
	}
}

func TestNoteHandler_ShareNote_UpdatePermission(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorUsername := "user-2"
	collaboratorID := "user2id"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	// First, share with "read" permission
	initialRequestBody := map[string]interface{}{
		"username":     collaboratorUsername,
		"permission":   "read",
		"note_version": 0,
	}
	initialBody, _ := json.Marshal(initialRequestBody)
	initialReq := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(initialBody))
	initialRR := httptest.NewRecorder()
	router.ServeHTTP(initialRR, initialReq)
	if initialRR.Code != http.StatusCreated {
		t.Fatalf("initial share failed: expected status %d; got %d", http.StatusCreated, initialRR.Code)
	}

	// Now, update the permission to "read-write"
	updateRequestBody := map[string]interface{}{
		"username":     collaboratorUsername,
		"permission":   "read-write",
		"note_version": 1,
	}
	updateBody, _ := json.Marshal(updateRequestBody)
	updateReq := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(updateBody))
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
	userDTO, err := uuc.Login(collaboratorUsername, "password")
	if err != nil {
		t.Fatalf("failed to login as collaborator: %v", err)
	}
	if len(userDTO.AccessibleNoteIDs) != 1 || userDTO.AccessibleNoteIDs[0] != noteID {
		t.Errorf("expected shared notes to contain note ID '%s'; got %v", noteID, userDTO.AccessibleNoteIDs)
	}
}

func TestNoteHandler_ShareNote_CollaboratorNotFound(t *testing.T) {
	// Arrange
	router, nc, _, _ := setupTest()
	ownerID := "owner-1"
	nonExistentUsername := "non-existent-user"
	noteID, err := nc.CreateNote("", "Test Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	requestBody := map[string]interface{}{
		"username":     nonExistentUsername,
		"permission":   "read",
		"note_version": 0,
	}
	body, _ := json.Marshal(requestBody)
	req := MockAuthenticatedRequestForTest(http.MethodPost, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
	note, err := nc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get note: %v", err)
	}
	if len(note.Collaborators) != 0 {
		t.Errorf("expected no collaborators; got %v", note.Collaborators)
	}
}

func TestNoteHandler_GetAccessibleNotesForUser(t *testing.T) {
	// Arrange
	router, nc, _, _ := setupTest()
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
	err = nc.ShareNote(sharedNoteID, ownerID, userID, "read", 0)
	if err != nil {
		t.Fatalf("failed to share note: %v", err)
	}

	req := MockAuthenticatedRequestForTest(http.MethodGet, "/notes/accessible-notes", userID, nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	var notes []GetNoteByIDResponse
	if err := json.NewDecoder(rr.Body).Decode(&notes); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(notes) != 2 {
		t.Fatalf("expected 2 accessible notes, but got %d", len(notes))
	}
}

func TestNoteHandler_RevokeAccess_Success(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorUsername1 := "user-1"
	collaboratorID1 := "user1id"
	collaboratorID2 := "user2id"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID1, "read", 0)
	nc.ShareNote(noteID, ownerID, collaboratorID2, "read", 1)
	uuc.AddAccessibleNote(ownerID, noteID)
	uuc.AddAccessibleNote(collaboratorID1, noteID)
	uuc.AddAccessibleNote(collaboratorID2, noteID)
	nc.TagNote(noteID, collaboratorID1, "test-keyword-1", 2)
	nc.TagNote(noteID, collaboratorID2, "test-keyword-2", 3)

	body, _ := json.Marshal(map[string]interface{}{
		"username":     collaboratorUsername1,
		"note_version": intPtr(4),
	})
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(body))
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
	user1DTO, _ := uuc.Login(collaboratorUsername1, "password")
	for _, accessibleNoteID := range user1DTO.AccessibleNoteIDs {
		if accessibleNoteID == noteID {
			t.Errorf("Expected collaborator 1 to no longer have access to the note, but they do")
		}
	}
	user2DTO, _ := uuc.Login("user-2", "password")
	found := false
	for _, accessibleNoteID := range user2DTO.AccessibleNoteIDs {
		if accessibleNoteID == noteID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected collaborator 2 to still have access to the note, but they do not")
	}
	ownerDTO, _ := uuc.Login(ownerID, "password")
	found = false
	for _, accessibleNoteID := range ownerDTO.AccessibleNoteIDs {
		if accessibleNoteID == noteID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected owner to still have access to the note, but they do not")
	}
}

func TestNoteHandler_RevokeAccess_NotOwner(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorUsername := "user-1"
	collaboratorID := "user1id"
	nonOwnerID := "user-2"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID, "read", 0)
	uuc.AddAccessibleNote(ownerID, noteID)
	uuc.AddAccessibleNote(collaboratorID, noteID)

	body, _ := json.Marshal(map[string]interface{}{"username": collaboratorUsername, "note_version": intPtr(1)})
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/shares", nonOwnerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d; got %d", http.StatusForbidden, rr.Code)
	}
	ownerDTO, _ := uuc.Login(ownerID, "password")
	found := false
	for _, accessibleNoteID := range ownerDTO.AccessibleNoteIDs {
		if accessibleNoteID == noteID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected owner to still have access to the note, but they do not")
	}
	collaboratorDTO, _ := uuc.Login(collaboratorUsername, "password")
	found = false
	for _, accessibleNoteID := range collaboratorDTO.AccessibleNoteIDs {
		if accessibleNoteID == noteID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected collaborator to still have access to the note, but they do not")
	}
}

func TestNoteHandler_RevokeAccess_UserNotFound(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user1id"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID, "read", 0)
	uuc.AddAccessibleNote(ownerID, noteID)
	uuc.AddAccessibleNote(collaboratorID, noteID)

	body, _ := json.Marshal(map[string]interface{}{"username": "non-existent-user", "note_version": intPtr(1)})
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNoteHandler_RevokeAccess_CollaboratorNotFound(t *testing.T) {
	// Arrange
	router, nc, _, uuc := setupTest()
	ownerID := "owner-1"
	collaboratorID := "user1id"
	noncollaboratorUsername := "user-2"
	noteID, _ := nc.CreateNote("", "Test Note", ownerID)
	nc.ShareNote(noteID, ownerID, collaboratorID, "read", 0)
	uuc.AddAccessibleNote(ownerID, noteID)
	uuc.AddAccessibleNote(collaboratorID, noteID)

	body, _ := json.Marshal(map[string]interface{}{"username": noncollaboratorUsername, "note_version": intPtr(1)})
	req := MockAuthenticatedRequestForTest(http.MethodDelete, "/notes/"+noteID+"/shares", ownerID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNoteHandler_DeleteNote_DeletesAssociatedContent(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTest()
	noteID, err := nuc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}
	contentID1, err := cuc.CreateContent(noteID, "", "Content 1", contentuc.TextContentType)
	if err != nil {
		t.Fatalf("setup: failed to create content 1: %v", err)
	}
	contentID2, err := cuc.CreateContent(noteID, "", "Content 2", contentuc.TextContentType)
	if err != nil {
		t.Fatalf("setup: failed to create content 2: %v", err)
	}
	if err := nuc.AddContent(noteID, contentID1, -1, 0); err != nil {
		t.Fatalf("setup: failed to add content 1 to note: %v", err)
	}
	if err := nuc.AddContent(noteID, contentID2, -1, 1); err != nil {
		t.Fatalf("setup: failed to add content 2 to note: %v", err)
	}

	requestBody := DeleteNoteRequest{
		NoteVersion: intPtr(2),
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d; got %d", http.StatusNoContent, rr.Code)
	}

	// Verify the note is deleted
	_, err = nuc.GetNoteByID(noteID)
	if !errors.Is(err, noteuc.ErrNoteNotFound) {
		t.Errorf("expected ErrNoteNotFound after deletion, but got %v", err)
	}

	// Verify the associated contents are deleted
	_, err = cuc.GetContentByID(contentID1)
	if !errors.Is(err, contentuc.ErrContentNotFound) {
		t.Errorf("expected content 1 to be deleted, but got error %v", err)
	}
	_, err = cuc.GetContentByID(contentID2)
	if !errors.Is(err, contentuc.ErrContentNotFound) {
		t.Errorf("expected content 2 to be deleted, but got error %v", err)
	}
}

func TestNoteHandler_ChangeTitle_Success(t *testing.T) {
	// Arrange
	router, nuc, _, _ := setupTest()
	ownerID := "owner-1"
	noteID, err := nuc.CreateNote("", "Original Title", ownerID)
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	newTitle := "New and Improved Title"
	requestBody := UpdateNoteRequest{
		Title:       newTitle,
		NoteVersion: intPtr(0),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	// Verify the note's title was updated
	updatedNote, err := nuc.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("failed to get updated note: %v", err)
	}
	if updatedNote.Title != newTitle {
		t.Errorf("expected note title '%s'; got '%s'", newTitle, updatedNote.Title)
	}
	if updatedNote.Version != 1 {
		t.Errorf("expected note version 1; got %d", updatedNote.Version)
	}
}
