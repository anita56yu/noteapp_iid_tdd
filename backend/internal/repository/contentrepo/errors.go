package contentrepo

import "errors"

var (
	// ErrContentNotFound is returned when a content is not found.
	ErrContentNotFound = errors.New("content not found")
	// ErrContentConflict is returned when a version conflict occurs.
	ErrContentConflict = errors.New("content conflict")
)
