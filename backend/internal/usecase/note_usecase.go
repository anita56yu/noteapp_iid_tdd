package usecase

import (
	"errors"
	"fmt"
	"noteapp/internal/domain"
	"noteapp/internal/repository"
)

// ErrInvalidID is returned when an invalid ID is provided.
var ErrInvalidID = errors.New("invalid ID")

// ErrNoteNotFound is returned when a note is not found.
var ErrNoteNotFound = errors.New("note not found")

// ErrNilNote is returned when a nil note is passed to a method.
var ErrNilNote = errors.New("note cannot be nil")

// ErrEmptyTitle is returned when a note is created with an empty title.
var ErrEmptyTitle = errors.New("title cannot be empty")

// ErrContentNotFound is returned when a content is not found.
var ErrContentNotFound = errors.New("content not found")

type ContentType string

const (
	// TextContentType represents a text content block.
	TextContentType ContentType = "text"
	// ImageContentType represents an image content block.
	ImageContentType ContentType = "image"
)

// NoteUsecase handles the business logic for notes.
type NoteUsecase struct {
	repo   repository.NoteRepository
	mapper *NoteMapper
}

// NewNoteUsecase creates a new NoteUsecase.
func NewNoteUsecase(repo repository.NoteRepository) *NoteUsecase {
	return &NoteUsecase{repo: repo, mapper: NewNoteMapper()}
}

// CreateNote creates a new note.
func (uc *NoteUsecase) CreateNote(id, title string) (string, error) {
	note, err := domain.NewNote(id, title)
	if err != nil {
		return "", uc.mapDomainError(err)
	}

	notePO := uc.mapper.ToPO(note)
	if err := uc.repo.Save(notePO); err != nil {
		return "", uc.mapRepositoryError(err)
	}

	return note.ID, nil
}

// GetNoteByID retrieves a note by its ID.
func (uc *NoteUsecase) GetNoteByID(id string) (*NoteDTO, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	notePO, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, uc.mapRepositoryError(err)
	}

	note := uc.mapper.ToDomain(notePO)
	return uc.mapper.toNoteDTO(note), nil
}

// DeleteNote deletes a note by its ID.
func (uc *NoteUsecase) DeleteNote(id string) error {
	if id == "" {
		return ErrInvalidID
	}

	if err := uc.repo.Delete(id); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

func (uc *NoteUsecase) AddContent(noteID, contentID, data string, contentType ContentType) (string, error) {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return "", uc.mapRepositoryError(err)
	}
	note := uc.mapper.ToDomain(notePO)

	newID := note.AddContent(contentID, data, mapToDomainContentType(contentType))

	updatedNotePO := uc.mapper.ToPO(note)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return "", uc.mapRepositoryError(err)
	}

	return newID, nil
}

// UpdateContent updates the content of a note.
func (uc *NoteUsecase) UpdateContent(noteID, contentID, data string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	note := uc.mapper.ToDomain(notePO)

	if err := note.UpdateContent(contentID, data); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(note)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

func (uc *NoteUsecase) mapRepositoryError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNoteNotFound):
		return ErrNoteNotFound
	case errors.Is(err, repository.ErrNilNote):
		return ErrNilNote
	default:
		return fmt.Errorf("an unexpected repositoryerror occurred: %w", err)
	}
}

func (uc *NoteUsecase) mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmptyTitle):
		return ErrEmptyTitle
	case errors.Is(err, domain.ErrContentNotFound):
		return ErrContentNotFound
	default:
		return fmt.Errorf("an unexpected domain error occurred: %w", err)
	}
}

func mapToDomainContentType(ct ContentType) domain.ContentType {
	switch ct {
	case TextContentType:
		return domain.TextContentType
	case ImageContentType:
		return domain.ImageContentType
	default:
		return domain.TextContentType // Default to text if unknown
	}
}
