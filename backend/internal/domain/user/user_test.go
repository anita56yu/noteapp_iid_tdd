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

		u, err := user.NewUser("", username, passwordHash)

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, username, u.Username())
		assert.Equal(t, passwordHash, u.PasswordHash())
		assert.NotZero(t, u.ID())
		assert.Empty(t, u.AccessibleNoteIDs())
	})

	t.Run("should return an error if username is empty", func(t *testing.T) {
		u, err := user.NewUser("", "", "testhash")

		assert.ErrorIs(t, err, user.ErrEmptyUsername)
		assert.Nil(t, u)
	})

	t.Run("should return an error if password hash is empty", func(t *testing.T) {
		u, err := user.NewUser("", "testuser", "")

		assert.ErrorIs(t, err, user.ErrEmptyPasswordHash)
		assert.Nil(t, u)
	})
}

func TestUserAccessibleNotes(t *testing.T) {
	user, _ := user.NewUser("", "testuser", "testhash")

	t.Run("should have no accessible notes initially", func(t *testing.T) {
		assert.Empty(t, user.AccessibleNoteIDs())
	})

	t.Run("should add accessible note IDs", func(t *testing.T) {
		noteID1 := "note1"
		noteID2 := "note2"

		user.AddAccessibleNoteID(noteID1)
		user.AddAccessibleNoteID(noteID2)

		expected := []string{noteID1, noteID2}
		assert.ElementsMatch(t, expected, user.AccessibleNoteIDs())
	})

	t.Run("should not add duplicate accessible note IDs", func(t *testing.T) {
		noteID1 := "note1"

		initialCount := len(user.AccessibleNoteIDs())
		user.AddAccessibleNoteID(noteID1)

		assert.Len(t, user.AccessibleNoteIDs(), initialCount)
	})

	t.Run("should remove accessible note IDs", func(t *testing.T) {
		noteIDToRemove := "note1"

		user.RemoveAccessibleNoteID(noteIDToRemove)

		assert.NotContains(t, user.AccessibleNoteIDs(), noteIDToRemove)
		assert.Len(t, user.AccessibleNoteIDs(), 1) // Only note2 should remain
	})

	t.Run("should not change anything if removing a non-existent note ID", func(t *testing.T) {
		noteIDNonExistent := "nonexistent"
		initialCount := len(user.AccessibleNoteIDs())

		user.RemoveAccessibleNoteID(noteIDNonExistent)

		assert.Len(t, user.AccessibleNoteIDs(), initialCount)
	})

	t.Run("modifying returned slice should not affect internal state", func(t *testing.T) {
		ids := user.AccessibleNoteIDs()
		ids = append(ids, "newID")

		assert.NotContains(t, user.AccessibleNoteIDs(), "newID")
	})
}
