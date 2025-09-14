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

type ContentType string

const (
	// TextContentType represents a text content block.
	TextContentType ContentType = "text"
	// ImageContentType represents an image content block.
	ImageContentType ContentType = "image"
)

// NoteUsecase handles the business logic for notes.
type NoteUsecase struct {
	repo repository.NoteRepository
}

// NewNoteUsecase creates a new NoteUsecase.
func NewNoteUsecase(repo repository.NoteRepository) *NoteUsecase {
	return &NoteUsecase{repo: repo}
}

// CreateNote creates a new note.
func (uc *NoteUsecase) CreateNote(id, title string) (string, error) {
	note, err := domain.NewNote(id, title)
	if err != nil {
		return "", uc.mapDomainError(err)
	}

	if err := uc.repo.Save(note); err != nil {
		return "", uc.mapRepositoryError(err)
	}

	return note.ID, nil
}

// GetNoteByID retrieves a note by its ID.
func (uc *NoteUsecase) GetNoteByID(id string) (*NoteDTO, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	note, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, uc.mapRepositoryError(err)
	}

	return toNoteDTO(note), nil
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
	note, err := uc.repo.FindByID(noteID)
	if err != nil {
		return "", uc.mapRepositoryError(err)
	}

	newID := note.AddContent(contentID, data, mapToDomainContentType(contentType))
	if err := uc.repo.Save(note); err != nil {
		return "", uc.mapRepositoryError(err)
	}

	return newID, nil
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
