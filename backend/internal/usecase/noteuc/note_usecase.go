package noteuc

import (
	"errors"
	"fmt"
	"noteapp/internal/domain/note"
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

// ErrIndexOutOfBounds is returned when an index is out of bounds in content insertion.
var ErrIndexOutOfBounds = errors.New("index out of bounds")

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
	repo   NoteRepository
	mapper *NoteMapper
}

// NewNoteUsecase creates a new NoteUsecase.
func NewNoteUsecase(repo NoteRepository) *NoteUsecase {
	return &NoteUsecase{repo: repo, mapper: NewNoteMapper()}
}

// HasWritePermission checks if the user has write permission for the given note.
func (uc *NoteUsecase) HasWritePermission(noteID, userID string) bool {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		// If note is not found or other repo error, assume no write permission.
		return false
	}
	n := uc.mapper.ToDomain(notePO)
	return n.HasWritePermission(userID)
}

// CreateNote creates a new note.
func (uc *NoteUsecase) CreateNote(id, title, ownerID string) (string, error) {
	n, err := note.NewNote(id, title, ownerID)
	if err != nil {
		return "", uc.mapDomainError(err)
	}

	notePO := uc.mapper.ToPO(n)
	eventPOs := uc.mapper.toEventPOs(n)
	if err := uc.repo.SaveWithEvent(notePO, eventPOs); err != nil {
		return "", err
	}

	return n.ID, nil
}

// GetNoteByID retrieves a note by its ID.
func (uc *NoteUsecase) GetNoteByID(id, userID string) (*NoteDTO, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	notePO, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	n := uc.mapper.ToDomain(notePO)
	if !n.HasReadPermission(userID) {
		return nil, ErrPermissionDenied
	}
	return uc.mapper.toNoteDTO(n), nil
}

// DeleteNote deletes a note by its ID.
func (uc *NoteUsecase) DeleteNote(id, userID string, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(id, version)
	if err != nil {
		return err
	}

	n := uc.mapper.ToDomain(notePO)
	if !n.IsOwner(userID) {
		return ErrPermissionDenied
	}

	if err := uc.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (uc *NoteUsecase) AddContent(noteID, contentID string, index, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
	}

	n := uc.mapper.ToDomain(notePO)

	if err := n.AddContentID(contentID, index); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return err
	}

	return nil
}

// ChangeTitle updates the title of a note.
func (uc *NoteUsecase) ChangeTitle(noteID, userID, newTitle string, version int) error {
	if noteID == "" {
		return ErrInvalidID
	}
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
	}

	n := uc.mapper.ToDomain(notePO)
	if err := n.ChangeTitle(userID, newTitle); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return err
	}

	return nil
}

func (uc *NoteUsecase) RemoveContent(noteID, contentID string, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
	}
	n := uc.mapper.ToDomain(notePO)

	if err := n.RemoveContentID(contentID); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return err
	}

	return nil
}

// TagNote adds a keyword to a note for a specific user.
func (uc *NoteUsecase) TagNote(noteID, userID, keywordStr string, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
	}
	n := uc.mapper.ToDomain(notePO)

	keyword, err := note.NewKeyword(keywordStr)
	if err != nil {
		return uc.mapDomainError(err)
	}

	n.AddKeyword(userID, keyword)

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return err
	}

	return nil
}

// UntagNote removes a keyword from a note for a specific user.
func (uc *NoteUsecase) UntagNote(noteID, userID, keywordStr string, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
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
		return err
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
		return nil, err
	}

	var noteDTOs []*NoteDTO
	for _, notePO := range notePOs {
		n := uc.mapper.ToDomain(notePO)
		noteDTOs = append(noteDTOs, uc.mapper.toNoteDTO(n))
	}

	return noteDTOs, nil
}

// ShareNote shares a note with another user.
func (uc *NoteUsecase) ShareNote(noteID, ownerID, collaboratorID, permission string, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

// GetAccessibleNotesForUser retrieves all notes that a user can access (owned or shared).
func (uc *NoteUsecase) GetAccessibleNotesForUser(userID string) ([]*NoteDTO, error) {
	notePOs, err := uc.repo.GetAccessibleNotesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var noteDTOs []*NoteDTO
	for _, notePO := range notePOs {
		n := uc.mapper.ToDomain(notePO)
		noteDTOs = append(noteDTOs, uc.mapper.toNoteDTO(n))
	}

	return noteDTOs, nil
}

// RevokeAccess revokes a collaborator's access to a note.
func (uc *NoteUsecase) RevokeAccess(noteID, ownerID, collaboratorID string, version int) error {
	notePO, err := uc.getNotePOAndCheckVersion(noteID, version)
	if err != nil {
		return err
	}

	n := uc.mapper.ToDomain(notePO)

	if err := n.RemoveCollaborator(ownerID, collaboratorID); err != nil {
		return uc.mapDomainError(err)
	}

	updatedNotePO := uc.mapper.ToPO(n)
	if err := uc.repo.Save(updatedNotePO); err != nil {
		return err
	}

	return nil
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
	case errors.Is(err, note.ErrIndexOutOfBounds):
		return ErrIndexOutOfBounds
	default:
		return fmt.Errorf("an unexpected domain error occurred: %w", err)
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

func (uc *NoteUsecase) getNotePOAndCheckVersion(noteID string, version int) (*NotePO, error) {
	notePO, err := uc.repo.FindByID(noteID)
	if err != nil {
		return nil, err
	}

	if notePO.Version != version {
		return nil, ErrConflict
	}

	return notePO, nil
}
