package usecase

import (
	"noteapp/internal/domain"
	"noteapp/internal/repository"
)

// NoteMapper handles mapping between domain.Note and other representations.
type NoteMapper struct{}

// NewNoteMapper creates a new NoteMapper.
func NewNoteMapper() *NoteMapper {
	return &NoteMapper{}
}

// ToPO converts a domain.Note to a repository.NotePO.
func (m *NoteMapper) ToPO(note *domain.Note) *repository.NotePO {
	contentPOs := make([]repository.ContentPO, len(note.Contents()))
	for i, content := range note.Contents() {
		contentPOs[i] = repository.ContentPO{
			ID:   content.ID,
			Type: string(content.Type),
			Data: content.Data,
		}
	}

	keywordPOs := make(map[string][]string)
	for userID, keywords := range note.Keywords() {
		for _, keyword := range keywords {
			keywordPOs[userID] = append(keywordPOs[userID], keyword.String())
		}
	}

	collaboratorPOs := make(map[string]string)
	for userID, permission := range note.Collaborators {
		collaboratorPOs[userID] = string(permission)
	}

	return &repository.NotePO{
		ID:            note.ID,
		OwnerID:       note.OwnerID,
		Title:         note.Title,
		Contents:      contentPOs,
		Keywords:      keywordPOs,
		Collaborators: collaboratorPOs,
	}
}

// ToDomain converts a repository.NotePO to a domain.Note.
func (m *NoteMapper) ToDomain(po *repository.NotePO) *domain.Note {
	note, _ := domain.NewNote(po.ID, po.Title, po.OwnerID)
	for _, contentPO := range po.Contents {
		note.AddContent(contentPO.ID, contentPO.Data, domain.ContentType(contentPO.Type))
	}
	for userID, keywords := range po.Keywords {
		for _, keywordStr := range keywords {
			keyword, _ := domain.NewKeyword(keywordStr)
			note.AddKeyword(userID, keyword)
		}
	}
	for userID, permission := range po.Collaborators {
		note.AddCollaborator(note.OwnerID, userID, domain.Permission(permission))
	}
	return note
}

func (m *NoteMapper) toNoteDTO(note *domain.Note) *NoteDTO {
	contents := []*ContentDTO{}
	for _, c := range note.Contents() {
		contents = append(contents, &ContentDTO{
			ID:   c.ID,
			Type: string(c.Type),
			Data: c.Data,
		})
	}

	keywords := make(map[string][]string)
	for userID, userKeywords := range note.Keywords() {
		for _, keyword := range userKeywords {
			keywords[userID] = append(keywords[userID], keyword.String())
		}
	}

	collaborators := make(map[string]Permission)
	for userID, permission := range note.Collaborators {
		collaborators[userID] = Permission(permission)
	}

	return &NoteDTO{
		ID:            note.ID,
		Title:         note.Title,
		OwnerID:       note.OwnerID,
		Contents:      contents,
		Keywords:      keywords,
		Collaborators: collaborators,
	}
}
