package contentuc_test

import (
	"noteapp/internal/repository/contentrepo"
	"noteapp/internal/usecase/contentuc"
	"testing"
)

func TestContentUsecase_CreateContent_WithInjectedID(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	contentID := "c1"
	data := "Test content"
	contentType := contentuc.TextContentType

	returnedID, err := usecase.CreateContent(noteID, contentID, data, contentType)
	if err != nil {
		t.Fatalf("CreateContent() returned an unexpected error: %v", err)
	}

	if returnedID != contentID {
		t.Errorf("Expected ID to be '%s', got '%s'", contentID, returnedID)
	}

	po, err := repo.GetByID(contentID)
	if err != nil {
		t.Fatalf("GetByID() returned an unexpected error: %v", err)
	}
	if po.NoteID != noteID {
		t.Errorf("Expected NoteID to be '%s', got '%s'", noteID, po.NoteID)
	}
	if po.Data != data {
		t.Errorf("Expected Data to be '%s', got '%s'", data, po.Data)
	}
	if po.Type != string(contentType) {
		t.Errorf("Expected Type to be '%s', got '%s'", contentType, po.Type)
	}
	if po.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", po.Version)
	}
}

func TestContentUsecase_CreateContent_WithGeneratedID(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := contentuc.TextContentType

	returnedID, err := usecase.CreateContent(noteID, "", data, contentType)
	if err != nil {
		t.Fatalf("CreateContent() returned an unexpected error: %v", err)
	}

	if returnedID == "" {
		t.Error("Expected a generated ID, but got an empty string")
	}

	po, err := repo.GetByID(returnedID)
	if err != nil {
		t.Fatalf("GetByID() returned an unexpected error: %v", err)
	}
	if po.NoteID != noteID {
		t.Errorf("Expected NoteID to be '%s', got '%s'", noteID, po.NoteID)
	}
	if po.Data != data {
		t.Errorf("Expected Data to be '%s', got '%s'", data, po.Data)
	}
	if po.Type != string(contentType) {
		t.Errorf("Expected Type to be '%s', got '%s'", contentType, po.Type)
	}
	if po.Version != 0 {
		t.Errorf("Expected Version to be 0, but got %d", po.Version)
	}
}

func TestContentUsecase_CreateContent_UnsupportedType(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := "unsupported"

	_, err := usecase.CreateContent(noteID, "", data, contentuc.ContentType(contentType))
	if err != contentuc.ErrUnsupportedContentType {
		t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrUnsupportedContentType, err)
	}
}

func TestContentUsecase_GetContentByID(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := contentuc.TextContentType

	id, _ := usecase.CreateContent(noteID, "", data, contentType)

	t.Run("should get content by id", func(t *testing.T) {
		dto, err := usecase.GetContentByID(id)
		if err != nil {
			t.Fatalf("GetContentByID() returned an unexpected error: %v", err)
		}
		if dto.ID != id {
			t.Errorf("Expected ID to be '%s', got '%s'", id, dto.ID)
		}
		if dto.Data != data {
			t.Errorf("Expected Data to be '%s', got '%s'", data, dto.Data)
		}
	})

	t.Run("should return not found error for non-existent id", func(t *testing.T) {
		_, err := usecase.GetContentByID("non-existent-id")
		if err != contentuc.ErrContentNotFound {
			t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrContentNotFound, err)
		}
	})

	t.Run("should return invalid id error for empty id", func(t *testing.T) {
		_, err := usecase.GetContentByID("")
		if err != contentuc.ErrInvalidID {
			t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrInvalidID, err)
		}
	})
}

func TestContentUsecase_UpdateContent(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := contentuc.TextContentType

	id, _ := usecase.CreateContent(noteID, "", data, contentType)

	updatedData := "Updated data"
	err := usecase.UpdateContent(id, updatedData, 0)
	if err != nil {
		t.Fatalf("UpdateContent() returned an unexpected error: %v", err)
	}

	po, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID() returned an unexpected error: %v", err)
	}
	if po.Data != updatedData {
		t.Errorf("Expected Data to be '%s', got '%s'", updatedData, po.Data)
	}
	if po.Version != 1 {
		t.Errorf("Expected Version to be 1, but got %d", po.Version)
	}
}

func TestContentUsecase_UpdateContent_NotFound(t *testing.T) {
	usecase := contentuc.NewContentUsecase(contentrepo.NewInMemoryContentRepository())

	err := usecase.UpdateContent("non-existent-id", "updated data", 0)
	if err != contentuc.ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrContentNotFound, err)
	}
}

func TestContentUsecase_UpdateContent_Conflict(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := contentuc.TextContentType

	id, _ := usecase.CreateContent(noteID, "", data, contentType)

	// Attempt to update with an incorrect version
	err := usecase.UpdateContent(id, "updated data", 99)
	if err != contentuc.ErrConflict {
		t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrConflict, err)
	}
}

func TestContentUsecase_DeleteContent(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := contentuc.TextContentType

	id, _ := usecase.CreateContent(noteID, "", data, contentType)

	err := usecase.DeleteContent(id, 0)
	if err != nil {
		t.Fatalf("DeleteContent() returned an unexpected error: %v", err)
	}

	_, err = repo.GetByID(id)
	if err != contentrepo.ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", contentrepo.ErrContentNotFound, err)
	}
}

func TestContentUsecase_DeleteContent_NotFound(t *testing.T) {
	usecase := contentuc.NewContentUsecase(contentrepo.NewInMemoryContentRepository())

	err := usecase.DeleteContent("non-existent-id", 0)
	if err != contentuc.ErrContentNotFound {
		t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrContentNotFound, err)
	}
}

func TestContentUsecase_DeleteContent_Conflict(t *testing.T) {
	repo := contentrepo.NewInMemoryContentRepository()
	usecase := contentuc.NewContentUsecase(repo)
	noteID := "n1"
	data := "Test content"
	contentType := contentuc.TextContentType

	id, _ := usecase.CreateContent(noteID, "", data, contentType)

	err := usecase.DeleteContent(id, 99)
	if err != contentuc.ErrConflict {
		t.Errorf("Expected error to be '%v', but got '%v'", contentuc.ErrConflict, err)
	}
}