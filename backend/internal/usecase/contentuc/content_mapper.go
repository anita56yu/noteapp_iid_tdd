package contentuc

import (
	"noteapp/internal/domain/content"
	"noteapp/internal/repository/contentrepo"
)

// ContentMapper handles mapping between domain.Content and other representations.
type ContentMapper struct{}

// NewContentMapper creates a new ContentMapper.
func NewContentMapper() *ContentMapper {
	return &ContentMapper{}
}

// ToPO converts a domain.Content to a repository.ContentPO.
func (m *ContentMapper) ToPO(c *content.Content) *contentrepo.ContentPO {
	return &contentrepo.ContentPO{
		ID:      c.ID,
		NoteID:  c.NoteID,
		Data:    c.Data,
		Type:    string(c.Type),
		Version: c.Version,
	}
}

// ToDomain converts a repository.ContentPO to a domain.Content.
func (m *ContentMapper) ToDomain(po *contentrepo.ContentPO) *content.Content {
	return content.NewContent(po.ID, po.NoteID, po.Data, content.ContentType(po.Type), po.Version)
}

// ToDTO converts a domain.Content to a ContentDTO.
func (m *ContentMapper) ToDTO(c *content.Content) *ContentDTO {
	return &ContentDTO{
		ID:      c.ID,
		NoteID:  c.NoteID,
		Data:    c.Data,
		Type:    string(c.Type),
		Version: c.Version,
	}
}
