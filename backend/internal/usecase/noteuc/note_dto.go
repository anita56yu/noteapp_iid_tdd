package noteuc

// Permission represents the permission level for a collaborator.
type Permission string

const (
	// PermissionRead allows a user to view a note.
	PermissionRead Permission = "read"
	// PermissionReadWrite allows a user to view and edit a note.
	PermissionReadWrite Permission = "read-write"
)

// ContentDTO represents a content block in a note.
type ContentDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// NoteDTO represents a note.
type NoteDTO struct {
	ID            string              `json:"id"`
	Title         string              `json:"title"`
	OwnerID       string              `json:"owner_id"`
	ContentIDs    []string            `json:"content_ids"`
	Contents      []*ContentDTO       `json:"contents"`
	Keywords      map[string][]string `json:"keywords"`
	Collaborators map[string]Permission `json:"collaborators"`
}