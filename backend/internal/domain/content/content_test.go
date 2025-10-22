package content_test

import (
	"noteapp/internal/domain/content"
	"testing"
)

func TestNewContent_WithID(t *testing.T) {
	id := "test-id"
	noteID := "test-note-id"
	c := content.NewContent(id, noteID, "Test Content", content.TextContentType, 0)

	if c.ID != id {
		t.Errorf("Expected ID to be %v, but got %v", id, c.ID)
	}
	if c.NoteID != noteID {
		t.Errorf("Expected NoteID to be %v, but got %v", noteID, c.NoteID)
	}
	if c.Data != "Test Content" {
		t.Errorf("Expected Data to be 'Test Content', but got '%s'", c.Data)
	}
	if c.Type != content.TextContentType {
		t.Errorf("Expected Type to be %v, but got %v", content.TextContentType, c.Type)
	}
	if c.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", c.Version)
	}
}

func TestNewContent_GenerateID(t *testing.T) {
	noteID := "test-note-id"
	c := content.NewContent("", noteID, "Test Content", content.TextContentType, 0)

	if c.ID == "" {
		t.Error("Expected ID to be generated, but it is empty")
	}
	if c.NoteID != noteID {
		t.Errorf("Expected NoteID to be %v, but got %v", noteID, c.NoteID)
	}
	if c.Data != "Test Content" {
		t.Errorf("Expected Data to be 'Test Content', but got '%s'", c.Data)
	}
	if c.Type != content.TextContentType {
		t.Errorf("Expected Type to be %v, but got %v", content.TextContentType, c.Type)
	}
	if c.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", c.Version)
	}
}
