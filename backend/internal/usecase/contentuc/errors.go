package contentuc

import "errors"

// ErrUnsupportedContentType is returned when an unsupported content type is provided.
var ErrUnsupportedContentType = errors.New("unsupported content type")

// ErrContentNotFound is returned when a content is not found.
var ErrContentNotFound = errors.New("content not found")

// ErrConflict is returned when a version conflict occurs.
var ErrConflict = errors.New("conflict")
