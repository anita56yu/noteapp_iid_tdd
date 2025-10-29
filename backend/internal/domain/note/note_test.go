package note

import (
	"testing"
)

func TestNewNote_ValidCreation_WithInjectedID(t *testing.T) {
	id := "test-id"
	title := "Test Note"
	ownerID := "owner-1"

	note, err := NewNoteWithVersion(id, title, ownerID, 0)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}

	if note.ID != id {
		t.Errorf("Expected ID to be '%s', but got '%s'", id, note.ID)
	}

	if note.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", note.Version)
	}

	if note.ContentIDs == nil {
		t.Fatalf("Expected ContentIDs to be an empty slice, but it was nil")
	}

	if len(note.ContentIDs) != 0 {
		t.Errorf("Expected ContentIDs to be empty, but got %d elements", len(note.ContentIDs))
	}
}

func TestNewNote_ValidCreation_WithGeneratedID(t *testing.T) {
	title := "Test Note"
	ownerID := "owner-1"
	note, err := NewNoteWithVersion("", title, ownerID, 0)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}
	if note.ID == "" {
		t.Errorf("Expected ID to be non-empty, but it was empty")
	}
	if note.Title != title {
		t.Errorf("Expected title to be '%s', but got '%s'", title, note.Title)
	}

	if note.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", note.Version)
	}

	if note.ContentIDs == nil {
		t.Fatalf("Expected ContentIDs to be an empty slice, but it was nil")
	}

	if len(note.ContentIDs) != 0 {
		t.Errorf("Expected ContentIDs to be empty, but got %d elements", len(note.ContentIDs))
	}
}

func TestNewNote_EmptyTitle(t *testing.T) {
	_, err := NewNote("", "", "owner-1")
	if err == nil {
		t.Fatal("Expected an error for empty title, but got nil")
	}
	if err != ErrEmptyTitle {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrEmptyTitle, err)
	}
}

func TestNote_AddContentID(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	contentID := "content-1"

	// Act
	note.AddContentID(contentID)

	// Assert
	if len(note.ContentIDs) != 1 {
		t.Fatalf("Expected 1 content ID, but got %d", len(note.ContentIDs))
	}
	if note.ContentIDs[0] != contentID {
		t.Errorf("Expected content ID to be '%s', but got '%s'", contentID, note.ContentIDs[0])
	}
}

func TestNote_RemoveContentID(t *testing.T) {
	t.Run("should remove content ID successfully", func(t *testing.T) {
		// Arrange
		note, _ := NewNote("note-1", "Test Note", "owner-1")
		note.AddContentID("content-1")
		note.AddContentID("content-2")

		// Act
		err := note.RemoveContentID("content-1")

		// Assert
		if err != nil {
			t.Fatalf("RemoveContentID returned an unexpected error: %v", err)
		}
		if len(note.ContentIDs) != 1 {
			t.Fatalf("Expected 1 content ID, but got %d", len(note.ContentIDs))
		}
		if note.ContentIDs[0] != "content-2" {
			t.Errorf("Expected remaining content ID to be 'content-2', but got '%s'", note.ContentIDs[0])
		}
	})

	t.Run("should return error when content ID not found", func(t *testing.T) {
		// Arrange
		note, _ := NewNote("note-1", "Test Note", "owner-1")
		note.AddContentID("content-1")

		// Act
		err := note.RemoveContentID("non-existent-id")

		// Assert
		if err != ErrContentNotFound {
			t.Errorf("Expected error to be '%v', but got '%v'", ErrContentNotFound, err)
		}
	})
}

func TestNote_AddKeyword(t *testing.T) {
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	userID := "user-1"
	keyword, _ := NewKeyword("test-keyword")

	note.AddKeyword(userID, keyword)

	keywords := note.UserKeywords(userID)
	if len(keywords) != 1 {
		t.Fatalf("Expected 1 keyword for user, but got %d", len(keywords))
	}
	if keywords[0] != keyword {
		t.Errorf("Expected keyword to be '%v', but got '%v'", keyword, keywords[0])
	}
}

func TestNote_RemoveKeyword_Success(t *testing.T) {
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	userID := "user-1"
	keyword1, _ := NewKeyword("keyword1")
	keyword2, _ := NewKeyword("keyword2")
	note.AddKeyword(userID, keyword1)
	note.AddKeyword(userID, keyword2)

	err := note.RemoveKeyword(userID, keyword1)
	if err != nil {
		t.Fatalf("RemoveKeyword returned an unexpected error: %v", err)
	}

	keywords := note.UserKeywords(userID)
	if len(keywords) != 1 {
		t.Fatalf("Expected 1 keyword remaining, but got %d", len(keywords))
	}
	if keywords[0] != keyword2 {
		t.Errorf("Expected remaining keyword to be '%v', but got '%v'", keyword2, keywords[0])
	}
}

func TestNote_RemoveKeyword_UserNotFound(t *testing.T) {
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	userID := "user-1"
	keyword, _ := NewKeyword("keyword")
	note.AddKeyword(userID, keyword)

	err := note.RemoveKeyword("non-existent-user", keyword)
	if err != ErrUserNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrUserNotFound, err)
	}
}

func TestNote_RemoveKeyword_KeywordNotFound(t *testing.T) {
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	userID := "user-1"
	keyword, _ := NewKeyword("keyword")
	nonExistentKeyword, _ := NewKeyword("non-existent-keyword")
	note.AddKeyword(userID, keyword)

	err := note.RemoveKeyword(userID, nonExistentKeyword)
	if err != ErrKeywordNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrKeywordNotFound, err)
	}
}

func TestNote_AddCollaborator_Success(t *testing.T) {
	note, _ := NewNote("note1", "Test Note", "owner1")
	callerID := "owner1"
	collaboratorID := "user1"
	permission := ReadWrite

	err := note.AddCollaborator(callerID, collaboratorID, permission)

	if err != nil {
		t.Fatalf("AddCollaborator returned an unexpected error: %v", err)
	}
	if len(note.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(note.Collaborators))
	}
	if p, ok := note.Collaborators[collaboratorID]; !ok || p != permission {
		t.Errorf("Expected collaborator to have permission '%v', but got '%v'", permission, p)
	}
}

func TestNote_AddCollaborator_NotOwner(t *testing.T) {
	note, _ := NewNote("note1", "Test Note", "owner1")
	callerID := "another_user"
	collaboratorID := "user1"
	permission := ReadWrite

	err := note.AddCollaborator(callerID, collaboratorID, permission)

	if err == nil {
		t.Fatal("Expected an error when a non-owner adds a collaborator, but got nil")
	}
	if err != ErrPermissionDenied {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrPermissionDenied, err)
	}
	if len(note.Collaborators) != 0 {
		t.Fatalf("Expected 0 collaborators, but got %d", len(note.Collaborators))
	}
}

func TestNote_RemoveCollaborator_Success(t *testing.T) {
	note, _ := NewNote("note1", "Test Note", "owner1")
	ownerID := "owner1"
	collaboratorID := "user1"
	note.AddCollaborator(ownerID, collaboratorID, ReadOnly)
	keyword, _ := NewKeyword("test")
	note.AddKeyword(collaboratorID, keyword)

	err := note.RemoveCollaborator(ownerID, collaboratorID)

	if err != nil {
		t.Fatalf("RemoveCollaborator returned an unexpected error: %v", err)
	}
	if len(note.Collaborators) != 0 {
		t.Errorf("Expected 0 collaborators, but got %d", len(note.Collaborators))
	}
	if _, ok := note.keywords[collaboratorID]; ok {
		t.Errorf("Expected collaborator's keywords to be removed, but they were not")
	}
}

func TestNote_RemoveCollaborator_NotOwner(t *testing.T) {
	note, _ := NewNote("note1", "Test Note", "owner1")
	ownerID := "owner1"
	collaboratorID := "user1"
	nonOwnerID := "user2"
	note.AddCollaborator(ownerID, collaboratorID, ReadOnly)

	err := note.RemoveCollaborator(nonOwnerID, collaboratorID)

	if err != ErrPermissionDenied {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrPermissionDenied, err)
	}
	if len(note.Collaborators) != 1 {
		t.Errorf("Expected 1 collaborator, but got %d", len(note.Collaborators))
	}
}

func TestNote_RemoveCollaborator_CollaboratorNotFound(t *testing.T) {
	note, _ := NewNote("note1", "Test Note", "owner1")
	ownerID := "owner1"
	collaboratorID := "user1"
	note.AddCollaborator(ownerID, collaboratorID, ReadOnly)

	err := note.RemoveCollaborator(ownerID, "non-existent-user")

	if err != ErrUserNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrUserNotFound, err)
	}
	if len(note.Collaborators) != 1 {
		t.Errorf("Expected 1 collaborator, but got %d", len(note.Collaborators))
	}
}
