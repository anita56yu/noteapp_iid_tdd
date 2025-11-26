package useruc_test

import (
	"testing"

	"noteapp/internal/domain/user"
	"noteapp/internal/repository/userrepo"
	"noteapp/internal/usecase/useruc"
)

func TestUserMapper_ToDomain(t *testing.T) {
	// Arrange
	mapper := &useruc.UserMapper{}
	userPO := &userrepo.UserPO{
		ID:                "test-id",
		Username:          "testuser",
		PasswordHash:      "testhash",
		AccessibleNoteIDs: []string{"note1", "note2"},
	}

	// Act
	userDomain, err := mapper.ToDomain(userPO)

	// Assert
	if err != nil {
		t.Fatalf("ToDomain() returned an unexpected error: %v", err)
	}

	if userDomain.ID() != userPO.ID {
		t.Errorf("Expected ID %s, got %s", userPO.ID, userDomain.ID())
	}
	if userDomain.Username() != userPO.Username {
		t.Errorf("Expected Username %s, got %s", userPO.Username, userDomain.Username())
	}
	if userDomain.PasswordHash() != userPO.PasswordHash {
		t.Errorf("Expected PasswordHash %s, got %s", userPO.PasswordHash, userDomain.PasswordHash())
	}
	if len(userDomain.AccessibleNoteIDs()) != len(userPO.AccessibleNoteIDs) {
		t.Errorf("Expected %d accessible note IDs, got %d", len(userPO.AccessibleNoteIDs), len(userDomain.AccessibleNoteIDs()))
	}
	for _, id := range userPO.AccessibleNoteIDs {
		found := false
		for _, domainID := range userDomain.AccessibleNoteIDs() {
			if id == domainID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected accessible note ID %s not found in domain user", id)
		}
	}
}

func TestUserMapper_FromDomain(t *testing.T) {
	// Arrange
	mapper := &useruc.UserMapper{}
	userDomain, _ := user.NewUser("test-id", "testuser", "testhash")
	userDomain.AddAccessibleNoteID("note1")
	userDomain.AddAccessibleNoteID("note2")

	// Act
	userPO := mapper.FromDomain(userDomain)

	// Assert
	if userPO.ID != userDomain.ID() {
		t.Errorf("Expected ID %s, got %s", userDomain.ID(), userPO.ID)
	}
	if userPO.Username != userDomain.Username() {
		t.Errorf("Expected Username %s, got %s", userDomain.Username(), userPO.Username)
	}
	if userPO.PasswordHash != userDomain.PasswordHash() {
		t.Errorf("Expected PasswordHash %s, got %s", userDomain.PasswordHash(), userPO.PasswordHash)
	}
	if len(userPO.AccessibleNoteIDs) != len(userDomain.AccessibleNoteIDs()) {
		t.Errorf("Expected %d accessible note IDs, got %d", len(userDomain.AccessibleNoteIDs()), len(userPO.AccessibleNoteIDs))
	}
	for _, id := range userDomain.AccessibleNoteIDs() {
		found := false
		for _, poID := range userPO.AccessibleNoteIDs {
			if id == poID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected accessible note ID %s not found in PO user", id)
		}
	}
}
