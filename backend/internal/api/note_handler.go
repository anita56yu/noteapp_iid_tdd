package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"noteapp/internal/usecase/contentuc"
	"noteapp/internal/usecase/noteuc"

	"github.com/go-chi/chi/v5"
)

// NoteHandler handles HTTP requests for notes.
type NoteHandler struct {
	noteUsecase    *noteuc.NoteUsecase
	contentUsecase *contentuc.ContentUsecase
}

// NewNoteHandler creates a new NoteHandler.
func NewNoteHandler(nuc *noteuc.NoteUsecase, cuc *contentuc.ContentUsecase) *NoteHandler {
	return &NoteHandler{noteUsecase: nuc, contentUsecase: cuc}
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
	Type        string `json:"type"`
	Data        string `json:"data"`
	NoteVersion *int   `json:"note_version"`
}

// UpdateContentRequest represents the request body for updating content in a note.
type UpdateContentRequest struct {
	Data    string `json:"data"`
	Version *int   `json:"version"`
}

// DeleteContentRequest represents the request body for deleting content in a note.
type DeleteContentRequest struct {
	ContentVersion *int `json:"content_version"`
	NoteVersion    *int `json:"note_version"`
}

// DeleteNoteRequest represents the request body for deleting a note.
type DeleteNoteRequest struct {
	NoteVersion *int `json:"note_version"`
}

// TagNoteRequest represents the request body for tagging a note.
type TagNoteRequest struct {
	Keyword     string `json:"keyword"`
	NoteVersion *int   `json:"note_version"`
}

// UntagNoteRequest represents the request body for untagging a note.
type UntagNoteRequest struct {
	NoteVersion *int `json:"note_version"`
}

// ShareNoteRequest represents the request body for sharing a note.
type ShareNoteRequest struct {
	UserID      string `json:"user_id"`
	Permission  string `json:"permission"`
	NoteVersion *int   `json:"note_version"`
}

type RevokeAccessRequest struct {
	UserID      string `json:"user_id"`
	NoteVersion *int   `json:"note_version"`
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

	noteID, err := h.noteUsecase.CreateNote("", req.Title, ownerID)
	if err != nil {
		if errors.Is(err, noteuc.ErrEmptyTitle) {
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

	note, err := h.noteUsecase.GetNoteByID(id)
	if err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrInvalidID):
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

	var req DeleteNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NoteVersion == nil {
		http.Error(w, "note_version is required", http.StatusBadRequest)
		return
	}



	// Now, delete the note itself.
	if err := h.noteUsecase.DeleteNote(id, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrInvalidID):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

		// First, delete all associated content.
	if err := h.contentUsecase.DeleteAllContentsByNoteID(id); err != nil {
		http.Error(w, "An internal error occurred while deleting content", http.StatusInternalServerError)
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

	if req.NoteVersion == nil {
		http.Error(w, "note_version is required", http.StatusBadRequest)
		return
	}

	contentType, err := mapToContentUsecaseContentType(req.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the content first.
	contentID, err := h.contentUsecase.CreateContent(noteID, "", req.Data, contentType)
	if err != nil {
		// Error handling for content creation can be added here.
		http.Error(w, "Failed to create content", http.StatusInternalServerError)
		return
	}

	// Then, add the content ID to the note.
	if err := h.noteUsecase.AddContent(noteID, contentID, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrInvalidID):
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
	contentID := chi.URLParam(r, "contentId")

	var req UpdateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Version == nil {
		http.Error(w, "version is required", http.StatusBadRequest)
		return
	}

	if err := h.contentUsecase.UpdateContent(contentID, req.Data, *req.Version); err != nil {
		switch {
		case errors.Is(err, contentuc.ErrContentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, contentuc.ErrConflict):
			http.Error(w, err.Error(), http.StatusConflict)
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

	var req DeleteContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ContentVersion == nil || req.NoteVersion == nil {
		http.Error(w, "content_version and note_version are required", http.StatusBadRequest)
		return
	}

	// First, remove the content ID from the note.
	if err := h.noteUsecase.RemoveContent(noteID, contentID, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound), errors.Is(err, noteuc.ErrContentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrInvalidID):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	// Then, delete the content itself.
	if err := h.contentUsecase.DeleteContent(contentID, *req.ContentVersion); err != nil {
		switch {
		case errors.Is(err, contentuc.ErrContentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
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

	if req.NoteVersion == nil {
		http.Error(w, "note_version is required", http.StatusBadRequest)
		return
	}

	if err := h.noteUsecase.TagNote(noteID, userID, req.Keyword, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrEmptyKeyword):
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

	notes, err := h.noteUsecase.FindNotesByKeyword(userID, keyword)
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

	var req UntagNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NoteVersion == nil {
		http.Error(w, "note_version is required", http.StatusBadRequest)
		return
	}

	if err := h.noteUsecase.UntagNote(noteID, userID, keyword, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound), errors.Is(err, noteuc.ErrUserNotFound), errors.Is(err, noteuc.ErrKeywordNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrEmptyKeyword):
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

	if req.NoteVersion == nil {
		http.Error(w, "note_version is required", http.StatusBadRequest)
		return
	}

	if err := h.noteUsecase.ShareNote(noteID, ownerID, req.UserID, req.Permission, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrPermissionDenied):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, noteuc.ErrUnsupportedPermissionType):
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

	notes, err := h.noteUsecase.GetAccessibleNotesForUser(userID)
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

	if req.NoteVersion == nil {
		http.Error(w, "note_version is required", http.StatusBadRequest)
		return
	}

	if err := h.noteUsecase.RevokeAccess(noteID, ownerID, req.UserID, *req.NoteVersion); err != nil {
		switch {
		case errors.Is(err, noteuc.ErrNoteNotFound), errors.Is(err, noteuc.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, noteuc.ErrPermissionDenied):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "An internal error occurred", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func mapToContentUsecaseContentType(ct string) (contentuc.ContentType, error) {
	switch ct {
	case "text":
		return contentuc.TextContentType, nil
	case "image":
		return contentuc.ImageContentType, nil
	default:
		return contentuc.TextContentType, ErrUnsupportedContentType
	}
}
