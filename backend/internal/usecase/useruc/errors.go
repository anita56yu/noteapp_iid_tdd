package useruc

import (
	"errors"

	"noteapp/internal/domain/user"
)

// ErrUsernameExists is returned when a user tries to register with a username that already exists.
var ErrUsernameExists = errors.New("username already exists")

// ErrInvalidCredentials is returned when login fails due to incorrect username or password.
var ErrInvalidCredentials = errors.New("invalid credentials")

func mapDomainError(err error) error {
	if errors.Is(err, user.ErrEmptyUsername) || errors.Is(err, user.ErrEmptyPasswordHash) {
		return ErrInvalidCredentials
	}
	return err
}
