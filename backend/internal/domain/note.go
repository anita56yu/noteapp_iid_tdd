package domain

import (
	"errors"

	"github.com/google/uuid"
)

// ErrEmptyTitle is returned when a note is created with an empty title.
var ErrEmptyTitle = errors.New("title cannot be empty")

// ErrContentNotFound is returned when a content is not found.
var ErrContentNotFound = errors.New("content not found")

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
	ID       string
	Title    string
	contents []Content
}

// NewNote creates a new Note instance.
// If id is empty, a new UUID will be generated.
func NewNote(id, title string) (*Note, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if id == "" {
		id = uuid.New().String()
	}

	return &Note{
		ID:       id,
		Title:    title,
		contents: []Content{},
	}, nil
}

// Contents returns a copy of the note's contents.
func (n *Note) Contents() []Content {
	// Return a copy to prevent modification of the internal slice.
	contentsCopy := make([]Content, len(n.contents))
	copy(contentsCopy, n.contents)
	return contentsCopy
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
