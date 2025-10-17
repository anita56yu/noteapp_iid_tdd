package usecase

// Permission represents the permission level of a user for a note.
type Permission string

const (
	// PermissionRead allows a user to read a note.
	PermissionRead Permission = "read"
	// PermissionReadWrite allows a user to read and write a note.
	PermissionReadWrite Permission = "read-write"
)

// NoteDTO represents a note data transfer object.
type NoteDTO struct {
	ID            string                `json:"id"`
	Title         string                `json:"title"`
	OwnerID       string                `json:"owner_id"`
	Contents      []*ContentDTO         `json:"contents"`
	Keywords      map[string][]string   `json:"keywords"`
	Collaborators map[string]Permission `json:"collaborators"`
}

// ContentDTO represents a content data transfer object.
type ContentDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
}