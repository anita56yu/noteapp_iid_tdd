package content

import "github.com/google/uuid"

// ContentType defines the type of content.
type ContentType string

const (
	// TextContentType is a text content type.
	TextContentType ContentType = "text/plain"
	// ImageContentType is an image content type.
	ImageContentType ContentType = "image/jpeg"
)

// Content represents a block of content within a note.
type Content struct {
	ID      string
	NoteID  string
	Data    string
	Type    ContentType
	Version int
}

// NewContent creates a new Content object.
// If id is an empty string, a new UUID will be generated.
func NewContent(id, noteID string, data string, contentType ContentType, version int) *Content {
	if id == "" {
		id = uuid.New().String()
	}
	return &Content{
		ID:      id,
		NoteID:  noteID,
		Data:    data,
		Type:    contentType,
		Version: version,
	}
}
