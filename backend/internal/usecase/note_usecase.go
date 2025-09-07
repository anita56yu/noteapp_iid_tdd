package usecase

import (
	"noteapp/internal/domain"
	"noteapp/internal/repository"
)

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
