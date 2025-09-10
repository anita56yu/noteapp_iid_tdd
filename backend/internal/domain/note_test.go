package domain

import (
	"testing"
)

func TestNewNote_ValidCreation_WithInjectedID(t *testing.T) {
	id := "test-id"
	title := "Test Note"
	content := "This is a test note."
	note, err := NewNote(id, title, content)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}
	if note.ID != id {
		t.Errorf("Expected ID to be '%s', but got '%s'", id, note.ID)
	}
	if note.Title != title {
		t.Errorf("Expected title to be '%s', but got '%s'", title, note.Title)
	}
	if note.Content != content {
		t.Errorf("Expected content to be '%s', but got '%s'", content, note.Content)
	}
}

func TestNewNote_ValidCreation_WithGeneratedID(t *testing.T) {
	title := "Test Note"
	content := "This is a test note."
	note, err := NewNote("", title, content)
	if err != nil {
		t.Fatalf("Failed to create a valid note: %v", err)
	}
	if note.ID == "" {
		t.Errorf("Expected ID to be non-empty, but it was empty")
	}
	if note.Title != title {
		t.Errorf("Expected title to be '%s', but got '%s'", title, note.Title)
	}
	if note.Content != content {
		t.Errorf("Expected content to be '%s', but got '%s'", content, note.Content)
	}
}

func TestNewNote_EmptyTitle(t *testing.T) {
	_, err := NewNote("", "", "some content")
	if err == nil {
		t.Fatal("Expected an error for empty title, but got nil")
	}
	if err != ErrEmptyTitle {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrEmptyTitle, err)
	}
}
