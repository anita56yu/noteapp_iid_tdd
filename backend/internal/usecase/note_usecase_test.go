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
	content := "Test Content"

	// Act
	returnedID, err := noteUsecase.CreateNote(id, title, content)
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
	content := "Test Content"

	// Act
	returnedID, err := noteUsecase.CreateNote("", title, content)
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
	_, err := noteUsecase.CreateNote("", "", "some content") // Empty title

	// Assert
	if err == nil {
		t.Fatal("Expected an error for empty title, but got nil")
	}
	expectedErr := "title cannot be empty"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
	}
}

func TestNoteUsecase_CreateNote_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := &mockNoteRepository{
		SaveFunc: func(note *domain.Note) error {
			return errors.New("database error")
		},
	}
	noteUsecase := NewNoteUsecase(mockRepo)

	// Act
	_, err := noteUsecase.CreateNote("test-id", "Test Title", "Test Content")

	// Assert
	if err == nil {
		t.Fatal("Expected a repository error, but got nil")
	}
	expectedErr := "database error"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
	}
}
