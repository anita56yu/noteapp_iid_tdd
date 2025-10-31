package contentrepo

import (
	"sync"
)

// InMemoryContentRepository is an in-memory implementation of ContentRepository.
type InMemoryContentRepository struct {
	mu       sync.RWMutex
	contents map[string]*ContentPO
}

// NewInMemoryContentRepository creates a new InMemoryContentRepository.
func NewInMemoryContentRepository() *InMemoryContentRepository {
	return &InMemoryContentRepository{
		contents: make(map[string]*ContentPO),
	}
}

// Save saves a content to the repository.
func (r *InMemoryContentRepository) Save(c *ContentPO) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.contents[c.ID]; ok {
		if existing.Version != c.Version {
			return ErrContentConflict
		}
		c.Version++
	} else {
		c.Version = 0
	}

	r.contents[c.ID] = c
	return nil
}

// GetByID retrieves a content by its ID.
func (r *InMemoryContentRepository) GetByID(id string) (*ContentPO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if c, ok := r.contents[id]; ok {
		// Return a copy to prevent race conditions
		copy := *c
		return &copy, nil
	}
	return nil, ErrContentNotFound
}

// GetAllByNoteID retrieves all contents for a given note ID.
func (r *InMemoryContentRepository) GetAllByNoteID(noteID string) ([]*ContentPO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*ContentPO
	for _, c := range r.contents {
		if c.NoteID == noteID {
			results = append(results, c)
		}
	}
	return results, nil
}

// Delete removes a content from the repository.
func (r *InMemoryContentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.contents[id]; !ok {
		return ErrContentNotFound
	}
	delete(r.contents, id)
	return nil
}

// DeleteAllByNoteID removes all contents associated with a given note ID.
func (r *InMemoryContentRepository) DeleteAllByNoteID(noteID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, c := range r.contents {
		if c.NoteID == noteID {
			delete(r.contents, id)
		}
	}
	return nil
}
