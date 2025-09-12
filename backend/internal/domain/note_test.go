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
