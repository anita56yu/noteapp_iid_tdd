package usecase

// NoteDTO represents a note for external consumers.
type NoteDTO struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
