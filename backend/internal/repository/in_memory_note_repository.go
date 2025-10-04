package repository

import "sync"

// InMemoryNoteRepository is an in-memory implementation of NoteRepository.
type InMemoryNoteRepository struct {
	notes   map[string]*NotePO
	mutexes map[string]*sync.Mutex
	mu      sync.RWMutex
}

// NewInMemoryNoteRepository creates a new InMemoryNoteRepository.
func NewInMemoryNoteRepository() *InMemoryNoteRepository {
	return &InMemoryNoteRepository{
		notes:   make(map[string]*NotePO),
		mutexes: make(map[string]*sync.Mutex),
	}
}

// Save saves a note to the repository.
func (r *InMemoryNoteRepository) Save(note *NotePO) error {
	if note == nil {
		return ErrNilNote
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.notes[note.ID]; !ok {
		r.mutexes[note.ID] = &sync.Mutex{}
	}
	r.notes[note.ID] = note
	return nil
}

// FindByID finds a note by its ID.
func (r *InMemoryNoteRepository) FindByID(id string) (*NotePO, error) {
	r.mu.RLock()
	note, ok := r.notes[id]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrNoteNotFound
	}

	// Return a copy to prevent external modification
	copiedNote := &NotePO{
		ID:    note.ID,
		Title: note.Title,
	}
	copiedNote.Contents = make([]ContentPO, len(note.Contents))
	copy(copiedNote.Contents, note.Contents)
	copiedNote.Keywords = make(map[string][]string)
	for k, v := range note.Keywords {
		copiedNote.Keywords[k] = append([]string(nil), v...)
	}
	return copiedNote, nil
}

// Delete removes a note from the repository.
func (r *InMemoryNoteRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.notes[id]; !ok {
		return ErrNoteNotFound
	}
	delete(r.notes, id)
	delete(r.mutexes, id)
	return nil
}

// FindByKeywordForUser finds notes by a specific keyword for a given user.
func (r *InMemoryNoteRepository) FindByKeywordForUser(userID, keyword string) ([]*NotePO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var foundNotes []*NotePO
	for _, note := range r.notes {
		if userKeywords, ok := note.Keywords[userID]; ok {
			for _, k := range userKeywords {
				if k == keyword {
					foundNotes = append(foundNotes, note)
					break
				}
			}
		}
	}
	return foundNotes, nil
}

// LockNoteForUpdate locks a note for updating.
func (r *InMemoryNoteRepository) LockNoteForUpdate(noteID string) error {
	r.mu.RLock()
	noteMutex, ok := r.mutexes[noteID]
	r.mu.RUnlock()
	if !ok {
		return ErrNoteNotFound
	}
	noteMutex.Lock()
	return nil
}

// UnlockNoteForUpdate unlocks a note for updating.
func (r *InMemoryNoteRepository) UnlockNoteForUpdate(noteID string) error {
	r.mu.RLock()
	noteMutex, ok := r.mutexes[noteID]
	r.mu.RUnlock()
	if !ok {
		return ErrNoteNotFound
	}
	noteMutex.Unlock()
	return nil
}
