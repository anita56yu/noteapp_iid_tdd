package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"noteapp/internal/usecase"

	"github.com/go-chi/chi/v5"
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

	noteID, err := h.usecase.CreateNote("", req.Title, req.Content)
	if err != nil {
		if errors.Is(err, usecase.ErrEmptyTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/notes/%s", noteID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateNoteResponse{ID: noteID})
}

// GetNoteByID is the handler for the GET /notes/{id} endpoint.
func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	note, err := h.usecase.GetNoteByID(id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrInvalidID):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

// DeleteNote is the handler for the DELETE /notes/{id} endpoint.
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.usecase.DeleteNote(id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrInvalidID):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
