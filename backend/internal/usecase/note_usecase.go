package usecase

import (
	"errors"
	"noteapp/internal/domain"
	"noteapp/internal/repository"
)

// ErrInvalidID is returned when an invalid ID is provided.
var ErrInvalidID = errors.New("invalid ID")

// NoteUsecase handles the business logic for notes.
type NoteUsecase struct {
	repo repository.NoteRepository
}

// NewNoteUsecase creates a new NoteUsecase.
func NewNoteUsecase(repo repository.NoteRepository) *NoteUsecase {
	return &NoteUsecase{repo: repo}
}

// CreateNote creates a new note.
func (uc *NoteUsecase) CreateNote(id, title, content string) (string, error) {
	note, err := domain.NewNote(id, title, content)
	if err != nil {
		return "", err
	}

	if err := uc.repo.Save(note); err != nil {
		return "", err
	}

	return note.ID, nil
}

// GetNoteByID retrieves a note by its ID.
func (uc *NoteUsecase) GetNoteByID(id string) (*NoteDTO, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	note, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &NoteDTO{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
	}, nil
}
