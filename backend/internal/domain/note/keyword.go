package note

import "errors"

// ErrEmptyKeyword is returned when a keyword is created with empty text.
var ErrEmptyKeyword = errors.New("keyword text cannot be empty")

// Keyword is a value object representing a keyword.
type Keyword struct {
	text string
}

// NewKeyword creates a new Keyword. It returns an error if the text is empty.
func NewKeyword(text string) (Keyword, error) {
	if text == "" {
		return Keyword{}, ErrEmptyKeyword
	}
	return Keyword{text: text}, nil
}

// String returns the keyword text.
func (k Keyword) String() string {
	return k.text
}