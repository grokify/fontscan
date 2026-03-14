package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/fontscan/font"
)

func findTestFont(t *testing.T) string {
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

func TestCheckText(t *testing.T) {
	path := findTestFont(t)

	f, err := font.LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	result := CheckText("Hello World", f)

	if result.Font == "" {
		t.Error("Expected Font to be set")
	}

	if result.TotalChars == 0 {
		t.Error("Expected TotalChars > 0")
	}

	if result.UniqueChars == 0 {
		t.Error("Expected UniqueChars > 0")
	}
}

func TestCheckText_Empty(t *testing.T) {
	path := findTestFont(t)

	f, err := font.LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	result := CheckText("", f)

	if result.TotalChars != 0 {
		t.Errorf("Expected TotalChars = 0 for empty string, got %d", result.TotalChars)
	}

	if !result.Supported {
		t.Error("Empty text should be supported")
	}
}

func TestCheckText_WhitespaceOnly(t *testing.T) {
	path := findTestFont(t)

	f, err := font.LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	result := CheckText("   \t\n\r   ", f)

	if result.TotalChars != 0 {
		t.Errorf("Expected TotalChars = 0 for whitespace-only, got %d", result.TotalChars)
	}

	if !result.Supported {
		t.Error("Whitespace-only text should be supported")
	}
}

func TestCheckFile(t *testing.T) {
	path := findTestFont(t)

	f, err := font.LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	// Create a temp file
	tmpfile, err := os.CreateTemp("", "fontscan-test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.WriteString("Hello World 123"); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := CheckFile(tmpfile.Name(), f)
	if err != nil {
		t.Fatalf("CheckFile failed: %v", err)
	}

	if result.File != tmpfile.Name() {
		t.Errorf("Expected File = %s, got %s", tmpfile.Name(), result.File)
	}
}

func TestCheckFile_NotFound(t *testing.T) {
	path := findTestFont(t)

	f, err := font.LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	_, err = CheckFile("/nonexistent/file.txt", f)
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestNewChecker(t *testing.T) {
	c := NewChecker()

	if c.Scanner == nil {
		t.Error("Expected Scanner to be initialized")
	}
}

func TestNewCheckerWithPaths(t *testing.T) {
	paths := []string{"/custom/path"}
	c := NewCheckerWithPaths(paths)

	if c.Scanner == nil {
		t.Error("Expected Scanner to be initialized")
	}
}

func TestIsSignificant(t *testing.T) {
	tests := []struct {
		r        rune
		expected bool
	}{
		{'A', true},
		{'a', true},
		{'0', true},
		{'!', true},
		{'→', true},
		{' ', false},     // Space
		{'\t', false},    // Tab
		{'\n', false},    // Newline
		{'\r', false},    // Carriage return
		{0x00, false},    // Null
		{0x1F, false},    // Unit separator (control char)
		{'\u00A0', true}, // Non-breaking space (should be checked)
	}

	for _, tt := range tests {
		got := IsSignificant(tt.r)
		if got != tt.expected {
			t.Errorf("IsSignificant(%U) = %v, want %v", tt.r, got, tt.expected)
		}
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		r        rune
		expected bool
	}{
		{'A', true},
		{'→', true},
		{' ', false},  // Space
		{'\t', false}, // Tab
		{'\n', false}, // Newline
		{0x00, false}, // Null
	}

	for _, tt := range tests {
		got := IsPrintable(tt.r)
		if got != tt.expected {
			t.Errorf("IsPrintable(%U) = %v, want %v", tt.r, got, tt.expected)
		}
	}
}

func TestMissingCharSorted(t *testing.T) {
	path := findTestFont(t)

	f, err := font.LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	// Text with characters in non-ascending codepoint order
	result := CheckText("z→a─", f)

	// Verify missing chars are sorted by codepoint
	for i := 1; i < len(result.Missing); i++ {
		if result.Missing[i].Codepoint < result.Missing[i-1].Codepoint {
			t.Errorf("Missing chars not sorted: %s before %s",
				result.Missing[i-1].Codepoint, result.Missing[i].Codepoint)
		}
	}
}
