package usecase

import (
	"noteapp/internal/domain"
	"testing"
)

func TestToNoteDTO(t *testing.T) {
	// Arrange
	note, err := domain.NewNote("note-1", "Test Title")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}
	note.AddContent("content-1", "Hello", domain.TextContentType)
	note.AddContent("content-2", "base64-encoded-image", domain.ImageContentType)

	// Act
	noteDTO := toNoteDTO(note)

	// Assert
	if noteDTO.ID != "note-1" {
		t.Errorf("Expected DTO ID to be 'note-1', got '%s'", noteDTO.ID)
	}
	if noteDTO.Title != "Test Title" {
		t.Errorf("Expected DTO title to be 'Test Title', got '%s'", noteDTO.Title)
	}
	if len(noteDTO.Contents) != 2 {
		t.Fatalf("Expected 2 content DTOs, got %d", len(noteDTO.Contents))
	}

	// Check content 1
	if noteDTO.Contents[0].ID != "content-1" {
		t.Errorf("Expected content 1 ID to be 'content-1', got '%s'", noteDTO.Contents[0].ID)
	}
	if noteDTO.Contents[0].Type != "text" {
		t.Errorf("Expected content 1 type to be 'text', got '%s'", noteDTO.Contents[0].Type)
	}
	if noteDTO.Contents[0].Data != "Hello" {
		t.Errorf("Expected content 1 data to be 'Hello', got '%s'", noteDTO.Contents[0].Data)
	}

	// Check content 2
	if noteDTO.Contents[1].ID != "content-2" {
		t.Errorf("Expected content 2 ID to be 'content-2', got '%s'", noteDTO.Contents[1].ID)
	}
	if noteDTO.Contents[1].Type != "image" {
		t.Errorf("Expected content 2 type to be 'image', got '%s'", noteDTO.Contents[1].Type)
	}
	if noteDTO.Contents[1].Data != "base64-encoded-image" {
		t.Errorf("Expected content 2 data to be 'base64-encoded-image', got '%s'", noteDTO.Contents[1].Data)
	}
}
