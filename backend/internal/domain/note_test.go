package domain

import (
	"testing"
)

func TestNewNote_ValidCreation_WithInjectedID(t *testing.T) {
	id := "test-id"
	title := "Test Note"

	note, err := NewNote(id, title)
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
	note, err := NewNote("", title)
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
	_, err := NewNote("", "")
	if err == nil {
		t.Fatal("Expected an error for empty title, but got nil")
	}
	if err != ErrEmptyTitle {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrEmptyTitle, err)
	}
}

func TestNote_AddContent_WithInjectedID(t *testing.T) {
	// Arrange
	note, _ := NewNote("note-1", "Test Note")
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
	note, _ := NewNote("note-1", "Test Note")
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
	note, _ := NewNote("note-1", "Test Note")
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
	note, _ := NewNote("note-1", "Test Note")
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
