package usecase

import (
	"errors"
	"noteapp/internal/domain"
	"noteapp/internal/repository"
	"testing"
)

// mockNoteRepository is a mock implementation of the NoteRepository for testing error cases.
type mockNoteRepository struct {
	SaveFunc func(note *domain.Note) error
}

func (m *mockNoteRepository) Save(note *domain.Note) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(note)
	}
	return nil
}
func (m *mockNoteRepository) FindByID(id string) (*domain.Note, error) { return nil, nil }
func (m *mockNoteRepository) Delete(id string) error                   { return nil }

func TestNoteUsecase_CreateNote_WithInjectedID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	id := "test-id"
	title := "Test Title"

	// Act
	returnedID, err := noteUsecase.CreateNote(id, title)
	if err != nil {
		t.Fatalf("CreateNote() returned an unexpected error: %v", err)
	}

	// Assert
	if returnedID != id {
		t.Errorf("Expected returned ID to be '%s', got '%s'", id, returnedID)
	}
	savedNote, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("Failed to find saved note: %v", err)
	}
	if savedNote.Title != title {
		t.Errorf("Expected saved note title to be '%s', got '%s'", title, savedNote.Title)
	}
}

func TestNoteUsecase_CreateNote_WithGeneratedID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	title := "Test Title"

	// Act
	returnedID, err := noteUsecase.CreateNote("", title)
	if err != nil {
		t.Fatalf("CreateNote() returned an unexpected error: %v", err)
	}

	// Assert
	if returnedID == "" {
		t.Error("Expected a generated ID, but got an empty string")
	}
	savedNote, err := repo.FindByID(returnedID)
	if err != nil {
		t.Fatalf("Failed to find saved note with generated ID: %v", err)
	}
	if savedNote.Title != title {
		t.Errorf("Expected saved note title to be '%s', got '%s'", title, savedNote.Title)
	}
}

func TestNoteUsecase_CreateNote_DomainError(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	_, err := noteUsecase.CreateNote("", "") // Empty title

	// Assert
	if err == nil {
		t.Fatal("Expected an error for empty title, but got nil")
	}
	if !errors.Is(err, ErrEmptyTitle) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrEmptyTitle, err)
	}
}

func TestNoteUsecase_CreateNote_NilNoteError(t *testing.T) {
	// Arrange
	mockRepo := &mockNoteRepository{
		SaveFunc: func(note *domain.Note) error {
			return repository.ErrNilNote
		},
	}
	noteUsecase := NewNoteUsecase(mockRepo)

	// Act
	_, err := noteUsecase.CreateNote("test-id", "Test Title")

	// Assert
	if err == nil {
		t.Fatal("Expected a repository error, but got nil")
	}

	if err != ErrNilNote {
		t.Errorf("Expected error message '%s', but got '%s'", ErrNilNote, err)
	}
}

func TestNoteUsecase_GetNoteByID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	id, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}

	// Act
	noteDTO, err := noteUsecase.GetNoteByID(id)
	if err != nil {
		t.Fatalf("GetNoteByID() returned an unexpected error: %v", err)
	}

	// Assert
	if noteDTO == nil {
		t.Fatal("Expected a note DTO, but got nil")
	}
	if noteDTO.ID != id {
		t.Errorf("Expected note ID to be '%s', got '%s'", id, noteDTO.ID)
	}
	if noteDTO.Title != "Test Title" {
		t.Errorf("Expected note title to be '%s', got '%s'", "Test Title", noteDTO.Title)
	}
}

func TestNoteUsecase_GetNoteByID_NotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	_, err := noteUsecase.GetNoteByID("non-existent-id")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_GetNoteByID_InvalidID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	_, err := noteUsecase.GetNoteByID("") // Empty ID

	// Assert
	if err == nil {
		t.Fatal("Expected an error for an invalid ID, but got nil")
	}
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrInvalidID, err)
	}
}

func TestNoteUsecase_DeleteNote_Success(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	id, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}

	// Act
	err = noteUsecase.DeleteNote(id)
	if err != nil {
		t.Fatalf("DeleteNote() returned an unexpected error: %v", err)
	}

	// Assert
	_, err = repo.FindByID(id)
	if !errors.Is(err, repository.ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", repository.ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_DeleteNote_NotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	err := noteUsecase.DeleteNote("non-existent-id")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_DeleteNote_InvalidID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	err := noteUsecase.DeleteNote("") // Empty ID

	// Assert
	if err == nil {
		t.Fatal("Expected an error for an invalid ID, but got nil")
	}
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrInvalidID, err)
	}
}
