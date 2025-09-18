package repository

// NoteRepository defines the interface for note persistence.
type NoteRepository interface {
	Save(note *NotePO) error
	FindByID(id string) (*NotePO, error)
	Delete(id string) error
}
