package contentrepo

// ContentPO represents the persistence model for a Content object.
type ContentPO struct {
	ID      string
	NoteID  string
	Data    string
	Type    string
	Version int
}
