package contentrepo

// ContentRepository defines the interface for content persistence.
type ContentRepository interface {
	Save(c *ContentPO) error
	GetByID(id string) (*ContentPO, error)
	GetAllByNoteID(noteID string) ([]*ContentPO, error)
	Delete(id string) error
	DeleteAllByNoteID(noteID string) error
}
