package noteuc

import (
	"errors"
	"fmt"
	"noteapp/internal/domain/note"
	"noteapp/internal/repository/noterepo"
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

// ErrEmptyKeyword is returned when a keyword is empty.
var ErrEmptyKeyword = errors.New("keyword cannot be empty")

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrKeywordNotFound is returned when a keyword is not found.
var ErrKeywordNotFound = errors.New("keyword not found")

// ErrPermissionDenied is returned when a user is not authorized to perform an action.
var ErrPermissionDenied = errors.New("permission denied")

// ErrUnsupportedPermissionType is returned when an unsupported permission type is provided.
var ErrUnsupportedPermissionType = errors.New("unsupported permission type")

// ErrConflict is returned when a version conflict occurs.
var ErrConflict = errors.New("conflict")

type ContentType string

const (
	// TextContentType represents a text content block.
	TextContentType ContentType = "text"
	// ImageContentType represents an image content block.
	ImageContentType ContentType = "image"
)

// NoteUsecase handles the business logic for notes.
type NoteUsecase struct {
	repo   noterepo.NoteRepository
	mapper *NoteMapper
}

// NewNoteUsecase creates a new NoteUsecase.
func NewNoteUsecase(repo noterepo.NoteRepository) *NoteUsecase {
	return &NoteUsecase{repo: repo, mapper: NewNoteMapper()}
}

// CreateNote creates a new note.
func (uc *NoteUsecase) CreateNote(id, title, ownerID string) (string, error) {
	n, err := note.NewNote(id, title, ownerID)
	if err != nil {
		return "", uc.mapDomainError(err)
	}

	notePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(notePO); err != nil {
		return "", uc.mapRepositoryError(err)
	}

	return n.ID, nil
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

	n := uc.mapper.ToDomain(notePO)
	return uc.mapper.toNoteDTO(n), nil
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

func (uc *NoteUsecase) CreateAndAddContent(noteID, contentID, data string, contentType ContentType) (string, error) {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return "", uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	newID := n.AddContent(contentID, data, mapToDomainContentType(contentType))

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return "", uc.mapRepositoryError(err)
	}

	return newID, nil
}

func (uc *NoteUsecase) AddContent(noteID, contentID string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	n.AddContentID(contentID)

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

// UpdateContent updates the content of a note.
func (uc *NoteUsecase) UpdateContent(noteID, contentID, data string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	if err := n.UpdateContent(contentID, data); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

// DeleteContent deletes a content from a note.
func (uc *NoteUsecase) DeleteAndRemoveContent(noteID, contentID string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	if err := n.DeleteContent(contentID); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

func (uc *NoteUsecase) RemoveContent(noteID, contentID string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	if err := n.RemoveContentID(contentID); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

// TagNote adds a keyword to a note for a specific user.
func (uc *NoteUsecase) TagNote(noteID, userID, keywordStr string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	keyword, err := note.NewKeyword(keywordStr)
	if err != nil {
		return uc.mapDomainError(err)
	}

	n.AddKeyword(userID, keyword)

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

// UntagNote removes a keyword from a note for a specific user.
func (uc *NoteUsecase) UntagNote(noteID, userID, keywordStr string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}
	n := uc.mapper.ToDomain(notePO)

	keyword, err := note.NewKeyword(keywordStr)
	if err != nil {
		return uc.mapDomainError(err)
	}

	if err := n.RemoveKeyword(userID, keyword); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

// FindNotesByTag finds notes by a specific tag for a given user.
func (uc *NoteUsecase) FindNotesByKeyword(userID, keyword string) ([]*NoteDTO, error) {
	if userID == "" || keyword == "" {
		return []*NoteDTO{}, nil
	}

	notePOs, err := uc.repo.FindByKeywordForUser(userID, keyword)
	if err != nil {
		return nil, uc.mapRepositoryError(err)
	}

	var noteDTOs []*NoteDTO
	for _, notePO := range notePOs {
		n := uc.mapper.ToDomain(notePO)
		noteDTOs = append(noteDTOs, uc.mapper.toNoteDTO(n))
	}

	return noteDTOs, nil
}

// ShareNote shares a note with another user.
func (uc *NoteUsecase) ShareNote(noteID, ownerID, collaboratorID, permission string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}

	n := uc.mapper.ToDomain(notePO)

	permissionType, err := mapToDomainPermissionType(permission)
	if err != nil {
		return err
	}

	if err := n.AddCollaborator(ownerID, collaboratorID, permissionType); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

// GetAccessibleNotesForUser retrieves all notes that a user can access (owned or shared).
func (uc *NoteUsecase) GetAccessibleNotesForUser(userID string) ([]*NoteDTO, error) {
	notePOs, err := uc.repo.GetAccessibleNotesByUserID(userID)
	if err != nil {
		return nil, uc.mapRepositoryError(err)
	}

	var noteDTOs []*NoteDTO
	for _, notePO := range notePOs {
		n := uc.mapper.ToDomain(notePO)
		noteDTOs = append(noteDTOs, uc.mapper.toNoteDTO(n))
	}

	return noteDTOs, nil
}

// RevokeAccess revokes a collaborator's access to a note.
func (uc *NoteUsecase) RevokeAccess(noteID, ownerID, collaboratorID string) error {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return uc.mapRepositoryError(err)
	}

	n := uc.mapper.ToDomain(notePO)

	if err := n.RemoveCollaborator(ownerID, collaboratorID); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return uc.mapRepositoryError(err)
	}

	return nil
}

func (uc *NoteUsecase) mapRepositoryError(err error) error {
	switch {
	case errors.Is(err, noterepo.ErrNoteNotFound):
		return ErrNoteNotFound
	case errors.Is(err, noterepo.ErrNilNote):
		return ErrNilNote
	case errors.Is(err, noterepo.ErrNoteConflict):
		return ErrConflict
	default:
		return fmt.Errorf("an unexpected repository error occurred: %w", err)
	}
}

func (uc *NoteUsecase) mapDomainError(err error) error {
	switch {
	case errors.Is(err, note.ErrEmptyTitle):
		return ErrEmptyTitle
	case errors.Is(err, note.ErrContentNotFound):
		return ErrContentNotFound
	case errors.Is(err, note.ErrEmptyKeyword):
		return ErrEmptyKeyword
	case errors.Is(err, note.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, note.ErrKeywordNotFound):
		return ErrKeywordNotFound
	case errors.Is(err, note.ErrPermissionDenied):
		return ErrPermissionDenied
	default:
		return fmt.Errorf("an unexpected domain error occurred: %w", err)
	}
}

func mapToDomainContentType(ct ContentType) note.ContentType {
	switch ct {
	case TextContentType:
		return note.TextContentType
	case ImageContentType:
		return note.ImageContentType
	default:
		return note.TextContentType // Default to text if unknown
	}
}

func mapToDomainPermissionType(p string) (note.Permission, error) {
	switch p {
	case "read":
		return note.ReadOnly, nil
	case "read-write":
		return note.ReadWrite, nil
	default:
		return note.ReadOnly, ErrUnsupportedPermissionType
	}
}
