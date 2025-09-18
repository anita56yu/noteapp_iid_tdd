package repository

// ContentPO represents the persistent state of a content block.
type ContentPO struct {
	ID   string
	Type string
	Data string
}

// NotePO represents the persistent state of a note.
type NotePO struct {
	ID       string
	Title    string
	Contents []ContentPO
}
