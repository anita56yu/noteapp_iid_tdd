package domain

import (
	"testing"
)

func TestNewKeyword_ValidCreation(t *testing.T) {
	text := "test-keyword"
	keyword, err := NewKeyword(text)

	if err != nil {
		t.Fatalf("Failed to create a valid keyword: %v", err)
	}

	if keyword.String() != text {
		t.Errorf("Expected keyword text to be '%s', but got '%s'", text, keyword.String())
	}
}

func TestNewKeyword_EmptyText(t *testing.T) {
	_, err := NewKeyword("")
	if err == nil {
		t.Fatal("Expected an error for empty keyword text, but got nil")
	}
}
