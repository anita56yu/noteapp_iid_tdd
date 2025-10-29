package noterepo

import (
	"sync"
)

// InMemoryNoteRepository is an in-memory implementation of NoteRepository.
type InMemoryNoteRepository struct {
	notes map[string]*NotePO
	mu    sync.RWMutex
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
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.notes[note.ID]; ok {
		if existing.Version != note.Version {
			return ErrNoteConflict
		}
		note.Version++
	} else {
		note.Version = 0
	}

	r.notes[note.ID] = note
	return nil
}

// FindByID retrieves a note by its ID.

func (r *InMemoryNoteRepository) FindByID(id string) (*NotePO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	note, ok := r.notes[id]
	if !ok {
		return nil, ErrNoteNotFound
	}
	// Return a copy to prevent race conditions in concurrent tests
	newNote := &NotePO{
		ID:            note.ID,
		OwnerID:       note.OwnerID,
		Title:         note.Title,
		Version:       note.Version,
		ContentIDs:    make([]string, len(note.ContentIDs)),
		Keywords:      make(map[string][]string),
		Collaborators: make(map[string]string),
	}
	copy(newNote.ContentIDs, note.ContentIDs)
	for k, v := range note.Keywords {
		newNote.Keywords[k] = append([]string{}, v...)
	}

	for k, v := range note.Collaborators {
		newNote.Collaborators[k] = v
	}
	return newNote, nil

}

// Delete removes a note from the repository.
func (r *InMemoryNoteRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.notes[id]; !ok {
		return ErrNoteNotFound
	}
	delete(r.notes, id)
	return nil
}

// TODO: add deep copy where necessary
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

// TODO: add deep copy where necessary
// GetAccessibleNotesByUserID retrieves all notes where the user is either the owner or a collaborator.
func (r *InMemoryNoteRepository) GetAccessibleNotesByUserID(userID string) ([]*NotePO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var accessibleNotes []*NotePO
	for _, note := range r.notes {
		if note.OwnerID == userID {
			accessibleNotes = append(accessibleNotes, note)
			continue
		}
		if _, ok := note.Collaborators[userID]; ok {
			accessibleNotes = append(accessibleNotes, note)
		}
	}
	return accessibleNotes, nil
}
