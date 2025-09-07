package repository

import "noteapp/internal/domain"

// NoteRepository defines the interface for note persistence.
type NoteRepository interface {
	Save(note *domain.Note) error
	FindByID(id string) (*domain.Note, error)
	Delete(id string) error
}
