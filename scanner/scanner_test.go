package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSystemFontPaths(t *testing.T) {
	paths := SystemFontPaths()

	if len(paths) == 0 {
		t.Error("Expected some system font paths")
	}

	// Verify paths look reasonable for the current OS
	switch runtime.GOOS {
	case "darwin":
		found := false
		for _, p := range paths {
			if p == "/System/Library/Fonts" || p == "/Library/Fonts" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected macOS font paths to include /System/Library/Fonts or /Library/Fonts")
		}
	case "linux":
		found := false
		for _, p := range paths {
			if p == "/usr/share/fonts" || p == "/usr/local/share/fonts" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected Linux font paths to include /usr/share/fonts")
		}
	}
}

func TestNewScanner(t *testing.T) {
	s := NewScanner()

	if len(s.Paths) == 0 {
		t.Error("NewScanner should have default paths")
	}
}

func TestNewScannerWithPaths(t *testing.T) {
	paths := []string{"/custom/path1", "/custom/path2"}
	s := NewScannerWithPaths(paths)

	if len(s.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(s.Paths))
	}
}

func TestScanner_Scan(t *testing.T) {
	s := NewScanner()

	fonts, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should find at least some fonts on a typical system
	if len(fonts) == 0 {
		t.Log("Warning: No fonts found. This may be expected in some CI environments.")
	}
}

func TestScanner_Scan_NonexistentPath(t *testing.T) {
	s := NewScannerWithPaths([]string{"/nonexistent/path"})

	fonts, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan should not error on nonexistent path: %v", err)
	}

	if len(fonts) != 0 {
		t.Errorf("Expected 0 fonts from nonexistent path, got %d", len(fonts))
	}
}

func TestIsFontFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"font.ttf", true},
		{"font.TTF", true},
		{"font.otf", true},
		{"font.OTF", true},
		{"font.ttc", true},
		{"font.TTC", true},
		{"font.txt", false},
		{"font.png", false},
		{"font", false},
		{".ttf", true}, // Hidden file with font extension
	}

	for _, tt := range tests {
		got := isFontFile(tt.path)
		if got != tt.expected {
			t.Errorf("isFontFile(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestScanner_FindByName(t *testing.T) {
	s := NewScanner()

	// Try to find a common font
	_, err := s.FindByName("Arial")
	// This may fail if Arial is not installed, which is fine
	if err != nil {
		t.Logf("Arial not found: %v (this may be expected)", err)
	}
}

func TestScanner_FindByName_NotFound(t *testing.T) {
	s := NewScanner()

	_, err := s.FindByName("NonExistentFont12345")
	if err == nil {
		t.Error("Expected error for nonexistent font")
	}
}

func TestFindByPath(t *testing.T) {
	// Try to find a test font
	candidates := []string{
		"/System/Library/Fonts/Helvetica.ttc",
		"/System/Library/Fonts/Monaco.ttf",
		"/Library/Fonts/Arial.ttf",
	}

	var fontPath string
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			fontPath = path
			break
		}
	}

	if fontPath == "" {
		t.Skip("No test font found")
	}

	f, err := FindByPath(fontPath)
	if err != nil {
		t.Fatalf("FindByPath failed: %v", err)
	}

	if f.Path != fontPath {
		t.Errorf("Expected path %s, got %s", fontPath, f.Path)
	}
}

func TestSupportedExtensions(t *testing.T) {
	expected := []string{".ttf", ".otf", ".ttc"}

	if len(SupportedExtensions) != len(expected) {
		t.Errorf("Expected %d extensions, got %d", len(expected), len(SupportedExtensions))
	}

	for _, ext := range expected {
		found := false
		for _, supported := range SupportedExtensions {
			if supported == ext {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s in SupportedExtensions", ext)
		}
	}
}

func TestScanner_HomeExpansion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	// Create a temp font directory in home
	testDir := filepath.Join(home, ".fontscan-test-"+t.Name())
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(testDir) }()

	// Use ~ path
	s := NewScannerWithPaths([]string{"~/.fontscan-test-" + t.Name()})

	// Should not error, even if no fonts found
	_, err = s.Scan()
	if err != nil {
		t.Errorf("Scan with ~ path failed: %v", err)
	}
}
