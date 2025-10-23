package noterepo

import "errors"

var (
	// ErrNoteNotFound is returned when a note is not found.
	ErrNoteNotFound = errors.New("note not found")
	// ErrNoteConflict is returned when a version conflict occurs.
	ErrNoteConflict = errors.New("note conflict")
	// ErrNilNote is returned when a nil note is passed.
	ErrNilNote = errors.New("note cannot be nil")
)
