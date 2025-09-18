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

	return &repository.NotePO{
		ID:       note.ID,
		Title:    note.Title,
		Contents: contentPOs,
	}
}

// ToDomain converts a repository.NotePO to a domain.Note.
func (m *NoteMapper) ToDomain(po *repository.NotePO) *domain.Note {
	note, _ := domain.NewNote(po.ID, po.Title)
	for _, contentPO := range po.Contents {
		note.AddContent(contentPO.ID, contentPO.Data, domain.ContentType(contentPO.Type))
	}
	return note
}

func (m *NoteMapper) toNoteDTO(note *domain.Note) *NoteDTO {
	contentDTOs := make([]ContentDTO, len(note.Contents()))
	for i, content := range note.Contents() {
		contentDTOs[i] = ContentDTO{
			ID:   content.ID,
			Type: string(content.Type),
			Data: content.Data,
		}
	}

	return &NoteDTO{
		ID:       note.ID,
		Title:    note.Title,
		Contents: contentDTOs,
	}
}
