package repository

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
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

func TestInMemoryNoteRepository_GetAccessibleNoteByUserID(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	ownerID := "user-1"
	otherUserID := "user-2"

	// Note owned by the user
	ownedNote := &NotePO{ID: "owned-note", OwnerID: ownerID}
	repo.Save(ownedNote)

	// Note shared with the user
	sharedNote := &NotePO{ID: "shared-note", OwnerID: otherUserID, Collaborators: map[string]string{ownerID: "read"}}
	repo.Save(sharedNote)

	// Note not related to the user
	otherNote := &NotePO{ID: "other-note", OwnerID: otherUserID}
	repo.Save(otherNote)

	// Act
	notes, err := repo.GetAccessibleNoteByUserID(ownerID)

	// Assert
	if err != nil {
		t.Fatalf("GetAccessibleNoteByUserID() returned an unexpected error: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("Expected 2 notes, but got %d", len(notes))
	}

	foundOwned := false
	foundShared := false
	for _, note := range notes {
		if note.ID == "owned-note" {
			foundOwned = true
		}
		if note.ID == "shared-note" {
			foundShared = true
		}
	}

	if !foundOwned {
		t.Error("Expected to find the owned note, but it was not returned")
	}
	if !foundShared {
		t.Error("Expected to find the shared note, but it was not returned")
	}
}

func TestInMemoryNoteRepository_ConcurrentContentAddOnSameNote(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	noteID := "concurrent-note"
	initialNote := &NotePO{
		ID:       noteID,
		Title:    "Initial Title",
		Contents: []ContentPO{},
	}
	repo.Save(initialNote)

	var wg sync.WaitGroup
	contentUser1 := ContentPO{ID: "content-1", Type: "text", Data: "Content from user 1"}
	contentUser2 := ContentPO{ID: "content-2", Type: "text", Data: "Content from user 2"}

	// Act
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := repo.LockNoteForUpdate(noteID)
		if err != nil {
			t.Errorf("Failed to lock note for update: %v", err)
			return
		}
		defer repo.UnlockNoteForUpdate(noteID)
		note, _ := repo.FindByID(noteID)
		note.Contents = append(note.Contents, contentUser1)
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
		repo.Save(note)
	}()
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond) // Ensure this runs slightly after the first goroutine starts
		err := repo.LockNoteForUpdate(noteID)
		if err != nil {
			t.Errorf("Failed to lock note for update: %v", err)
			return
		}
		defer repo.UnlockNoteForUpdate(noteID)
		note, _ := repo.FindByID(noteID)
		note.Contents = append(note.Contents, contentUser2)
		repo.Save(note)
	}()
	wg.Wait()

	// Assert
	finalNote, _ := repo.FindByID(noteID)
	if len(finalNote.Contents) != 2 {
		t.Errorf("Expected 2 content items, got %d. A content update was lost.", len(finalNote.Contents))
	}
	if (finalNote.Contents[0].ID != "content-1" && finalNote.Contents[0].ID != "content-2") || (finalNote.Contents[1].ID != "content-1" && finalNote.Contents[1].ID != "content-2") {
		t.Errorf("Content IDs do not match expected values: %v", finalNote.Contents)
	}
}

// TestInMemoryNoteRepository_ConcurrentMapReadWrite simulates simultaneous reading and writing to the notes map.
// This test will fail with a race condition if the map access is not protected by a mutex.
func TestInMemoryNoteRepository_ConcurrentMapReadWrite(t *testing.T) {
	// Arrange
	repo := NewInMemoryNoteRepository()
	var wg sync.WaitGroup
	const numOperations = 100

	// Pre-populate one note to read
	const readNoteID = "read-note-id"
	repo.Save(&NotePO{ID: readNoteID, Title: "A note to be read"})

	// Act
	wg.Add(2)

	// Goroutine 1: Continuously writes new notes to the map
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			noteID := fmt.Sprintf("note-%d", i)
			repo.Save(&NotePO{ID: noteID, Title: "New Note"})
		}
	}()

	// Goroutine 2: Continuously reads an existing note from the map
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			_, err := repo.FindByID(readNoteID)
			if err != nil && !errors.Is(err, ErrNoteNotFound) {
				t.Errorf("FindByID returned an unexpected error: %v", err)
			}
		}
	}()

	wg.Wait()

	// Assert
	// The primary assertion is that the test completes without a race condition panic.
	// We can do a final check to ensure the writes happened.
	repo.mu.RLock()
	finalCount := len(repo.notes)
	repo.mu.RUnlock()

	if finalCount != numOperations+1 { // +1 for the initial read-note
		t.Errorf("Expected %d notes in the repository, but got %d", numOperations+1, finalCount)
	}
}
