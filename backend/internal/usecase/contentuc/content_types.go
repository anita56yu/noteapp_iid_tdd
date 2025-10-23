package contentuc

// ContentType defines the type of content in a note.
type ContentType string

const (
	// TextContentType represents a text content block.
	TextContentType ContentType = "text"
	// ImageContentType represents an image content block.
	ImageContentType ContentType = "image"
)
