package noterepo

import "time"

// ContentPO represents the persistent state of a content block.
type ContentPO struct {
	ID   string
	Type string
	Data string
}

// NotePO represents the persistent state of a note.
type NotePO struct {
	ID            string
	OwnerID       string
	Title         string
	Version       int
	ContentIDs    []string
	Keywords      map[string][]string
	Collaborators map[string]string
}

type NoteEventPO struct {
	EventType  string                 `json:"event_type"`
	EventID    string                 `json:"event_id"`
	OccurredAt time.Time              `json:"occurred_at"`
	NoteID     string                 `json:"note_id"`
	Payload    map[string]interface{} `json:"payload"`
}
