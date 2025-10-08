package domain

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

// ContentType defines the type of content in a note.
type ContentType string

const (
	// TextContentType represents a text content block.
	TextContentType ContentType = "text"
	// ImageContentType represents an image content block.
	ImageContentType ContentType = "image"
)

// Content represents a block of content within a note.
type Content struct {
	ID   string
	Type ContentType
	Data string
}

// Note represents a note in the application.
type Note struct {
	ID            string
	OwnerID       string
	Title         string
	contents      []Content
	keywords      map[string][]Keyword
	Collaborators map[string]Permission
}

// NewNote creates a new Note instance.
// If id is empty, a new UUID will be generated.
func NewNote(id, title, ownerID string) (*Note, error) {
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
		contents:      []Content{},
		keywords:      make(map[string][]Keyword),
		Collaborators: make(map[string]Permission),
	}, nil
}

// AddCollaborator adds a user to the note's collaborators with a specific permission.
func (n *Note) AddCollaborator(callerID, collaboratorID string, permission Permission) error {
	if n.OwnerID != callerID {
		return ErrPermissionDenied
	}
	n.Collaborators[collaboratorID] = permission
	return nil
}

// Contents returns a copy of the note's contents.
func (n *Note) Contents() []Content {
	// Return a copy to prevent modification of the internal slice.
	contentsCopy := make([]Content, len(n.contents))
	copy(contentsCopy, n.contents)
	return contentsCopy
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

// AddContent adds a new content block to the note.
// If id is empty, a new UUID will be generated.
func (n *Note) AddContent(id, data string, contentType ContentType) string {
	if id == "" {
		id = uuid.New().String()
	}

	newContent := Content{
		ID:   id,
		Type: contentType,
		Data: data,
	}

	n.contents = append(n.contents, newContent)

	return id
}

// AddKeyword adds a new keyword to the note for a specific user.
func (n *Note) AddKeyword(userID string, keyword Keyword) {
	n.keywords[userID] = append(n.keywords[userID], keyword)
}

// UpdateContent updates an existing content block in the note.
func (n *Note) UpdateContent(id, data string) error {
	for i, content := range n.contents {
		if content.ID == id {
			n.contents[i].Data = data
			return nil
		}
	}
	return ErrContentNotFound
}

// DeleteContent removes a content block from the note.
func (n *Note) DeleteContent(id string) error {
	for i, content := range n.contents {
		if content.ID == id {
			n.contents = append(n.contents[:i], n.contents[i+1:]...)
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
