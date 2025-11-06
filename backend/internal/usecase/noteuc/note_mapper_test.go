package noteuc

import (
	"noteapp/internal/domain/note"
	"noteapp/internal/repository/noterepo"
	"testing"
)

func TestToNoteDTO(t *testing.T) {
	// Arrange
	n, err := note.NewNoteWithVersion("note-1", "Test Title", "owner-1", 1)
	if err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}
	n.AddContentID("content-1", -1)
	n.AddContentID("content-2", -1)
	keyword1, _ := note.NewKeyword("keyword1")
	keyword2, _ := note.NewKeyword("keyword2")
	n.AddKeyword("user-1", keyword1)
	n.AddKeyword("user-1", keyword2)

	// Act
	mapper := NewNoteMapper()
	noteDTO := mapper.toNoteDTO(n)

	// Assert
	if noteDTO.ID != "note-1" {
		t.Errorf("Expected DTO ID to be 'note-1', got '%s'", noteDTO.ID)
	}
	if noteDTO.Title != "Test Title" {
		t.Errorf("Expected DTO title to be 'Test Title', got '%s'", noteDTO.Title)
	}
	if noteDTO.Version != 1 {
		t.Errorf("Expected DTO version to be 1, got %d", noteDTO.Version)
	}
	if len(noteDTO.ContentIDs) != 2 {
		t.Fatalf("Expected 2 content IDs, got %d", len(noteDTO.ContentIDs))
	}
	if noteDTO.ContentIDs[0] != "content-1" || noteDTO.ContentIDs[1] != "content-2" {
		t.Errorf("Expected ContentIDs to be [content-1, content-2], got %v", noteDTO.ContentIDs)
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
	n, _ := note.NewNote("note-1", "Test Note", "owner-1")
	n.AddContentID("content-1", -1)
	keyword, _ := note.NewKeyword("test-keyword")
	n.AddKeyword("user-1", keyword)

	mapper := NewNoteMapper()

	// Act
	po := mapper.ToPO(n)

	// Assert
	if po.ID != "note-1" {
		t.Errorf("Expected ID to be 'note-1', but got '%s'", po.ID)
	}
	if po.Title != "Test Note" {
		t.Errorf("Expected Title to be 'Test Note', but got '%s'", po.Title)
	}
	if len(po.ContentIDs) != 1 {
		t.Fatalf("Expected 1 content ID, but got %d", len(po.ContentIDs))
	}
	if po.ContentIDs[0] != "content-1" {
		t.Errorf("Expected content ID to be 'content-1', but got '%s'", po.ContentIDs[0])
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
	po := &noterepo.NotePO{
		ID:         "note-1",
		Title:      "Test Note",
		OwnerID:    "owner-1",
		ContentIDs: []string{"content-1"},
		Keywords: map[string][]string{
			"user-1": {"test-keyword"},
		},
	}

	expectedNote, _ := note.NewNote("note-1", "Test Note", "owner-1")
	expectedNote.AddContentID("content-1", -1)
	keyword, _ := note.NewKeyword("test-keyword")
	expectedNote.AddKeyword("user-1", keyword)

	mapper := NewNoteMapper()

	// Act
	n := mapper.ToDomain(po)

	// Assert
	if n.ID != expectedNote.ID {
		t.Errorf("Expected note ID to be '%s', got '%s'", expectedNote.ID, n.ID)
	}
	if n.Title != expectedNote.Title {
		t.Errorf("Expected note title to be '%s', got '%s'", expectedNote.Title, n.Title)
	}

	if len(n.ContentIDs) != 1 {
		t.Fatalf("Expected 1 content ID, but got %d", len(n.ContentIDs))
	}
	if n.ContentIDs[0] != expectedNote.ContentIDs[0] {
		t.Errorf("Expected content ID to be 'content-1', but got '%s'", n.ContentIDs[0])
	}

	keywords := n.UserKeywords("user-1")
	if len(keywords) != 1 {
		t.Fatalf("Expected 1 keyword, but got %d", len(keywords))
	}
	if keywords[0].String() != "test-keyword" {
		t.Errorf("Expected keyword to be 'test-keyword', but got '%s'", keywords[0].String())
	}
}

func TestNoteMapper_ToPO_WithCollaborators(t *testing.T) {
	// Arrange
	n, _ := note.NewNote("note-1", "Test Note", "owner-1")
	n.AddCollaborator("owner-1", "user-1", note.ReadWrite)
	mapper := NewNoteMapper()

	// Act
	po := mapper.ToPO(n)

	// Assert
	if po.OwnerID != "owner-1" {
		t.Errorf("Expected OwnerID to be 'owner-1', but got '%s'", po.OwnerID)
	}
	if len(po.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(po.Collaborators))
	}
	if permission, ok := po.Collaborators["user-1"]; !ok || permission != string(note.ReadWrite) {
		t.Errorf("Expected collaborator 'user-1' to have permission '%s', but got '%s'", note.ReadWrite, permission)
	}
}

func TestNoteMapper_ToDomain_WithCollaborators(t *testing.T) {
	// Arrange
	po := &noterepo.NotePO{
		ID:      "note-1",
		OwnerID: "owner-1",
		Title:   "Test Note",
		Collaborators: map[string]string{
			"user-1": string(note.ReadWrite),
		},
	}
	mapper := NewNoteMapper()

	// Act
	n := mapper.ToDomain(po)

	// Assert
	if n.OwnerID != "owner-1" {
		t.Errorf("Expected OwnerID to be 'owner-1', but got '%s'", n.OwnerID)
	}
	if len(n.Collaborators) != 1 {
		t.Fatalf("Expected 1 collaborator, but got %d", len(n.Collaborators))
	}
	if permission, ok := n.Collaborators["user-1"]; !ok || permission != note.ReadWrite {
		t.Errorf("Expected collaborator 'user-1' to have permission '%s', but got '%s'", note.ReadWrite, permission)
	}
}

func TestNoteMapper_ToDTO_WithCollaborators(t *testing.T) {
	// Arrange
	n, _ := note.NewNote("note-1", "Test Note", "owner-1")
	n.AddCollaborator("owner-1", "user-1", note.ReadWrite)
	mapper := NewNoteMapper()

	// Act
	dto := mapper.toNoteDTO(n)

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
