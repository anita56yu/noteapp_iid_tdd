package repository

import (
	"errors"
	"noteapp/internal/domain"
	"testing"
)

func TestInMemoryNoteRepository_SaveAndFindByID_Success(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	note, err := domain.NewNote("test-id", "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("Failed to create a new note for testing: %v", err)
	}

	// Act
	err = repo.Save(note)
	if err != nil {
		t.Fatalf("Save() returned an unexpected error: %v", err)
	}

	// Assert
	foundNote, err := repo.FindByID("test-id")
	if err != nil {
		t.Fatalf("FindByID() returned an unexpected error: %v", err)
	}
	if foundNote == nil {
		t.Fatal("FindByID() returned nil, expected a note")
	}
	if foundNote.ID != note.ID {
		t.Errorf("Expected ID %s, got %s", note.ID, foundNote.ID)
	}
	if foundNote.Title != note.Title {
		t.Errorf("Expected Title %s, got %s", note.Title, foundNote.Title)
	}
}

func TestInMemoryNoteRepository_FindByID_NotFound(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()

	// Act
	_, err := repo.FindByID("non-existent-id")

	// Assert
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error %v, got %v", ErrNoteNotFound, err)
	}
}

func TestInMemoryNoteRepository_Save_UpdateExisting(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	note, _ := domain.NewNote("test-id", "Original Title", "Original Content")
	repo.Save(note)
	updatedNote, _ := domain.NewNote("test-id", "Updated Title", "Updated Content")

	// Act
	err := repo.Save(updatedNote)
	if err != nil {
		t.Fatalf("Save() returned an unexpected error on update: %v", err)
	}

	// Assert
	foundNote, _ := repo.FindByID("test-id")
	if foundNote.Title != "Updated Title" {
		t.Errorf("Expected updated title, got '%s'", foundNote.Title)
	}
}

func TestInMemoryNoteRepository_Save_NilNote(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()

	// Act
	err := repo.Save(nil)

	// Assert
	if !errors.Is(err, ErrNilNote) {
		t.Errorf("Expected error %v, got %v", ErrNilNote, err)
	}
}

func TestInMemoryNoteRepository_Delete(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	note, _ := domain.NewNote("test-id", "Test Title", "Test Content")
	repo.Save(note)

	// Act
	err := repo.Delete("test-id")
	if err != nil {
		t.Fatalf("Delete() returned an unexpected error: %v", err)
	}

	// Assert
	_, err = repo.FindByID("test-id")
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error %v after delete, got %v", ErrNoteNotFound, err)
	}
}

func TestInMemoryNoteRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()

	// Act
	err := repo.Delete("non-existent-id")

	// Assert
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error %v, got %v", ErrNoteNotFound, err)
	}
}
