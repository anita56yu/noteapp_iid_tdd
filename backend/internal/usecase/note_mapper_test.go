package usecase

import (
	"noteapp/internal/domain"
	"noteapp/internal/repository"
	"testing"
)

func TestToNoteDTO(t *testing.T) {
	// Arrange
	note, err := domain.NewNote("note-1", "Test Title", "owner-1")
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}
	note.AddContent("content-1", "Hello", domain.TextContentType)
	note.AddContent("content-2", "base64-encoded-image", domain.ImageContentType)
	keyword1, _ := domain.NewKeyword("keyword1")
	keyword2, _ := domain.NewKeyword("keyword2")
	note.AddKeyword("user-1", keyword1)
	note.AddKeyword("user-1", keyword2)

	// Act
	mapper := NewNoteMapper()
	noteDTO := mapper.toNoteDTO(note)

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

	// Check keywords
	if len(noteDTO.Keywords) != 1 {
		t.Fatalf("Expected 1 user with keywords, got %d", len(noteDTO.Keywords))
	}
	if len(noteDTO.Keywords["user-1"]) != 2 {
		t.Fatalf("Expected 2 keywords for user-1, got %d", len(noteDTO.Keywords["user-1"]))
	}
	if noteDTO.Keywords["user-1"][0] != "keyword1" {
		t.Errorf("Expected keyword 1 to be 'keyword1', got '%s'", noteDTO.Keywords["user-1"][0])
	}
	if noteDTO.Keywords["user-1"][1] != "keyword2" {
		t.Errorf("Expected keyword 2 to be 'keyword2', got '%s'", noteDTO.Keywords["user-1"][1])
	}
}

func TestNoteMapper_ToPO(t *testing.T) {
	// Arrange
	note, _ := domain.NewNote("note-1", "Test Note", "owner-1")
	note.AddContent("content-1", "Hello", domain.TextContentType)
	keyword, _ := domain.NewKeyword("test-keyword")
	note.AddKeyword("user-1", keyword)

	mapper := NewNoteMapper()

	// Act
	po := mapper.ToPO(note)

	// Assert
	if po.ID != "note-1" {
		t.Errorf("Expected ID to be 'note-1', but got '%s'", po.ID)
	}
	if po.Title != "Test Note" {
		t.Errorf("Expected Title to be 'Test Note', but got '%s'", po.Title)
	}
	if len(po.Contents) != 1 {
		t.Fatalf("Expected 1 content block, but got %d", len(po.Contents))
	}
	if po.Contents[0].ID != "content-1" {
		t.Errorf("Expected content ID to be 'content-1', but got '%s'", po.Contents[0].ID)
	}
	if len(po.Keywords["user-1"]) != 1 {
		t.Fatalf("Expected 1 keyword for user-1, but got %d", len(po.Keywords["user-1"]))
	}
	if po.Keywords["user-1"][0] != "test-keyword" {
		t.Errorf("Expected keyword to be 'test-keyword', but got '%s'", po.Keywords["user-1"][0])
	}
}

func TestNoteMapper_ToDomain(t *testing.T) {
	// Arrange
	po := &repository.NotePO{
		ID:      "note-1",
		Title:   "Test Note",
		OwnerID: "owner-1",
		Contents: []repository.ContentPO{
			{ID: "content-1", Type: "text", Data: "Hello"},
		},
		Keywords: map[string][]string{
			"user-1": {"test-keyword"},
		},
	}

	expectedNote, _ := domain.NewNote("note-1", "Test Note", "owner-1")
	expectedNote.AddContent("content-1", "Hello", domain.TextContentType)
	keyword, _ := domain.NewKeyword("test-keyword")
	expectedNote.AddKeyword("user-1", keyword)

	mapper := NewNoteMapper()

	// Act
	note := mapper.ToDomain(po)

	// Assert
	if note.ID != expectedNote.ID {
		t.Errorf("Expected note ID to be '%s', got '%s'", expectedNote.ID, note.ID)
	}
	if note.Title != expectedNote.Title {
		t.Errorf("Expected note title to be '%s', got '%s'", expectedNote.Title, note.Title)
	}

	contents := note.Contents()
	expectedContents := expectedNote.Contents()
	if len(contents) != len(expectedContents) {
		t.Fatalf("Expected %d content blocks, got %d", len(expectedContents), len(contents))
	}

	if contents[0].ID != expectedContents[0].ID {
		t.Errorf("Expected content ID to be '%s', got '%s'", expectedContents[0].ID, contents[0].ID)
	}
	if contents[0].Type != expectedContents[0].Type {
		t.Errorf("Expected content type to be '%s', got '%s'", expectedContents[0].Type, contents[0].Type)
	}
	if contents[0].Data != expectedContents[0].Data {
		t.Errorf("Expected content data to be '%s', got '%s'", expectedContents[0].Data, contents[0].Data)
	}

	keywords := note.UserKeywords("user-1")
	if len(keywords) != 1 {
		t.Fatalf("Expected 1 keyword, but got %d", len(keywords))
	}
	if keywords[0].String() != "test-keyword" {
		t.Errorf("Expected keyword to be 'test-keyword', but got '%s'", keywords[0].String())
	}
}

func TestNoteMapper_ToPO_WithCollaborators(t *testing.T) {
	// Arrange
	note, _ := domain.NewNote("note-1", "Test Note", "owner-1")
	note.AddCollaborator("owner-1", "user-1", domain.ReadWrite)
	mapper := NewNoteMapper()

	// Act
	po := mapper.ToPO(note)

	// Assert
	if po.OwnerID != "owner-1" {
		t.Errorf("Expected OwnerID to be 'owner-1', but got '%s'", po.OwnerID)
	}
	if len(po.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(po.Collaborators))
	}
	if permission, ok := po.Collaborators["user-1"]; !ok || permission != string(domain.ReadWrite) {
		t.Errorf("Expected collaborator 'user-1' to have permission '%s', but got '%s'", domain.ReadWrite, permission)
	}
}

func TestNoteMapper_ToDomain_WithCollaborators(t *testing.T) {
	// Arrange
	po := &repository.NotePO{
		ID:      "note-1",
		OwnerID: "owner-1",
		Title:   "Test Note",
		Collaborators: map[string]string{
			"user-1": string(domain.ReadWrite),
		},
	}
	mapper := NewNoteMapper()

	// Act
	note := mapper.ToDomain(po)

	// Assert
	if note.OwnerID != "owner-1" {
		t.Errorf("Expected OwnerID to be 'owner-1', but got '%s'", note.OwnerID)
	}
	if len(note.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(note.Collaborators))
	}
	if permission, ok := note.Collaborators["user-1"]; !ok || permission != domain.ReadWrite {
		t.Errorf("Expected collaborator 'user-1' to have permission '%s', but got '%s'", domain.ReadWrite, permission)
	}
}

func TestNoteMapper_ToDTO_WithCollaborators(t *testing.T) {
	// Arrange
	note, _ := domain.NewNote("note-1", "Test Note", "owner-1")
	note.AddCollaborator("owner-1", "user-1", domain.ReadWrite)
	mapper := NewNoteMapper()

	// Act
	dto := mapper.toNoteDTO(note)

	// Assert
	if dto.OwnerID != "owner-1" {
		t.Errorf("Expected OwnerID to be 'owner-1', but got '%s'", dto.OwnerID)
	}
	if len(dto.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(dto.Collaborators))
	}
	if permission, ok := dto.Collaborators["user-1"]; !ok || permission != PermissionReadWrite {
		t.Errorf("Expected collaborator 'user-1' to have permission '%s', but got '%s'", PermissionReadWrite, permission)
	}
}
