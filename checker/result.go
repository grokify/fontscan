// Package checker provides text-to-font compatibility checking.
package checker

import (
	"fmt"
	"unicode"

	"golang.org/x/text/unicode/runenames"
)

// CheckResult represents the result of checking text against a font.
type CheckResult struct {
	File         string        `json:"file,omitempty"`
	Font         string        `json:"font"`
	FontPath     string        `json:"fontPath,omitempty"`
	Supported    bool          `json:"supported"`
	TotalChars   int           `json:"totalChars"`
	UniqueChars  int           `json:"uniqueChars"`
	MissingCount int           `json:"missingCount"`
	Missing      []MissingChar `json:"missing,omitempty"`
}

// MissingChar represents a character not supported by a font.
type MissingChar struct {
	Char      string `json:"char"`
	Codepoint string `json:"codepoint"`
	Name      string `json:"name,omitempty"`
	Count     int    `json:"count"`
}

// NewMissingChar creates a MissingChar from a rune and its occurrence count.
func NewMissingChar(r rune, count int) MissingChar {
	return MissingChar{
		Char:      string(r),
		Codepoint: fmt.Sprintf("U+%04X", r),
		Name:      runenames.Name(r),
		Count:     count,
	}
}

// RecommendResult represents fonts recommended for a document.
type RecommendResult struct {
	File          string            `json:"file"`
	TotalChars    int               `json:"totalChars"`
	UniqueChars   int               `json:"uniqueChars"`
	FullSupport   []FontMatch       `json:"fullSupport,omitempty"`
	PartialSupport []PartialMatch   `json:"partialSupport,omitempty"`
}

// FontMatch represents a font that fully supports a document.
type FontMatch struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Family string `json:"family,omitempty"`
	Style  string `json:"style,omitempty"`
}

// PartialMatch represents a font with partial support.
type PartialMatch struct {
	FontMatch
	Supported    int     `json:"supported"`
	Total        int     `json:"total"`
	Percentage   float64 `json:"percentage"`
	MissingCount int     `json:"missingCount"`
}

// IsPrintable returns true if the character is a printable character
// (not a control character or whitespace).
func IsPrintable(r rune) bool {
	return unicode.IsPrint(r) && !unicode.IsSpace(r)
}

// IsSignificant returns true if the character should be checked for font support.
// This excludes ASCII control characters but includes special Unicode characters.
func IsSignificant(r rune) bool {
	// Skip ASCII control characters and basic whitespace
	if r < 0x20 {
		return false
	}
	// Skip basic space, but keep non-breaking space and other special spaces
	if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
		return false
	}
	return true
}
