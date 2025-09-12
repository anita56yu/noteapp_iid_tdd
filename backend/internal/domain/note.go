package domain

import (
	"errors"
	"github.com/google/uuid"
)

// ErrEmptyTitle is returned when a note is created with an empty title.
var ErrEmptyTitle = errors.New("title cannot be empty")

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