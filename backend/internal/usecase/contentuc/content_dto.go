package contentuc

// ContentDTO represents the data transfer object for a Content.
type ContentDTO struct {
	ID      string `json:"id"`
	NoteID  string `json:"noteId"`
	Data    string `json:"data"`
	Type    string `json:"type"`
	Version int    `json:"version"`
}
