package useruc

// UserDTO represents a user's data transfer object.
type UserDTO struct {
	ID                string   `json:"id"`
	Username          string   `json:"username"`
	AccessibleNoteIDs []string `json:"accessibleNoteIDs"`
}
