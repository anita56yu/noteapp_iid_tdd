package contentuc

import (
	"errors"
	"fmt"
	"noteapp/internal/domain/content"
	"noteapp/internal/repository/contentrepo"
)

// ContentUsecase handles the business logic for content.
type ContentUsecase struct {
	repo   contentrepo.ContentRepository
	mapper *ContentMapper
}

// NewContentUsecase creates a new ContentUsecase.
func NewContentUsecase(repo contentrepo.ContentRepository) *ContentUsecase {
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
		return "", uc.mapRepositoryError(err)
	}
	return c.ID, nil
}

// UpdateContent updates a content.
func (uc *ContentUsecase) UpdateContent(id, data string) error {
	po, err := uc.repo.GetByID(id)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	c := uc.mapper.ToDomain(po)

	// For now, we only support updating the data.
	c.Data = data

	po = uc.mapper.ToPO(c)
	if err := uc.repo.Save(po); err != nil {
		return uc.mapRepositoryError(err)
	}
	return nil
}

// DeleteContent deletes a content.
func (uc *ContentUsecase) DeleteContent(id string) error {
	if err := uc.repo.Delete(id); err != nil {
		return uc.mapRepositoryError(err)
	}
	return nil
}

func (uc *ContentUsecase) mapRepositoryError(err error) error {
	switch {
	case errors.Is(err, contentrepo.ErrContentNotFound):
		return ErrContentNotFound
	case errors.Is(err, contentrepo.ErrContentConflict):
		return ErrConflict
	default:
		return fmt.Errorf("an unexpected repository error occurred: %w", err)
	}
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
