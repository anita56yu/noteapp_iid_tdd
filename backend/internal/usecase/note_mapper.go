package usecase

import "noteapp/internal/domain"

// toNoteDTO converts a domain.Note object to a NoteDTO object.
func toNoteDTO(note *domain.Note) *NoteDTO {
	if note == nil {
		return nil
	}

	contentDTOs := make([]ContentDTO, 0, len(note.Contents()))
	for _, content := range note.Contents() {
		contentDTOs = append(contentDTOs, ContentDTO{
			ID:   content.ID,
			Type: string(content.Type),
			Data: content.Data,
		})
	}

	return &NoteDTO{
		ID:       note.ID,
		Title:    note.Title,
		Contents: contentDTOs,
	}
}
