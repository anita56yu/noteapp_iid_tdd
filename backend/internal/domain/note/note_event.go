package note

import (
	"time"
)

type NoteEvent struct {
	EventType  string
	EventID    string
	NoteID     string
	Payload    map[string]interface{}
	OccurredAt time.Time
}
