package userrepo

import "errors"

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrUsernameInUseByExistingUser is returned when a user with the same username already exists.
var ErrUsernameInUseByExistingUser = errors.New("username in use by an existing user")

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	Save(user *UserPO) error
	FindByID(id string) (*UserPO, error)
	FindByUsername(username string) (*UserPO, error)
}
