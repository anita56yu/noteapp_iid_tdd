package usecase

// ContentDTO represents a content block for external consumers.
type ContentDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// NoteDTO represents a note for external consumers.
type NoteDTO struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	Contents []ContentDTO `json:"contents"`
}