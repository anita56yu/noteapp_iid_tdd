package noterepo

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
	Contents      []ContentPO
	Keywords      map[string][]string
	Collaborators map[string]string
}
