package user_test

import (
	"testing"

	"noteapp/internal/domain/user"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("should create a new user with valid data", func(t *testing.T) {
		username := "testuser"
		passwordHash := "testhash"

		u, err := user.NewUser(username, passwordHash)

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, username, u.Username())
		assert.Equal(t, passwordHash, u.PasswordHash())
		assert.NotZero(t, u.ID())
	})

	t.Run("should return an error if username is empty", func(t *testing.T) {
		u, err := user.NewUser("", "testhash")

		assert.ErrorIs(t, err, user.ErrEmptyUsername)
		assert.Nil(t, u)
	})

	t.Run("should return an error if password hash is empty", func(t *testing.T) {
		u, err := user.NewUser("testuser", "")

		assert.ErrorIs(t, err, user.ErrEmptyPasswordHash)
		assert.Nil(t, u)
	})
}
