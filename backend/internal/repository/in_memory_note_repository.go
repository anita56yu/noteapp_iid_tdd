package repository

// InMemoryNoteRepository is an in-memory implementation of NoteRepository.
type InMemoryNoteRepository struct {
	notes map[string]*NotePO
}

// NewInMemoryNoteRepository creates a new InMemoryNoteRepository.
func NewInMemoryNoteRepository() *InMemoryNoteRepository {
	return &InMemoryNoteRepository{
		notes: make(map[string]*NotePO),
	}
}

// Save saves a note to the repository.
func (r *InMemoryNoteRepository) Save(note *NotePO) error {
	if note == nil {
		return ErrNilNote
	}
	r.notes[note.ID] = note
	return nil
}

// FindByID finds a note by its ID.
func (r *InMemoryNoteRepository) FindByID(id string) (*NotePO, error) {
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
