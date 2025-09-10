package domain

import (
	"errors"
	"github.com/google/uuid"
)

// ErrEmptyTitle is returned when a note is created with an empty title.
var ErrEmptyTitle = errors.New("title cannot be empty")

// Note represents a note in the application.
type Note struct {
	ID      string
	Title   string
	Content string
}

// NewNote creates a new Note instance.
// If id is empty, a new UUID will be generated.
func NewNote(id, title, content string) (*Note, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if id == "" {
		id = uuid.New().String()
	}

	return &Note{
		ID:      id,
		Title:   title,
		Content: content,
	}, nil
}
