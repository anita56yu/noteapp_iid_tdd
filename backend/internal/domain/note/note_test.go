package note

import (
	"testing"
)

func TestNewNote_ValidCreation_WithInjectedID(t *testing.T) {
	id := "test-id"
	title := "Test Note"
	ownerID := "owner-1"

	note, err := NewNote(id, title, ownerID)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}

	if note.ID != id {
		t.Errorf("Expected ID to be '%s', but got '%s'", id, note.ID)
	}

	if note.contents == nil {
		t.Fatalf("Expected contents to be an empty slice, but it was nil")
	}

	if len(note.contents) != 0 {
		t.Errorf("Expected contents to be empty, but got %d elements", len(note.contents))
	}
}

func TestNewNote_ValidCreation_WithGeneratedID(t *testing.T) {
	title := "Test Note"
	ownerID := "owner-1"
	note, err := NewNote("", title, ownerID)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}
	if note.ID == "" {
		t.Errorf("Expected ID to be non-empty, but it was empty")
	}
	if note.Title != title {
		t.Errorf("Expected title to be '%s', but got '%s'", title, note.Title)
	}
	if note.contents == nil {
		t.Fatalf("Expected contents to be an empty slice, but it was nil")
	}

	if len(note.contents) != 0 {
		t.Errorf("Expected contents to be empty, but got %d elements", len(note.contents))
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

func TestNote_AddContent_WithInjectedID(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	contentID := "content-1"
	contentData := "Hello, world!"

	// Act
	newid := note.AddContent(contentID, contentData, TextContentType)

	// Assert
	contents := note.Contents()
	if newid != contentID {
		t.Errorf("Expected returned content ID to be '%s', but got '%s'", contentID, newid)
	}
	if len(contents) != 1 {
		t.Fatalf("Expected 1 content block, but got %d", len(contents))
	}
	if contents[0].ID != contentID {
		t.Errorf("Expected content ID to be '%s', but got '%s'", contentID, contents[0].ID)
	}
	if contents[0].Data != contentData {
		t.Errorf("Expected content data to be '%s', but got '%s'", contentData, contents[0].Data)
	}
}

func TestNote_AddContent_WithGeneratedID(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	contentData := "Hello, world!"

	// Act
	note.AddContent("", contentData, TextContentType)

	// Assert
	contents := note.Contents()
	if len(contents) != 1 {
		t.Fatalf("Expected 1 content block, but got %d", len(contents))
	}
	if contents[0].ID == "" {
		t.Error("Expected content ID to be non-empty")
	}
	if contents[0].Data != contentData {
		t.Errorf("Expected content data to be '%s', but got '%s'", contentData, contents[0].Data)
	}
}

func TestNote_UpdateContent_NormalCase(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	contentID := "content-1"
	initialData := "Initial content"
	note.AddContent(contentID, initialData, TextContentType)
	updatedData := "Updated content"

	// Act
	err := note.UpdateContent(contentID, updatedData)

	// Assert
	if err != nil {
		t.Fatalf("UpdateContent returned an unexpected error: %v", err)
	}
	contents := note.Contents()
	if contents[0].Data != updatedData {
		t.Errorf("Expected content data to be '%s', but got '%s'", updatedData, contents[0].Data)
	}
}

func TestNote_UpdateContent_NotFound(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	note.AddContent("content-1", "Initial data", TextContentType)

	// Act
	err := note.UpdateContent("non-existent-id", "Updated data")

	// Assert
	if err == nil {
		t.Fatal("Expected an error when updating non-existent content, but got nil")
	}
	if err != ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrContentNotFound, err)
	}
}

func TestNote_DeleteContent_Success(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	contentID1 := note.AddContent("", "Content 1", TextContentType)
	contentID2 := note.AddContent("", "Content 2", TextContentType)

	// Act
	err := note.DeleteContent(contentID1)

	// Assert
	if err != nil {
		t.Fatalf("DeleteContent returned an unexpected error: %v", err)
	}
	contents := note.Contents()
	if len(contents) != 1 {
		t.Fatalf("Expected 1 content block after deletion, but got %d", len(contents))
	}
	if contents[0].ID != contentID2 {
		t.Errorf("Expected remaining content ID to be '%s', but got '%s'", contentID2, contents[0].ID)
	}
}

func TestNote_DeleteContent_NotFound(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note", "owner-1")
	note.AddContent("content-1", "Initial data", TextContentType)

	// Act
	err := note.DeleteContent("non-existent-id")

	// Assert
	if err == nil {
		t.Fatal("Expected an error when deleting non-existent content, but got nil")
	}
	if err != ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrContentNotFound, err)
	}
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
