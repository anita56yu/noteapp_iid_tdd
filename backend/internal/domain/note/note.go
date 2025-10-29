package note

import (
	"errors"

	"github.com/google/uuid"
)

// ErrEmptyTitle is returned when a note is created with an empty title.
var ErrEmptyTitle = errors.New("title cannot be empty")

// ErrContentNotFound is returned when a content is not found.
var ErrContentNotFound = errors.New("content not found")

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrKeywordNotFound is returned when a keyword is not found.
var ErrKeywordNotFound = errors.New("keyword not found")

// ErrPermissionDenied is returned when a user is not authorized to perform an action.
var ErrPermissionDenied = errors.New("permission denied")

// Permission defines the access level for a collaborator.
type Permission string

const (
	// ReadOnly allows a user to view a note.
	ReadOnly Permission = "read"
	// ReadWrite allows a user to view and edit a note.
	ReadWrite Permission = "read-write"
)

// Note represents a note in the application.
type Note struct {
	ID            string
	OwnerID       string
	Title         string
	Version       int
	ContentIDs    []string
	keywords      map[string][]Keyword
	Collaborators map[string]Permission
}

// NewNoteWithVersion creates a new Note instance with a specific version.
// If id is empty, a new UUID will be generated.
func NewNoteWithVersion(id, title, ownerID string, version int) (*Note, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if id == "" {
		id = uuid.New().String()
	}

	return &Note{
		ID:            id,
		OwnerID:       ownerID,
		Title:         title,
		Version:       version,
		ContentIDs:    []string{},
		keywords:      make(map[string][]Keyword),
		Collaborators: make(map[string]Permission),
	}, nil
}

// NewNote creates a new Note instance with a default version of 0.
// If id is empty, a new UUID will be generated.
func NewNote(id, title, ownerID string) (*Note, error) {
	return NewNoteWithVersion(id, title, ownerID, 0)
}

// AddCollaborator adds a user to the note's collaborators with a specific permission.
func (n *Note) AddCollaborator(callerID, collaboratorID string, permission Permission) error {
	if n.OwnerID != callerID {
		return ErrPermissionDenied
	}
	n.Collaborators[collaboratorID] = permission
	return nil
}

// RemoveCollaborator removes a collaborator from the note.
func (n *Note) RemoveCollaborator(callerID string, collaboratorID string) error {
	if n.OwnerID != callerID {
		return ErrPermissionDenied
	}
	if _, ok := n.Collaborators[collaboratorID]; !ok {
		return ErrUserNotFound
	}
	delete(n.Collaborators, collaboratorID)
	delete(n.keywords, collaboratorID)
	return nil
}

// Keywords returns a deep copy of the note's keywords.
func (n *Note) Keywords() map[string][]Keyword {
	keywordsCopy := make(map[string][]Keyword)
	for userID, keywords := range n.keywords {
		userKeywordsCopy := make([]Keyword, len(keywords))
		copy(userKeywordsCopy, keywords)
		keywordsCopy[userID] = userKeywordsCopy
	}
	return keywordsCopy
}

// UserKeywords returns a copy of the note's keywords for a specific user.
func (n *Note) UserKeywords(userID string) []Keyword {
	// Return a copy to prevent modification of the internal slice.
	keywordsCopy := make([]Keyword, len(n.keywords[userID]))
	copy(keywordsCopy, n.keywords[userID])
	return keywordsCopy
}

// AddContentID adds a new content ID to the note.
func (n *Note) AddContentID(id string) {
	n.ContentIDs = append(n.ContentIDs, id)
}

// AddKeyword adds a new keyword to the note for a specific user.
func (n *Note) AddKeyword(userID string, keyword Keyword) {
	n.keywords[userID] = append(n.keywords[userID], keyword)
}

// RemoveContentID removes a content ID from the note.
func (n *Note) RemoveContentID(id string) error {
	for i, contentID := range n.ContentIDs {
		if contentID == id {
			n.ContentIDs = append(n.ContentIDs[:i], n.ContentIDs[i+1:]...)
			return nil
		}
	}
	return ErrContentNotFound
}

// RemoveKeyword removes a keyword from the note for a specific user.
func (n *Note) RemoveKeyword(userID string, keyword Keyword) error {
	userKeywords, ok := n.keywords[userID]
	if !ok {
		return ErrUserNotFound
	}

	for i, k := range userKeywords {
		if k == keyword {
			n.keywords[userID] = append(userKeywords[:i], userKeywords[i+1:]...)
			return nil
		}
	}
	return ErrKeywordNotFound
}
