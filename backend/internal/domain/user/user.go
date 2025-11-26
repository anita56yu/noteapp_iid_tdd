package user

import "github.com/google/uuid"

// User represents a user in the system.
type User struct {
	id                string
	username          string
	passwordHash      string
	accessibleNoteIDs map[string]struct{}
}

// NewUser creates a new User.
// If id is empty, a new UUID will be generated.
func NewUser(id, username, passwordHash string) (*User, error) {
	// In a real application, you'd have more robust validation here.
	if username == "" {
		return nil, ErrEmptyUsername
	}
	if passwordHash == "" {
		return nil, ErrEmptyPasswordHash
	}

	if id == "" {
		id = uuid.New().String()
	}

	return &User{
		id:                id,
		username:          username,
		passwordHash:      passwordHash,
		accessibleNoteIDs: make(map[string]struct{}),
	}, nil
}

// ID returns the user's ID.
func (u *User) ID() string {
	return u.id
}

// Username returns the user's username.
func (u *User) Username() string {
	return u.username
}

// PasswordHash returns the user's password hash.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// AccessibleNoteIDs returns the IDs of notes accessible to the user.
func (u *User) AccessibleNoteIDs() []string {
	ids := make([]string, 0, len(u.accessibleNoteIDs))
	for id := range u.accessibleNoteIDs {
		ids = append(ids, id)
	}
	return ids
}

// AddAccessibleNoteID adds a note ID to the user's accessible notes.
func (u *User) AddAccessibleNoteID(noteID string) {
	u.accessibleNoteIDs[noteID] = struct{}{}
}

// RemoveAccessibleNoteID removes a note ID from the user's accessible notes.
func (u *User) RemoveAccessibleNoteID(noteID string) {
	delete(u.accessibleNoteIDs, noteID)
}
