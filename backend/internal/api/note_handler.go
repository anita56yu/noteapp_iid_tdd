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
	Title string `json:"title"`
}

// CreateNoteResponse represents the response body for creating a note.
type CreateNoteResponse struct {
	ID string `json:"id"`
}

// AddContentRequest represents the request body for adding content to a note.
type AddContentRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// UpdateContentRequest represents the request body for updating content in a note.
type UpdateContentRequest struct {
	Data string `json:"data"`
}

// TagNoteRequest represents the request body for tagging a note.
type TagNoteRequest struct {
	Keyword string `json:"keyword"`
}

// ShareNoteRequest represents the request body for sharing a note.
type ShareNoteRequest struct {
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
}

type RevokeAccessRequest struct {
	UserID string `json:"user_id"`
}

var ErrUnsupportedContentType = errors.New("unsupported content type")

// CreateNote is the handler for the POST /notes endpoint.
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For now, we'll use a placeholder ownerID.
	// This will be replaced with actual user authentication later.
	ownerID := "placeholder-owner-id"

	noteID, err := h.usecase.CreateNote("", req.Title, ownerID)
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

// AddContent is the handler for the POST /notes/{id}/contents endpoint.
func (h *NoteHandler) AddContent(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")

	var req AddContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contentType, err := mapToDomainContentType(req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentID, err := h.usecase.AddContent(noteID, "", req.Data, contentType)
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID string `json:"id"`
	}{ID: contentID})
}

// UpdateContent is the handler for the PUT /notes/{id}/contents/{contentId} endpoint.
func (h *NoteHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")
	contentID := chi.URLParam(r, "contentId")

	var req UpdateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.usecase.UpdateContent(noteID, contentID, req.Data)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound), errors.Is(err, usecase.ErrContentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrInvalidID):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteContent is the handler for the DELETE /notes/{id}/contents/{contentId} endpoint.
func (h *NoteHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "id")
	contentID := chi.URLParam(r, "contentId")

	err := h.usecase.DeleteContent(noteID, contentID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound), errors.Is(err, usecase.ErrContentNotFound):
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

// TagNote is the handler for the POST /users/{userID}/notes/{noteID}/keyword endpoint.
func (h *NoteHandler) TagNote(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	noteID := chi.URLParam(r, "noteID")

	var req TagNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.TagNote(noteID, userID, req.Keyword); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrEmptyKeyword):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// FindNotesByKeyword is the handler for the GET /users/{userID}/notes?keyword={keyword} endpoint.
func (h *NoteHandler) FindNotesByKeyword(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	keyword := r.URL.Query().Get("keyword")

	notes, err := h.usecase.FindNotesByKeyword(userID, keyword)
	if err != nil {
		http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

// UntagNote is the handler for the DELETE /users/{userID}/notes/{noteID}/keyword/{keyword} endpoint.
func (h *NoteHandler) UntagNote(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	noteID := chi.URLParam(r, "noteID")
	keyword := chi.URLParam(r, "keyword")

	if err := h.usecase.UntagNote(noteID, userID, keyword); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound), errors.Is(err, usecase.ErrUserNotFound), errors.Is(err, usecase.ErrKeywordNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrEmptyKeyword):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ShareNote is the handler for the POST /users/{ownerID}/notes/{noteID}/shares endpoint.
func (h *NoteHandler) ShareNote(w http.ResponseWriter, r *http.Request) {
	ownerID := chi.URLParam(r, "ownerID")
	noteID := chi.URLParam(r, "noteID")

	var req ShareNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.ShareNote(noteID, ownerID, req.UserID, req.Permission); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrPermissionDenied):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, usecase.ErrUnsupportedPermissionType):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAccessibleNotesForUser is the handler for the GET /users/{userID}/notes endpoint.
func (h *NoteHandler) GetAccessibleNotesForUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	notes, err := h.usecase.GetAccessibleNotesForUser(userID)
	if err != nil {
		http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

// RevokeAccess is the handler for the DELETE /users/{ownerID}/notes/{noteID}/shares endpoint.
func (h *NoteHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	ownerID := chi.URLParam(r, "ownerID")
	noteID := chi.URLParam(r, "noteID")

	var req RevokeAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.RevokeAccess(noteID, ownerID, req.UserID); err != nil {
		switch {
		case errors.Is(err, usecase.ErrNoteNotFound), errors.Is(err, usecase.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, usecase.ErrPermissionDenied):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapToDomainContentType(ct string) (usecase.ContentType, error) {
	switch ct {
	case "text":
		return usecase.TextContentType, nil
	case "image":
		return usecase.ImageContentType, nil
	default:
		return usecase.TextContentType, ErrUnsupportedContentType
	}
}
