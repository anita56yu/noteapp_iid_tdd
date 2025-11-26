package userrepo

// UserPO represents the persistent state of a user.
type UserPO struct {
	ID                string
	Username          string
	PasswordHash      string
	AccessibleNoteIDs []string
}
