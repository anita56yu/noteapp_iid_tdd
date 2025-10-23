package contentrepo_test

import (
	"noteapp/internal/repository/contentrepo"
	"sync"
	"testing"
	"time"
)

func TestInMemoryContentRepository_Save_Conflict(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	content := &contentrepo.ContentPO{
		ID:      "c1",
		NoteID:  "n1",
		Data:    "Test data",
		Type:    "text/plain",
		Version: 0,
	}
	repo.Save(content)

	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error

	// Goroutine 1
	go func() {
		defer wg.Done()
		savedContent, _ := repo.GetByID("c1")
		savedContent.Data = "Goroutine 1"
		time.Sleep(10 * time.Millisecond)
		err1 = repo.Save(savedContent)
	}()

	// Goroutine 2
	go func() {
		defer wg.Done()
		savedContent, _ := repo.GetByID("c1")
		savedContent.Data = "Goroutine 2"
		time.Sleep(20 * time.Millisecond)
		err2 = repo.Save(savedContent)
	}()

	wg.Wait()

	if (err1 == nil && err2 != contentrepo.ErrContentConflict) || (err2 == nil && err1 != contentrepo.ErrContentConflict) {
		t.Errorf("Expected one of the saves to fail with a conflict error, but got err1: %v, err2: %v", err1, err2)
	}
}

func TestInMemoryContentRepository_Save(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	content := &contentrepo.ContentPO{
		ID:     "c1",
		NoteID: "n1",
		Data:   "Test data",
		Type:   "text/plain",
	}

	err := repo.Save(content)
	if err != nil {
		t.Fatalf("Save returned an unexpected error: %v", err)
	}

	saved, err := repo.GetByID("c1")
	if err != nil {
		t.Fatalf("GetByID returned an unexpected error: %v", err)
	}

	if saved.ID != content.ID {
		t.Errorf("Expected ID to be '%s', but got '%s'", content.ID, saved.ID)
	}
	if saved.NoteID != content.NoteID {
		t.Errorf("Expected NoteID to be '%s', but got '%s'", content.NoteID, saved.NoteID)
	}
	if saved.Data != content.Data {
		t.Errorf("Expected Data to be '%s', but got '%s'", content.Data, saved.Data)
	}
	if saved.Type != content.Type {
		t.Errorf("Expected Type to be '%s', but got '%s'", content.Type, saved.Type)
	}
	if saved.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", saved.Version)
	}
}

func TestInMemoryContentRepository_GetAllByNoteID(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	c1 := &contentrepo.ContentPO{ID: "c1", NoteID: "n1"}
	c2 := &contentrepo.ContentPO{ID: "c2", NoteID: "n1"}
	c3 := &contentrepo.ContentPO{ID: "c3", NoteID: "n2"}
	repo.Save(c1)
	repo.Save(c2)
	repo.Save(c3)

	contents, err := repo.GetAllByNoteID("n1")
	if err != nil {
		t.Fatalf("GetAllByNoteID returned an unexpected error: %v", err)
	}
	if len(contents) != 2 {
		t.Errorf("Expected 2 contents, but got %d", len(contents))
	}
}

func TestInMemoryContentRepository_GetByID(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	content := &contentrepo.ContentPO{
		ID:     "c1",
		NoteID: "n1",
		Data:   "Test data",
		Type:   "text/plain",
	}
	repo.Save(content)

	saved, err := repo.GetByID("c1")
	if err != nil {
		t.Fatalf("GetByID returned an unexpected error: %v", err)
	}
	if saved.ID != content.ID {
		t.Errorf("Expected ID to be '%s', but got '%s'", content.ID, saved.ID)
	}
}

func TestInMemoryContentRepository_GetByID_NotFound(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	_, err := repo.GetByID("non-existent-id")
	if err != contentrepo.ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", contentrepo.ErrContentNotFound, err)
	}
}

func TestInMemoryContentRepository_Delete(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	c1 := &contentrepo.ContentPO{ID: "c1", NoteID: "n1"}
	repo.Save(c1)

	err := repo.Delete("c1")
	if err != nil {
		t.Fatalf("Delete returned an unexpected error: %v", err)
	}

	_, err = repo.GetByID("c1")
	if err != contentrepo.ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", contentrepo.ErrContentNotFound, err)
	}

	err = repo.Delete("c1")
	if err != contentrepo.ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", contentrepo.ErrContentNotFound, err)
	}
}
