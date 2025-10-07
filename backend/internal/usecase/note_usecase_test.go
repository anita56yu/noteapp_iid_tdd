package usecase

import (
	"errors"
	"noteapp/internal/domain"
	"noteapp/internal/repository"
	"testing"
)

// mockNoteRepository is a mock implementation of the NoteRepository for testing error cases.
type mockNoteRepository struct {
	SaveFunc                 func(note *repository.NotePO) error
	FindByIDFunc             func(id string) (*repository.NotePO, error)
	DeleteFunc               func(id string) error
	FindByKeywordForUserFunc func(userID, keyword string) ([]*repository.NotePO, error)
	LockNoteForUpdateFunc    func(noteID string) error
	UnlockNoteForUpdateFunc  func(noteID string) error
}

func (m *mockNoteRepository) Save(note *repository.NotePO) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(note)
	}
	return nil
}
func (m *mockNoteRepository) FindByID(id string) (*repository.NotePO, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *mockNoteRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *mockNoteRepository) FindByKeywordForUser(userID, keyword string) ([]*repository.NotePO, error) {
	if m.FindByKeywordForUserFunc != nil {
		return m.FindByKeywordForUserFunc(userID, keyword)
	}
	return nil, nil
}
func (m *mockNoteRepository) LockNoteForUpdate(noteID string) error {
	if m.LockNoteForUpdateFunc != nil {
		return m.LockNoteForUpdateFunc(noteID)
	}
	return nil
}
func (m *mockNoteRepository) UnlockNoteForUpdate(noteID string) error {
	if m.UnlockNoteForUpdateFunc != nil {
		return m.UnlockNoteForUpdateFunc(noteID)
	}
	return nil
}

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
		SaveFunc: func(note *repository.NotePO) error {
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

func TestNoteUsecase_AddContent_WithInjectedID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	mapper := NewNoteMapper()
	id, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	content_id := "content-1"

	// Act
	newContentId, err := noteUsecase.AddContent(id, content_id, "New Content", TextContentType)

	// Assert
	if err != nil {
		t.Fatalf("AddContent() returned an unexpected error: %v", err)
	}
	if newContentId != content_id {
		t.Errorf("Expected content ID to be '%s', got '%s'", content_id, newContentId)
	}
	notePO, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	note := mapper.ToDomain(notePO)
	if len(note.Contents()) != 1 {
		t.Errorf("Expected 1 content block, got %d", len(note.Contents()))
	}
	if note.Contents()[0].Data != "New Content" {
		t.Errorf("Expected content to be 'New Content', got '%s'", note.Contents()[0].Data)
	}
}

func TestNoteUsecase_AddContent_WithGeneratedID(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	mapper := NewNoteMapper()
	id, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}

	// Act
	contentId, err := noteUsecase.AddContent(id, "", "New Content", TextContentType)

	// Assert
	if err != nil {
		t.Fatalf("AddContent() returned an unexpected error: %v", err)
	}
	notePO, err := repo.FindByID(id)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	note := mapper.ToDomain(notePO)
	if len(note.Contents()) != 1 {
		t.Errorf("Expected 1 content block, got %d", len(note.Contents()))
	}
	if note.Contents()[0].ID != contentId {
		t.Errorf("Expected content ID to be '%s', got '%s'", contentId, note.Contents()[0].ID)
	}
	if note.Contents()[0].Data != "New Content" {
		t.Errorf("Expected content to be 'New Content', got '%s'", note.Contents()[0].Data)
	}
}

func TestNoteUsecase_AddContent_NoteNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	_, err := noteUsecase.AddContent("non-existent-id", "", "New Content", TextContentType)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_GetNoteByID_WithMultipleContents(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	mapper := NewNoteMapper()
	note, err := domain.NewNote("note-1", "Test Title")
	if err != nil {
		t.Fatalf("NewNote() failed: %v", err)
	}
	note.AddContent("content-1", "Content 1", domain.TextContentType)
	note.AddContent("content-2", "Content 2", domain.ImageContentType)
	repo.Save(mapper.ToPO(note))

	// Act
	noteDTO, err := noteUsecase.GetNoteByID("note-1")

	// Assert
	if err != nil {
		t.Fatalf("GetNoteByID() returned an unexpected error: %v", err)
	}
	if len(noteDTO.Contents) != 2 {
		t.Errorf("Expected 2 content blocks, got %d", len(noteDTO.Contents))
	}
	if noteDTO.Contents[0].ID != "content-1" {
		t.Errorf("Expected content 1 ID to be 'content-1', got '%s'", noteDTO.Contents[0].ID)
	}
	if noteDTO.Contents[0].Data != "Content 1" {
		t.Errorf("Expected content 1 to be 'Content 1', got '%s'", noteDTO.Contents[0].Data)
	}
	if noteDTO.Contents[0].Type != "text" {
		t.Errorf("Expected content 1 type to be 'text', got '%s'", noteDTO.Contents[0].Type)
	}
	if noteDTO.Contents[1].ID != "content-2" {
		t.Errorf("Expected content 2 ID to be 'content-2', got '%s'", noteDTO.Contents[1].ID)
	}
	if noteDTO.Contents[1].Data != "Content 2" {
		t.Errorf("Expected content 2 to be 'Content 2', got '%s'", noteDTO.Contents[1].Data)
	}
	if noteDTO.Contents[1].Type != "image" {
		t.Errorf("Expected content 2 type to be 'image', got '%s'", noteDTO.Contents[1].Type)
	}
}

func TestNoteUsecase_UpdateContent_NormalCase(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	mapper := NewNoteMapper()
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	contentID, err := noteUsecase.AddContent(noteID, "", "Initial Content", TextContentType)
	if err != nil {
		t.Fatalf("AddContent() failed: %v", err)
	}
	updatedData := "Updated Content"

	// Act
	err = noteUsecase.UpdateContent(noteID, contentID, updatedData)

	// Assert
	if err != nil {
		t.Fatalf("UpdateContent() returned an unexpected error: %v", err)
	}
	notePO, err := repo.FindByID(noteID)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	note := mapper.ToDomain(notePO)
	contents := note.Contents()
	if len(contents) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(contents))
	}
	if contents[0].Data != updatedData {
		t.Errorf("Expected content to be '%s', got '%s'", updatedData, contents[0].Data)
	}
}

func TestNoteUsecase_UpdateContent_NoteNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	err := noteUsecase.UpdateContent("non-existent-note-id", "content-id", "Updated data")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_UpdateContent_ContentNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}

	// Act
	err = noteUsecase.UpdateContent(noteID, "non-existent-content-id", "Updated data")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for non-existent content, but got nil")
	}
	if !errors.Is(err, ErrContentNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrContentNotFound, err)
	}
}

func TestNoteUsecase_DeleteContent_Success(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	contentID1, err := noteUsecase.AddContent(noteID, "", "Content 1", TextContentType)
	if err != nil {
		t.Fatalf("AddContent() failed: %v", err)
	}
	contentID2, err := noteUsecase.AddContent(noteID, "", "Content 2", TextContentType)
	if err != nil {
		t.Fatalf("AddContent() failed: %v", err)
	}

	// Act
	err = noteUsecase.DeleteContent(noteID, contentID1)

	// Assert
	if err != nil {
		t.Fatalf("DeleteContent() returned an unexpected error: %v", err)
	}
	note, err := noteUsecase.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("GetNoteByID() failed: %v", err)
	}
	if len(note.Contents) != 1 {
		t.Errorf("Expected 1 content block, got %d", len(note.Contents))
	}
	if note.Contents[0].ID != contentID2 {
		t.Errorf("Expected remaining content ID to be '%s', got '%s'", contentID2, note.Contents[0].ID)
	}
	if note.Contents[0].Data != "Content 2" {
		t.Errorf("Expected content to be 'Content 2', got '%s'", note.Contents[0].Data)
	}
}

func TestNoteUsecase_DeleteContent_NoteNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)

	// Act
	err := noteUsecase.DeleteContent("non-existent-note-id", "content-id")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_DeleteContent_ContentNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}

	// Act
	err = noteUsecase.DeleteContent(noteID, "non-existent-content-id")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for non-existent content, but got nil")
	}
	if !errors.Is(err, ErrContentNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrContentNotFound, err)
	}
}

func TestNoteUsecase_TagNote(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	userID := "user-1"
	keyword := "test-keyword"

	// Act
	err = noteUsecase.TagNote(noteID, userID, keyword)

	// Assert
	if err != nil {
		t.Fatalf("TagNote() returned an unexpected error: %v", err)
	}
	notePO, err := repo.FindByID(noteID)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	if len(notePO.Keywords[userID]) != 1 {
		t.Fatalf("Expected 1 keyword, but got %d", len(notePO.Keywords[userID]))
	}
	if notePO.Keywords[userID][0] != keyword {
		t.Errorf("Expected keyword to be '%s', got '%s'", keyword, notePO.Keywords[userID][0])
	}
}

func TestNoteUsecase_TagNote_EmptyKeyword(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	userID := "user-1"

	// Act
	err = noteUsecase.TagNote(noteID, userID, "") // Empty keyword

	// Assert
	if err == nil {
		t.Fatal("Expected an error for empty keyword, but got nil")
	}
	if !errors.Is(err, ErrEmptyKeyword) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrEmptyKeyword, err)
	}
}

func TestNoteUsecase_FindNotesByKeyword(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	note1, _ := noteUsecase.CreateNote("", "Note 1")
	note2, _ := noteUsecase.CreateNote("", "Note 2")
	note3, _ := noteUsecase.CreateNote("", "Note 3")
	noteUsecase.TagNote(note1, "user-1", "go")
	noteUsecase.TagNote(note1, "user-1", "testing")
	noteUsecase.TagNote(note1, "user-2", "go")
	noteUsecase.TagNote(note2, "user-1", "testing")
	noteUsecase.TagNote(note3, "user-2", "java")
	noteUsecase.TagNote(note3, "user-2", "testing")

	// Act
	notes, err := noteUsecase.FindNotesByKeyword("user-1", "testing")

	// Assert
	if err != nil {
		t.Fatalf("FindNotesByKeyword() returned an unexpected error: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("Expected 2 notes, got %d", len(notes))
	}

	// Check that the correct notes are returned, regardless of order.
	noteIDs := make(map[string]bool)
	for _, note := range notes {
		noteIDs[note.ID] = true
	}
	if !noteIDs[note1] {
		t.Errorf("Expected to find note with ID %s", note1)
	}
	if !noteIDs[note2] {
		t.Errorf("Expected to find note with ID %s", note2)
	}
	if noteIDs[note3] {
		t.Errorf("Did not expect to find note with ID %s", note3)
	}
}

func TestNoteUsecase_FindNotesByKeyword_NoResults(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	note1, _ := noteUsecase.CreateNote("", "Note 1")
	noteUsecase.TagNote(note1, "user-1", "go")

	// Act
	notes, err := noteUsecase.FindNotesByKeyword("user-1", "testing")

	// Assert
	if err != nil {
		t.Fatalf("FindNotesByKeyword() returned an unexpected error: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("Expected 0 notes, got %d", len(notes))
	}
}

func TestNoteUsecase_FindNotesByKeyword_EmptyKeyword(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	note1, _ := noteUsecase.CreateNote("", "Note 1")
	noteUsecase.TagNote(note1, "user-1", "go")

	// Act
	notes, err := noteUsecase.FindNotesByKeyword("user-1", "")

	// Assert
	if err != nil {
		t.Fatalf("FindNotesByKeyword() returned an unexpected error: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("Expected 0 notes, got %d", len(notes))
	}
}

func TestNoteUsecase_UntagNote_Success(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	userID1 := "user-1"
	userID2 := "user-2"
	keywordToKeep := "go"
	keywordToRemove := "testing"
	noteUsecase.TagNote(noteID, userID1, keywordToKeep)
	noteUsecase.TagNote(noteID, userID1, keywordToRemove)
	noteUsecase.TagNote(noteID, userID2, keywordToKeep)
	noteUsecase.TagNote(noteID, userID2, keywordToRemove)

	// Act
	err = noteUsecase.UntagNote(noteID, userID1, keywordToRemove)

	// Assert
	if err != nil {
		t.Fatalf("UntagNote() returned an unexpected error: %v", err)
	}
	notePO, err := repo.FindByID(noteID)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	if len(notePO.Keywords[userID1]) != 1 {
		t.Fatalf("Expected 1 keyword for user1, but got %d", len(notePO.Keywords[userID1]))
	}
	if notePO.Keywords[userID1][0] != keywordToKeep {
		t.Errorf("Expected keyword for user1 to be '%s', got '%s'", keywordToKeep, notePO.Keywords[userID1][0])
	}
	if len(notePO.Keywords[userID2]) != 2 {
		t.Fatalf("Expected 2 keywords for user2, but got %d", len(notePO.Keywords[userID2]))
	}
}

func TestNoteUsecase_UntagNote_UserNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	userID := "user-1"
	keyword := "go"
	noteUsecase.TagNote(noteID, userID, keyword)

	// Act
	err = noteUsecase.UntagNote(noteID, "non-existent-user", keyword)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent user, but got nil")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrUserNotFound, err)
	}
}

func TestNoteUsecase_UntagNote_KeywordNotFound(t *testing.T) {
	// Arrange
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	noteID, err := noteUsecase.CreateNote("", "Test Title")
	if err != nil {
		t.Fatalf("CreateNote() failed: %v", err)
	}
	userID := "user-1"
	keyword := "go"
	noteUsecase.TagNote(noteID, userID, keyword)

	// Act
	err = noteUsecase.UntagNote(noteID, userID, "non-existent-keyword")

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent keyword, but got nil")
	}
	if !errors.Is(err, ErrKeywordNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrKeywordNotFound, err)
	}
}
