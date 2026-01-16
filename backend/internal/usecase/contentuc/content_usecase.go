package contentuc

import (
	"noteapp/internal/domain/content"
)

// ContentUsecase handles the business logic for content.
type ContentUsecase struct {
	repo   ContentRepository
	mapper *ContentMapper
}

// NewContentUsecase creates a new ContentUsecase.
func NewContentUsecase(repo ContentRepository) *ContentUsecase {
	return &ContentUsecase{repo: repo, mapper: NewContentMapper()}
}

// CreateContent creates a new content.
func (uc *ContentUsecase) CreateContent(noteID, contentID, data string, contentType ContentType) (string, error) {
	domainContentType, err := mapToDomainContentType(contentType)
	if err != nil {
		return "", err
	}
	c := content.NewContent(contentID, noteID, data, domainContentType, 0)
	po := uc.mapper.ToPO(c)
	if err := uc.repo.Save(po); err != nil {
		return "", err
	}
	return c.ID, nil
}

// GetContentByID retrieves a content by its ID.
func (uc *ContentUsecase) GetContentByID(id string) (*ContentDTO, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	po, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	c := uc.mapper.ToDomain(po)
	return uc.mapper.ToDTO(c), nil
}

// UpdateContent updates a content.
func (uc *ContentUsecase) UpdateContent(id, data string, version int) error {
	po, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if po.Version != version {
		return ErrConflict
	}

	c := uc.mapper.ToDomain(po)

	// For now, we only support updating the data.
	c.Data = data

	po = uc.mapper.ToPO(c)
	if err := uc.repo.Save(po); err != nil {
		return err
	}
	return nil
}

// DeleteContent deletes a content.
func (uc *ContentUsecase) DeleteContent(id string, version int) error {
	po, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	if po.Version != version {
		return ErrConflict
	}

	if err := uc.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

// DeleteAllContentsByNoteID deletes all content associated with a given note ID.
func (uc *ContentUsecase) DeleteAllContentsByNoteID(noteID string) error {
	if err := uc.repo.DeleteAllByNoteID(noteID); err != nil {
		return err
	}
	return nil
}

func mapToDomainContentType(ct ContentType) (content.ContentType, error) {
	switch ct {
	case TextContentType:
		return content.TextContentType, nil
	case ImageContentType:
		return content.ImageContentType, nil
	default:
		return "", ErrUnsupportedContentType
	}
}
