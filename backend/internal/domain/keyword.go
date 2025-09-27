package domain

import "errors"

// Keyword is a value object representing a keyword.
type Keyword struct {
	text string
}

// NewKeyword creates a new Keyword. It returns an error if the text is empty.
func NewKeyword(text string) (Keyword, error) {
	if text == "" {
		return Keyword{}, errors.New("keyword text cannot be empty")
	}
	return Keyword{text: text}, nil
}

// String returns the keyword text.
func (k Keyword) String() string {
	return k.text
}
