package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/repository/contentrepo"
	"noteapp/internal/repository/noterepo"
	"noteapp/internal/usecase/contentuc"
	"noteapp/internal/usecase/noteuc"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func setupTestForBroadcast() (*chi.Mux, *noteuc.NoteUsecase, *contentuc.ContentUsecase, *ConnectionManager) {
	noteRepo := noterepo.NewInMemoryNoteRepository()
	contentRepo := contentrepo.NewInMemoryContentRepository()
	nuc := noteuc.NewNoteUsecase(noteRepo)
	cuc := contentuc.NewContentUsecase(contentRepo)
	handler := NewNoteHandler(nuc, cuc)

	router := chi.NewRouter()
	router.Post("/notes", handler.CreateNote)
	router.Get("/notes/{id}", handler.GetNoteByID)
	router.Delete("/notes/{id}", handler.DeleteNote)
	router.Put("/notes/{id}", handler.UpdateNote)
	router.Post("/notes/{id}/contents", handler.AddContent)
	router.Put("/notes/{id}/contents/{contentId}", handler.UpdateContent)
	router.Delete("/notes/{id}/contents/{contentId}", handler.DeleteContent)
	router.Post("/users/{userID}/notes/{noteID}/keyword", handler.TagNote)
	router.Get("/users/{userID}/notes", handler.FindNotesByKeyword)
	router.Delete("/users/{userID}/notes/{noteID}/keyword/{keyword}", handler.UntagNote)
	router.Post("/users/{ownerID}/notes/{noteID}/shares", handler.ShareNote)
	router.Delete("/users/{ownerID}/notes/{noteID}/shares", handler.RevokeAccess)
	router.Get("/users/{userID}/accessible-notes", handler.GetAccessibleNotesForUser)

	router.Get("/ws/notes/{noteID}", handler.HandleWebSocket)
	return router, nuc, cuc, handler.connManager
}

func setUpNoteWithContents(nuc *noteuc.NoteUsecase, cuc *contentuc.ContentUsecase) (string, string, string) {
	noteID, err := nuc.CreateNote("", "Test Title", "owner-1")
	if err != nil {
		fmt.Printf("setup: failed to create note: %v", err)
	}

	contentID, err := cuc.CreateContent(noteID, "", "Test content", "text")
	if err != nil {
		fmt.Printf("setup: failed to create content: %v", err)
	}
	err = nuc.AddContent(noteID, contentID, -1, 0)
	if err != nil {
		fmt.Printf("setup: failed to add content to note: %v", err)
	}

	contentID1, err := cuc.CreateContent(noteID, "", "Test content", "text")
	if err != nil {
		fmt.Printf("setup: failed to create content: %v", err)
	}
	err = nuc.AddContent(noteID, contentID1, -1, 1)
	if err != nil {
		fmt.Printf("setup: failed to add content to note: %v", err)
	}

	return noteID, contentID, contentID1
}

func TestNoteHandler_WebSocket_BroadcastOnUpdate(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTestForBroadcast()
	noteID, _, _ := setUpNoteWithContents(nuc, cuc)

	// Act
	// In a separate goroutine, listen for a message on the WebSocket connection.
	msgChan := make(chan []byte)
	errChan := make(chan error)
	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/notes/" + noteID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		errChan <- fmt.Errorf("failed to dial websocket: %v", err)
		return
	}
	defer conn.Close()
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Failed to read message: %v", err)
			close(msgChan)
			return
		}
		msgChan <- msg
	}()
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Failed to read message: %v", err)
			close(msgChan)
			return
		}
		msgChan <- msg
	}()

	// Update the note by adding content.
	requestBody := AddContentRequest{
		Type:        "text",
		Data:        "Test content",
		NoteVersion: intPtr(2),
		Index:       intPtr(2),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/notes/"+noteID+"/contents", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Fatalf("failed to add content: status %d", rr.Code)
	}

	if len(errChan) != 0 {
		select {
		case err := <-errChan:
			t.Fatalf("error in websocket goroutine: %v", err)
		default:
		}
	}

	// Wait for the message from the WebSocket.
	select {
	case msg := <-msgChan:
		var event WebSocketEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			t.Fatalf("failed to unmarshal websocket message: %v", err)
		}
		if event.Type != "add_content" {
			t.Errorf("expected event type 'add_content', got '%s'", event.Type)
		}
		if event.NoteID != noteID {
			t.Errorf("expected note ID '%s', got '%s'", noteID, event.NoteID)
		}
		if event.NoteVersion != 3 {
			t.Errorf("expected version 3, got %d", event.NoteVersion)
		}
		if event.ContentID == "" {
			t.Error("expected non-empty content ID in event")
		}
		if event.Data != "Test content" {
			t.Errorf("expected data 'Test content', got '%s'", event.Data)
		}
		if event.ContentType != "text" {
			t.Errorf("expected content type 'text', got '%s'", event.ContentType)
		}
		if event.ContentVersion != 0 {
			t.Errorf("expected content version 0, got %d", event.ContentVersion)
		}
		if event.Index != 2 {
			t.Errorf("expected index 2, got %d", event.Index)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for websocket message")
	}
}

func TestNoteHandler_WebSocket_BroadcastOnDelete(t *testing.T) {
	// Arrange
	router, nuc, cuc, connManager := setupTestForBroadcast()
	noteID, _, _ := setUpNoteWithContents(nuc, cuc)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/notes/" + noteID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer conn.Close()
	errChan := make(chan error)
	msgChan := make(chan []byte)
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Failed to read message: %v", err)
			close(msgChan)
			return
		}
		msgChan <- msg
		if connManager.connections[noteID] != nil {
			errChan <- fmt.Errorf("connection was not removed after note deletion")
			return
		}
	}()

	// Act
	// Delete the note.
	requestBody := DeleteNoteRequest{
		NoteVersion: intPtr(2),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Fatalf("failed to delete note: status %d", rr.Code)
	}
	if len(errChan) != 0 {
		select {
		case err := <-errChan:
			t.Fatalf("error in websocket goroutine: %v", err)
		default:
		}
	}

	// Wait for the message from the WebSocket.
	select {
	case msg := <-msgChan:
		var event WebSocketEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			t.Fatalf("failed to unmarshal websocket message: %v", err)
		}
		if event.Type != "delete_note" {
			t.Errorf("expected event type 'delete_note', got '%s'", event.Type)
		}
		if event.NoteID != noteID {
			t.Errorf("expected note ID '%s', got '%s'", noteID, event.NoteID)
		}
		if event.NoteVersion != 3 {
			t.Errorf("expected version 3, got %d", event.NoteVersion)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for websocket message")
	}
}

func TestNoteHandler_WebSocket_BroadcastOnUpdateContent(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTestForBroadcast()
	noteID, contentID, _ := setUpNoteWithContents(nuc, cuc)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/notes/" + noteID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer conn.Close()

	msgChan := make(chan []byte)
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Failed to read message: %v", err)
			close(msgChan)
			return
		}
		msgChan <- msg
	}()

	// Act
	// Update the content.
	requestBody := UpdateContentRequest{
		Data:           "Updated content",
		ContentVersion: intPtr(0),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/"+noteID+"/contents/"+contentID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Fatalf("failed to update content: status %d", rr.Code)
	}

	select {
	case msg := <-msgChan:
		var event WebSocketEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			t.Fatalf("failed to unmarshal websocket message: %v", err)
		}
		if event.Type != "update_content" {
			t.Errorf("expected event type 'update_content', got '%s'", event.Type)
		}
		if event.NoteID != noteID {
			t.Errorf("expected note ID '%s', got '%s'", noteID, event.NoteID)
		}
		if event.ContentID != contentID {
			t.Errorf("expected content ID '%s', got '%s'", contentID, event.ContentID)
		}
		if event.Data != "Updated content" {
			t.Errorf("expected data 'Updated content', got '%s'", event.Data)
		}
		if event.ContentVersion != 1 {
			t.Errorf("expected version 1, got %d", event.ContentVersion)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for websocket message")
	}
}

func TestNoteHandler_WebSocket_BroadcastOnDeleteContent(t *testing.T) {
	// Arrange
	router, nuc, cuc, _ := setupTestForBroadcast()
	noteID, contentID, _ := setUpNoteWithContents(nuc, cuc)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/notes/" + noteID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer conn.Close()

	msgChan := make(chan []byte)
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Failed to read message: %v", err)
			close(msgChan)
			return
		}
		msgChan <- msg
	}()

	// Act
	// Delete the content.
	requestBody := DeleteContentRequest{
		ContentVersion: intPtr(0),
		NoteVersion:    intPtr(2),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/notes/"+noteID+"/contents/"+contentID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusNoContent {
		t.Fatalf("failed to delete content: status %d", rr.Code)
	}

	select {
	case msg := <-msgChan:
		var event WebSocketEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			t.Fatalf("failed to unmarshal websocket message: %v", err)
		}
		if event.Type != "delete_content" {
			t.Errorf("expected event type 'delete_content', got '%s'", event.Type)
		}
		if event.NoteID != noteID {
			t.Errorf("expected note ID '%s', got '%s'", noteID, event.NoteID)
		}
		if event.ContentID != contentID {
			t.Errorf("expected content ID '%s', got '%s'", contentID, event.ContentID)
		}
		if event.NoteVersion != 3 {
			t.Errorf("expected version 3, got %d", event.NoteVersion)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for websocket message")
	}
}

func TestNoteHandler_UpdateNote_Broadcast(t *testing.T) {
	// Arrange
	router, nuc, _, _ := setupTestForBroadcast()
	noteID, err := nuc.CreateNote("", "Original Title", "owner-1")
	if err != nil {
		t.Fatalf("setup: failed to create note: %v", err)
	}

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/notes/" + noteID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial websocket: %v", err)
	}
	defer conn.Close()

	msgChan := make(chan []byte)
	go func() {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Failed to read message: %v", err)
			close(msgChan)
			return
		}
		msgChan <- msg
	}()

	// Act
	newTitle := "Updated Title for Broadcast"
	requestBody := UpdateNoteRequest{
		Title:       newTitle,
		NoteVersion: intPtr(0),
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Fatalf("failed to update note: status %d", rr.Code)
	}

	select {
	case msg := <-msgChan:
		var event WebSocketEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			t.Fatalf("failed to unmarshal websocket message: %v", err)
		}
		if event.Type != "update_note" {
			t.Errorf("expected event type 'update_note', got '%s'", event.Type)
		}
		if event.NoteID != noteID {
			t.Errorf("expected note ID '%s', got '%s'", noteID, event.NoteID)
		}
		if event.Data != newTitle {
			t.Errorf("expected title '%s', got '%s'", newTitle, event.Data)
		}
		if event.NoteVersion != 1 {
			t.Errorf("expected version 1, got %d", event.NoteVersion)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for websocket message")
	}
}
