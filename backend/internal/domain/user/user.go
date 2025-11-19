package user

import "github.com/google/uuid"

// User represents a user in the system.
type User struct {
	id           uuid.UUID
	username     string
	passwordHash string
}

// NewUser creates a new User.
func NewUser(username, passwordHash string) (*User, error) {
	// In a real application, you'd have more robust validation here.
	if username == "" {
		return nil, ErrEmptyUsername
	}
	if passwordHash == "" {
		return nil, ErrEmptyPasswordHash
	}

	return &User{
		id:           uuid.New(),
		username:     username,
		passwordHash: passwordHash,
	}, nil
}

// ID returns the user's ID.
func (u *User) ID() uuid.UUID {
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
