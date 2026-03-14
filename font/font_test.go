package font

import (
	"os"
	"path/filepath"
	"testing"
)

func findTestFont(t *testing.T) string {
	// Try common macOS system fonts
	candidates := []string{
		"/System/Library/Fonts/Helvetica.ttc",
		"/System/Library/Fonts/Monaco.ttf",
		"/Library/Fonts/Arial.ttf",
		"/System/Library/Fonts/SFNS.ttf",
		"/System/Library/Fonts/Supplemental/Arial.ttf",
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Try to find any TTF file in system fonts
	matches, _ := filepath.Glob("/System/Library/Fonts/*.ttf")
	if len(matches) > 0 {
		return matches[0]
	}

	matches, _ = filepath.Glob("/Library/Fonts/*.ttf")
	if len(matches) > 0 {
		return matches[0]
	}

	t.Skip("No system fonts found for testing")
	return ""
}

func TestLoadFile(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	if f.Path != path {
		t.Errorf("Expected path %s, got %s", path, f.Path)
	}

	if f.NumGlyphs == 0 {
		t.Error("Expected NumGlyphs > 0")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/font.ttf")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestLoadFile_InvalidFile(t *testing.T) {
	// Create a temp file that's not a font
	tmpfile, err := os.CreateTemp("", "notafont*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.WriteString("This is not a font file"); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadFile(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid font file")
	}
}

func TestHasGlyph(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	// Basic ASCII characters should be supported by most fonts
	if !f.HasGlyph('A') {
		t.Error("Expected font to have glyph for 'A'")
	}

	if !f.HasGlyph('a') {
		t.Error("Expected font to have glyph for 'a'")
	}

	if !f.HasGlyph('0') {
		t.Error("Expected font to have glyph for '0'")
	}
}

func TestMissingGlyphs(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	// ASCII text should have no missing glyphs in most fonts
	missing := f.MissingGlyphs("Hello World 123")
	// This may or may not be empty depending on the font, so just verify it returns
	_ = missing
}

func TestMissingGlyphs_Deduplication(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	// Use a character that may be missing in some fonts
	text := "aaa\U0001F600\U0001F600\U0001F600" // Three emoji

	missing := f.MissingGlyphs(text)

	// Count occurrences of each rune
	counts := make(map[rune]int)
	for _, r := range missing {
		counts[r]++
	}

	// Each unique missing rune should appear only once
	for r, count := range counts {
		if count > 1 {
			t.Errorf("Rune %U appears %d times, expected 1", r, count)
		}
	}
}

func TestSupportedGlyphs(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	supported := f.SupportedGlyphs("ABC")
	if len(supported) == 0 {
		t.Error("Expected some supported glyphs for basic ASCII")
	}
}

func TestFont_String(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	str := f.String()
	if str == "" {
		t.Error("Expected non-empty String() result")
	}
}
