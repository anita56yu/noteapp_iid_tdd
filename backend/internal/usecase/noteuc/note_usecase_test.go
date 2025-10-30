package noteuc

import (
	"errors"
	"noteapp/internal/repository/noterepo"
	"testing"
)

// mockNoteRepository is a mock implementation of the NoteRepository for testing error cases.
type mockNoteRepository struct {
	SaveFunc                       func(note *noterepo.NotePO) error
	FindByIDFunc                   func(id string) (*noterepo.NotePO, error)
	DeleteFunc                     func(id string) error
	FindByKeywordForUserFunc       func(userID, keyword string) ([]*noterepo.NotePO, error)
	GetAccessibleNotesByUserIDFunc func(userID string) ([]*noterepo.NotePO, error)
}

func (m *mockNoteRepository) Save(note *noterepo.NotePO) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(note)
	}
	return nil
}
func (m *mockNoteRepository) FindByID(id string) (*noterepo.NotePO, error) {
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
func (m *mockNoteRepository) FindByKeywordForUser(userID, keyword string) ([]*noterepo.NotePO, error) {
	if m.FindByKeywordForUserFunc != nil {
		return m.FindByKeywordForUserFunc(userID, keyword)
	}
	return nil, nil
}
func (m *mockNoteRepository) GetAccessibleNotesByUserID(userID string) ([]*noterepo.NotePO, error) {
	if m.GetAccessibleNotesByUserIDFunc != nil {
		return m.GetAccessibleNotesByUserIDFunc(userID)
	}
	return nil, nil
}

func setUpRepositoryAndUsecase() (*noterepo.InMemoryNoteRepository, *NoteUsecase) {
	repo := noterepo.NewInMemoryNoteRepository()
	noteUsecase := NewNoteUsecase(repo)
	return repo, noteUsecase
}

func setUpRepositoryAndUsecaseWithNote() (*noterepo.InMemoryNoteRepository, *NoteUsecase, string) {
	repo, noteUsecase := setUpRepositoryAndUsecase()
	noteID, _ := noteUsecase.CreateNote("", "Test Title", "owner-1")
	return repo, noteUsecase, noteID
}

func setUpRepositoryAndUsecaseWithNoteAndContents() (*noterepo.InMemoryNoteRepository, *NoteUsecase, string, string, string) {
	repo, noteUsecase := setUpRepositoryAndUsecase()
	noteID, _ := noteUsecase.CreateNote("", "Test Title", "owner-1")
	contentID := "content-1"
	contentID1 := "content-2"
	noteUsecase.AddContent(noteID, contentID, 0)
	noteUsecase.AddContent(noteID, contentID1, 1)
	return repo, noteUsecase, noteID, contentID, contentID1
}

func setUpRepositoryAndUsecaseWithTaggedNotes() (*noterepo.InMemoryNoteRepository, *NoteUsecase, string, string, string) {
	repo, noteUsecase := setUpRepositoryAndUsecase()
	note1, _ := noteUsecase.CreateNote("", "Note 1", "owner-1")
	note2, _ := noteUsecase.CreateNote("", "Note 2", "owner-2")
	note3, _ := noteUsecase.CreateNote("", "Note 3", "owner-3")
	noteUsecase.ShareNote(note2, "owner-2", "user-1", "read", 0)
	noteUsecase.ShareNote(note2, "owner-2", "user-2", "read", 1)
	noteUsecase.TagNote(note1, "user-1", "go", 0)
	noteUsecase.TagNote(note1, "user-1", "testing", 1)
	noteUsecase.TagNote(note1, "user-2", "go", 2)
	noteUsecase.TagNote(note2, "user-1", "testing", 2)
	noteUsecase.TagNote(note2, "user-2", "java", 3)
	noteUsecase.TagNote(note3, "user-2", "java", 0)
	noteUsecase.TagNote(note3, "user-2", "testing", 1)

	return repo, noteUsecase, note1, note2, note3
}

func TestNoteUsecase_CreateNote_WithInjectedID(t *testing.T) {
	// Arrange
	repo, noteUsecase := setUpRepositoryAndUsecase()
	id := "test-id"
	title := "Test Title"
	ownerID := "owner-1"

	// Act
	returnedID, err := noteUsecase.CreateNote(id, title, ownerID)
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
	if savedNote.OwnerID != ownerID {
		t.Errorf("Expected saved note owner ID to be '%s', got '%s'", ownerID, savedNote.OwnerID)
	}
}

func TestNoteUsecase_CreateNote_WithGeneratedID(t *testing.T) {
	// Arrange
	repo, noteUsecase := setUpRepositoryAndUsecase()
	title := "Test Title"
	ownerID := "owner-1"

	// Act
	returnedID, err := noteUsecase.CreateNote("", title, ownerID)
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
	if savedNote.OwnerID != ownerID {
		t.Errorf("Expected saved note owner ID to be '%s', got '%s'", ownerID, savedNote.OwnerID)
	}
}

func TestNoteUsecase_CreateNote_DomainError(t *testing.T) {
	// Arrange
	_, noteUsecase := setUpRepositoryAndUsecase()

	// Act
	_, err := noteUsecase.CreateNote("", "", "owner-1") // Empty title

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
		SaveFunc: func(note *noterepo.NotePO) error {
			return noterepo.ErrNilNote
		},
	}
	noteUsecase := NewNoteUsecase(mockRepo)

	// Act
	_, err := noteUsecase.CreateNote("test-id", "Test Title", "owner-1")

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
	_, noteUsecase, id := setUpRepositoryAndUsecaseWithNote()

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
	_, noteUsecase, _ := setUpRepositoryAndUsecaseWithNote()

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
	_, noteUsecase, _ := setUpRepositoryAndUsecaseWithNote()

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
	repo, noteUsecase, id := setUpRepositoryAndUsecaseWithNote()

	// Act
	err := noteUsecase.DeleteNote(id, 0)
	if err != nil {
		t.Fatalf("DeleteNote() returned an unexpected error: %v", err)
	}

	// Assert
	_, err = repo.FindByID(id)
	if !errors.Is(err, noterepo.ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", noterepo.ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_DeleteNote_NotFound(t *testing.T) {
	// Arrange
	_, noteUsecase := setUpRepositoryAndUsecase()

	// Act
	err := noteUsecase.DeleteNote("non-existent-id", 0)

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
	_, noteUsecase, _ := setUpRepositoryAndUsecaseWithNote()

	// Act
	err := noteUsecase.DeleteNote("", 0) // Empty ID

	// Assert
	if err == nil {
		t.Fatal("Expected an error for an invalid ID, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_DeleteNote_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()

	// Act
	err := noteUsecase.DeleteNote(noteID, 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}

func TestNoteUsecase_GetNoteByID_WithMultipleContents(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, contentID, contentID1 := setUpRepositoryAndUsecaseWithNoteAndContents()

	// Act
	noteDTO, err := noteUsecase.GetNoteByID(noteID)

	// Assert
	if err != nil {
		t.Fatalf("GetNoteByID() returned an unexpected error: %v", err)
	}
	if len(noteDTO.ContentIDs) != 2 {
		t.Errorf("Expected 2 content blocks, got %d", len(noteDTO.ContentIDs))
	}
	if noteDTO.ContentIDs[0] != contentID {
		t.Errorf("Expected content 1 ID to be '%s', got '%s'", contentID, noteDTO.ContentIDs[0])
	}
	if noteDTO.ContentIDs[1] != contentID1 {
		t.Errorf("Expected content 2 ID to be '%s', got '%s'", contentID1, noteDTO.ContentIDs[1])
	}
}

func TestNoteUsecase_RemoveContent_Success(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, contentID1, contentID2 := setUpRepositoryAndUsecaseWithNoteAndContents()

	// Act
	err := noteUsecase.RemoveContent(noteID, contentID1, 2)

	// Assert
	if err != nil {
		t.Fatalf("RemoveContent() returned an unexpected error: %v", err)
	}
	note, err := noteUsecase.GetNoteByID(noteID)
	if err != nil {
		t.Fatalf("GetNoteByID() failed: %v", err)
	}
	if len(note.ContentIDs) != 1 {
		t.Errorf("Expected 1 content ID, got %d", len(note.ContentIDs))
	}
	if note.ContentIDs[0] != contentID2 {
		t.Errorf("Expected remaining content ID to be '%s', got '%s'", contentID2, note.ContentIDs[0])
	}
}

func TestNoteUsecase_RemoveContent_NoteNotFound(t *testing.T) {
	// Arrange
	_, noteUsecase, _, _, _ := setUpRepositoryAndUsecaseWithNoteAndContents()

	// Act
	err := noteUsecase.RemoveContent("non-existent-note-id", "content-id", 0)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_RemoveContent_ContentNotFound(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithNoteAndContents()

	// Act
	err := noteUsecase.RemoveContent(noteID, "non-existent-content-id", 2)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for non-existent content, but got nil")
	}
	if !errors.Is(err, ErrContentNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrContentNotFound, err)
	}
}

func TestNoteUsecase_RemoveContent_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, contentID, _ := setUpRepositoryAndUsecaseWithNoteAndContents()

	// Act
	err := noteUsecase.RemoveContent(noteID, contentID, 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}

func TestNoteUsecase_TagNote(t *testing.T) {
	// Arrange
	repo, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	userID := "user-1"
	keyword := "test-keyword"

	// Act
	err := noteUsecase.TagNote(noteID, userID, keyword, 0)

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
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	userID := "user-1"

	// Act
	err := noteUsecase.TagNote(noteID, userID, "", 0) // Empty keyword

	// Assert
	if err == nil {
		t.Fatal("Expected an error for empty keyword, but got nil")
	}
	if !errors.Is(err, ErrEmptyKeyword) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrEmptyKeyword, err)
	}
}

func TestNoteUsecase_TagNote_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	userID := "user-1"
	keyword := "test-keyword"

	// Act
	err := noteUsecase.TagNote(noteID, userID, keyword, 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}

func TestNoteUsecase_FindNotesByKeyword(t *testing.T) {
	// Arrange
	_, noteUsecase, note1, note2, note3 := setUpRepositoryAndUsecaseWithTaggedNotes()

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
	_, noteUsecase, _, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()

	// Act
	notes, err := noteUsecase.FindNotesByKeyword("user-1", "java")

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
	_, noteUsecase, _, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()

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
	repo, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	keywordToRemove := "go"
	keywordToKeep := "testing"
	userID1 := "user-1"
	userID2 := "user-2"

	// Act
	err := noteUsecase.UntagNote(noteID, userID1, keywordToRemove, 3)

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
	if len(notePO.Keywords[userID2]) != 1 {
		t.Fatalf("Expected 1 keyword for user2, but got %d", len(notePO.Keywords[userID2]))
	}
}

func TestNoteUsecase_UntagNote_UserNotFound(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	keyword := "go"

	// Act
	err := noteUsecase.UntagNote(noteID, "non-existent-user", keyword, 3)

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
	_, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	userID := "user-1"

	// Act
	err := noteUsecase.UntagNote(noteID, userID, "non-existent-keyword", 3)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent keyword, but got nil")
	}
	if !errors.Is(err, ErrKeywordNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrKeywordNotFound, err)
	}
}

func TestNoteUsecase_UntagNote_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	userID := "user-1"
	keyword := "go"

	// Act
	err := noteUsecase.UntagNote(noteID, userID, keyword, 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}

func TestNoteUsecase_ShareNote_NotOwner(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	collaboratorID := "collaborator-1"

	// Act
	err := noteUsecase.ShareNote(noteID, "not-the-owner", collaboratorID, "read", 0)

	// Assert
	if err == nil {
		t.Fatal("Expected an error when sharing a note by a non-owner, but got nil")
	}
	if !errors.Is(err, ErrPermissionDenied) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrPermissionDenied, err)
	}
}

func TestNoteUsecase_ShareNote_InvalidPermission(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	collaboratorID := "collaborator-1"
	ownerID := "owner-1"

	// Act
	err := noteUsecase.ShareNote(noteID, ownerID, collaboratorID, "invalid-permission", 0)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for invalid permission, but got nil")
	}
}

func TestNoteUsecase_ShareNote_Success(t *testing.T) {
	// Arrange
	repo, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	collaboratorID := "collaborator-1"
	ownerID := "owner-1"

	// Act
	err := noteUsecase.ShareNote(noteID, ownerID, collaboratorID, "read", 0)

	// Assert
	if err != nil {
		t.Fatalf("ShareNote() returned an unexpected error: %v", err)
	}
	notePO, err := repo.FindByID(noteID)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	if len(notePO.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(notePO.Collaborators))
	}
	if _, ok := notePO.Collaborators[collaboratorID]; !ok {
		t.Errorf("Expected collaborator with ID '%s' to be in the map", collaboratorID)
	}
	if notePO.Collaborators[collaboratorID] != "read" {
		t.Errorf("Expected collaborator permission to be 'read', got '%s'", notePO.Collaborators[collaboratorID])
	}
}

func TestNoteUsecase_ShareNote_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	collaboratorID := "collaborator-1"
	ownerID := "owner-1"

	// Act
	err := noteUsecase.ShareNote(noteID, ownerID, collaboratorID, "read", 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}

func TestNoteUsecase_GetAccessibleNotesForUser(t *testing.T) {
	// Arrange
	_, noteUsecase, ownedNoteID, sharedNoteID, unrelatedNoteID := setUpRepositoryAndUsecaseWithTaggedNotes()
	noteUsecase.ShareNote(sharedNoteID, "owner-2", "owner-1", "read", 4)

	// Act
	notes, err := noteUsecase.GetAccessibleNotesForUser("owner-1")

	// Assert
	if err != nil {
		t.Fatalf("GetAccessibleNotesForUser() returned an unexpected error: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("Expected 2 accessible notes, but got %d", len(notes))
	}

	// Check that the correct notes are returned
	returnedIDs := make(map[string]bool)
	for _, note := range notes {
		returnedIDs[note.ID] = true
	}
	if !returnedIDs[ownedNoteID] {
		t.Errorf("Expected to find owned note with ID %s", ownedNoteID)
	}
	if !returnedIDs[sharedNoteID] {
		t.Errorf("Expected to find shared note with ID %s", sharedNoteID)
	}
	if returnedIDs[unrelatedNoteID] {
		t.Errorf("Did not expect to find unrelated note with ID %s", unrelatedNoteID)
	}
}

func TestNoteUsecase_RevokeAccess_Success(t *testing.T) {
	// Arrange
	_, noteUsecase, _, noteID, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	ownerID := "owner-2"
	collaboratorID1 := "user-1"
	collaboratorID2 := "user-2"

	// Act
	err := noteUsecase.RevokeAccess(noteID, ownerID, collaboratorID1, 4)

	// Assert
	if err != nil {
		t.Fatalf("RevokeAccess() returned an unexpected error: %v", err)
	}
	note, _ := noteUsecase.GetNoteByID(noteID)
	if _, ok := note.Collaborators[collaboratorID1]; ok {
		t.Errorf("Expected collaborator 1 to be removed, but they still exist")
	}
	if _, ok := note.Collaborators[collaboratorID2]; !ok {
		t.Errorf("Expected collaborator 2 to remain, but they were removed")
	}
	if _, ok := note.Keywords[collaboratorID1]; ok {
		t.Errorf("Expected collaborator's keywords to be removed, but they still exist")
	}
	if _, ok := note.Keywords[collaboratorID2]; !ok {
		t.Errorf("Expected collaborator 2's keywords to remain, but they were removed")
	}
}

func TestNoteUsecase_RevokeAccess_NotOwner(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	nonOwnerID := "user-1"
	collaboratorID := "user-2"

	// Act
	err := noteUsecase.RevokeAccess(noteID, nonOwnerID, collaboratorID, 3)

	// Assert
	if err != ErrPermissionDenied {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrPermissionDenied, err)
	}
}

func TestNoteUsecase_RevokeAccess_CollaboratorNotFound(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID, _, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	ownerID := "owner-1"

	// Act
	err := noteUsecase.RevokeAccess(noteID, ownerID, "non-existent-user", 3)

	// Assert
	if err != ErrUserNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrUserNotFound, err)
	}
}

func TestNoteUsecase_RevokeAccess_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, _, noteID, _ := setUpRepositoryAndUsecaseWithTaggedNotes()
	ownerID := "owner-2"
	collaboratorID := "user-1"

	// Act
	err := noteUsecase.RevokeAccess(noteID, ownerID, collaboratorID, 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}

func TestNoteUsecase_AddContent_Success(t *testing.T) {
	// Arrange
	repo, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	contentID := "new-content-id"

	// Act
	err := noteUsecase.AddContent(noteID, contentID, 0)

	// Assert
	if err != nil {
		t.Fatalf("AddContent() returned an unexpected error: %v", err)
	}
	notePO, err := repo.FindByID(noteID)
	if err != nil {
		t.Fatalf("Failed to find note: %v", err)
	}
	if len(notePO.ContentIDs) != 1 {
		t.Fatalf("Expected 1 content ID, but got %d", len(notePO.ContentIDs))
	}
	if notePO.ContentIDs[0] != contentID {
		t.Errorf("Expected content ID to be '%s', but got '%s'", contentID, notePO.ContentIDs[0])
	}
}

func TestNoteUsecase_AddContent_NoteNotFound(t *testing.T) {
	// Arrange
	_, noteUsecase, _ := setUpRepositoryAndUsecaseWithNote()
	contentID := "new-content-id"

	// Act
	err := noteUsecase.AddContent("non-existent-id", contentID, 0)

	// Assert
	if err == nil {
		t.Fatal("Expected an error for a non-existent note, but got nil")
	}
	if !errors.Is(err, ErrNoteNotFound) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrNoteNotFound, err)
	}
}

func TestNoteUsecase_AddContent_Conflict(t *testing.T) {
	// Arrange
	_, noteUsecase, noteID := setUpRepositoryAndUsecaseWithNote()
	contentID := "new-content-id"

	// Act
	err := noteUsecase.AddContent(noteID, contentID, 99) // Incorrect version

	// Assert
	if err == nil {
		t.Fatal("Expected a conflict error, but got nil")
	}
	if !errors.Is(err, ErrConflict) {
		t.Errorf("Expected error to be '%v', but got '%v'", ErrConflict, err)
	}
}
