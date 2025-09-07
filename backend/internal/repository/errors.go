package repository

import "errors"

var (
	// ErrNoteNotFound is returned when a note is not found.
	ErrNoteNotFound = errors.New("note not found")
	// ErrNilNote is returned when a nil note is passed to a method.
	ErrNilNote = errors.New("note cannot be nil")
)
