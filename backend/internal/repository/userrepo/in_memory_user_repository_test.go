package userrepo

import (
	"errors"
	"noteapp/internal/usecase/useruc"
	"testing"
)

func TestInMemoryUserRepository_SaveAndFindByID_Success(t *testing.T) {
	// Arrange
	repo := NewInMemoryUserRepository()
	user := &useruc.UserPO{
		ID:                "test-id",
		Username:          "testuser",
		PasswordHash:      "hash",
		AccessibleNoteIDs: []string{"n1", "n2"},
	}

	// Act
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	// Assert
	foundUser, err := repo.FindByID("test-id")
	if err != nil {
		t.Fatalf("FindByID() returned an unexpected error: %v", err)
	}
	if foundUser == nil {
		t.Fatal("FindByID() returned nil, expected a user")
	}
	if foundUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, foundUser.ID)
	}
	if foundUser.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, foundUser.Username)
	}
	if len(foundUser.AccessibleNoteIDs) != 2 {
		t.Fatalf("Expected 2 accessible note IDs, got %d", len(foundUser.AccessibleNoteIDs))
	}
	if foundUser.AccessibleNoteIDs[0] != "n1" || foundUser.AccessibleNoteIDs[1] != "n2" {
		t.Errorf("Expected AccessibleNoteIDs to be [n1, n2], got %v", foundUser.AccessibleNoteIDs)
	}
}

func TestInMemoryUserRepository_FindByID_NotFound(t *testing.T) {
	// Arrange
	repo := NewInMemoryUserRepository()

	// Act
	_, err := repo.FindByID("non-existent-id")

	// Assert
	if !errors.Is(err, useruc.ErrInvalidCredentials) {
		t.Errorf("Expected error %v, got %v", useruc.ErrInvalidCredentials, err)
	}
}

func TestInMemoryUserRepository_FindByUsername_Success(t *testing.T) {
	// Arrange
	repo := NewInMemoryUserRepository()
	user := &useruc.UserPO{
		ID:           "test-id",
		Username:     "testuser",
		PasswordHash: "hash",
	}
	repo.Save(user)

	// Act
	foundUser, err := repo.FindByUsername("testuser")

	// Assert
	if err != nil {
		t.Fatalf("FindByUsername() returned an unexpected error: %v", err)
	}
	if foundUser == nil {
		t.Fatal("FindByUsername() returned nil, expected a user")
	}
	if foundUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, foundUser.ID)
	}
}

func TestInMemoryUserRepository_FindByUsername_NotFound(t *testing.T) {
	// Arrange
	repo := NewInMemoryUserRepository()

	// Act
	_, err := repo.FindByUsername("non-existent-user")

	// Assert
	if !errors.Is(err, useruc.ErrInvalidCredentials) {
		t.Errorf("Expected error %v, got %v", useruc.ErrInvalidCredentials, err)
	}
}

func TestInMemoryUserRepository_DeepCopy(t *testing.T) {
	// Arrange
	repo := NewInMemoryUserRepository()
	user := &useruc.UserPO{
		ID:                "test-id",
		Username:          "testuser",
		AccessibleNoteIDs: []string{"n1"},
	}
	repo.Save(user)

	// Act
	foundUser, _ := repo.FindByID("test-id")
	foundUser.AccessibleNoteIDs[0] = "modified"

	// Assert
	refetchedUser, _ := repo.FindByID("test-id")
	if refetchedUser.AccessibleNoteIDs[0] == "modified" {
		t.Error("Modifying returned user PO affected the repository state")
	}
}

func TestInMemoryUserRepository_Save_DuplicateUsername(t *testing.T) {
	// Arrange
	repo := NewInMemoryUserRepository()
	user1 := &useruc.UserPO{
		ID:           "user1-id",
		Username:     "duplicateusername",
		PasswordHash: "hash1",
	}
	user2 := &useruc.UserPO{
		ID:           "user2-id",
		Username:     "duplicateusername",
		PasswordHash: "hash2",
	}

	// Act
	err1 := repo.Save(user1)
	err2 := repo.Save(user2)

	// Assert
	if err1 != nil {
		t.Fatalf("Save() for user1 returned an unexpected error: %v", err1)
	}

	if !errors.Is(err2, useruc.ErrUsernameExists) {
		t.Errorf("Expected error %v, got %v", useruc.ErrUsernameExists, err2)
	}
}
