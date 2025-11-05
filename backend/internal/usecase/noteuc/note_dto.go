package noteuc

// Permission represents the permission level for a collaborator.
type Permission string

const (
	// PermissionRead allows a user to view a note.
	PermissionRead Permission = "read"
	// PermissionReadWrite allows a user to view and edit a note.
	PermissionReadWrite Permission = "read-write"
)

// NoteDTO represents a note.
type NoteDTO struct {
	ID            string                `json:"id"`
	Title         string                `json:"title"`
	OwnerID       string                `json:"owner_id"`
	Version       int                   `json:"version"`
	ContentIDs    []string              `json:"content_ids"`
	Keywords      map[string][]string   `json:"keywords"`
	Collaborators map[string]Permission `json:"collaborators"`
}
