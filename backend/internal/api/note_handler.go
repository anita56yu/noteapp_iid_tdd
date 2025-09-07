package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"noteapp/internal/usecase"
)

// NoteHandler handles HTTP requests for notes.
type NoteHandler struct {
	usecase *usecase.NoteUsecase
}

// NewNoteHandler creates a new NoteHandler.
func NewNoteHandler(uc *usecase.NoteUsecase) *NoteHandler {
	return &NoteHandler{usecase: uc}
}

// CreateNoteRequest represents the request body for creating a note.
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateNoteResponse represents the response body for creating a note.
type CreateNoteResponse struct {
	ID string `json:"id"`
}

// CreateNote is the handler for the POST /notes endpoint.
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// We pass an empty string for the ID to let the usecase generate it.
	noteID, err := h.usecase.CreateNote("", req.Title, req.Content)
	if err != nil {
		// In a real app, we'd check the error type to return the correct status code.
		// For now, we'll assume any error from the usecase is a bad request.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/notes/%s", noteID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateNoteResponse{ID: noteID})
}
