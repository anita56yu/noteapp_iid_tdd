package noteuc

import (
	domainnote "noteapp/internal/domain/note"
	"noteapp/internal/repository/noterepo"
)

// NoteMapper handles mapping between domain.Note and other representations.
type NoteMapper struct{}

// NewNoteMapper creates a new NoteMapper.
func NewNoteMapper() *NoteMapper {
	return &NoteMapper{}
}

// ToPO converts a domain.Note to a repository.NotePO.
func (m *NoteMapper) ToPO(note *domainnote.Note) *noterepo.NotePO {
	contentPOs := make([]noterepo.ContentPO, len(note.Contents()))
	for i, content := range note.Contents() {
		contentPOs[i] = noterepo.ContentPO{
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

	return &noterepo.NotePO{
		ID:            note.ID,
		OwnerID:       note.OwnerID,
		Title:         note.Title,
		Version:       note.Version,
		Contents:      contentPOs,
		Keywords:      keywordPOs,
		Collaborators: collaboratorPOs,
	}
}

// ToDomain converts a repository.NotePO to a domain.Note.
func (m *NoteMapper) ToDomain(po *noterepo.NotePO) *domainnote.Note {
	note, _ := domainnote.NewNoteWithVersion(po.ID, po.Title, po.OwnerID, po.Version)
	for _, contentPO := range po.Contents {
		note.AddContent(contentPO.ID, contentPO.Data, domainnote.ContentType(contentPO.Type))
	}
	for userID, keywords := range po.Keywords {
		for _, keywordStr := range keywords {
			keyword, _ := domainnote.NewKeyword(keywordStr)
			note.AddKeyword(userID, keyword)
		}
	}
	for userID, permission := range po.Collaborators {
		note.AddCollaborator(note.OwnerID, userID, domainnote.Permission(permission))
	}
	return note
}

func (m *NoteMapper) toNoteDTO(note *domainnote.Note) *NoteDTO {
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
