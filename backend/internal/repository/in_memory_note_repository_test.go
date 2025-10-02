package repository

import (
	"errors"
	"testing"
)

func TestInMemoryNoteRepository_SaveAndFindByID_Success(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	note := &NotePO{ID: "test-id", Title: "Test Title"}

	// Act
	err := repo.Save(note)
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
	note := &NotePO{ID: "test-id", Title: "Original Title"}
	repo.Save(note)
	updatedNote := &NotePO{ID: "test-id", Title: "Updated Title"}

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
	note := &NotePO{ID: "test-id", Title: "Test Title"}
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

func TestInMemoryNoteRepository_FindByKeywordForUser(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	note1 := &NotePO{
		ID:    "note-1",
		Title: "Note 1",
		Keywords: map[string][]string{
			"user-1": {"go", "testing"},
			"user-2": {"go"},
		},
	}
	note2 := &NotePO{
		ID:    "note-2",
		Title: "Note 2",
		Keywords: map[string][]string{
			"user-1": {"testing"},
		},
	}
	note3 := &NotePO{
		ID:    "note-3",
		Title: "Note 3",
		Keywords: map[string][]string{
			"user-2": {"java", "testing"},
		},
	}
	repo.Save(note1)
	repo.Save(note2)
	repo.Save(note3)

	// Act
	notes, err := repo.FindByKeywordForUser("user-1", "testing")
	if err != nil {
		t.Fatalf("FindByKeywordForUser() returned an unexpected error: %v", err)
	}

	// Assert
	if len(notes) != 2 {
		t.Fatalf("Expected 2 notes, got %d", len(notes))
	}
	if notes[0].ID != "note-1" && notes[0].ID != "note-2" {
		t.Errorf("Expected note-1 or note-2, got %s", notes[0].ID)
	}
	if notes[1].ID != "note-1" && notes[1].ID != "note-2" {
		t.Errorf("Expected note-1 or note-2, got %s", notes[1].ID)
	}
}
