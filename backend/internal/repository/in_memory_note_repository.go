package repository

import (
	"noteapp/internal/domain"
)

// InMemoryNoteRepository is an in-memory implementation of NoteRepository.
type InMemoryNoteRepository struct {
	notes map[string]*domain.Note
}

// NewInMemoryNoteRepository creates a new InMemoryNoteRepository.
func NewInMemoryNoteRepository() *InMemoryNoteRepository {
	return &InMemoryNoteRepository{
		notes: make(map[string]*domain.Note),
	}
}

// Save saves a note to the repository.
func (r *InMemoryNoteRepository) Save(note *domain.Note) error {
	if note == nil {
		return ErrNilNote
	}
	r.notes[note.ID] = note
	return nil
}

// FindByID finds a note by its ID.
func (r *InMemoryNoteRepository) FindByID(id string) (*domain.Note, error) {
	note, ok := r.notes[id]
	if !ok {
		return nil, ErrNoteNotFound
	}
	return note, nil
}

// Delete removes a note from the repository.
func (r *InMemoryNoteRepository) Delete(id string) error {
	if _, ok := r.notes[id]; !ok {
		return ErrNoteNotFound
	}
	delete(r.notes, id)
	return nil
}
