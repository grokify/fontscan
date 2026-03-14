package pandoc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckResult_HasFontSettings(t *testing.T) {
	tests := []struct {
		name     string
		result   CheckResult
		expected bool
	}{
		{
			name:     "empty settings",
			result:   CheckResult{FontSettings: map[string]string{}},
			expected: false,
		},
		{
			name:     "nil settings",
			result:   CheckResult{FontSettings: nil},
			expected: false,
		},
		{
			name: "with settings",
			result: CheckResult{
				FontSettings: map[string]string{"mainfont": "Arial"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.HasFontSettings(); got != tt.expected {
				t.Errorf("HasFontSettings() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCheckOptions(t *testing.T) {
	opts := CheckOptions{
		Paths:        []string{"/custom/path"},
		WarnDefaults: true,
	}

	if len(opts.Paths) != 1 || opts.Paths[0] != "/custom/path" {
		t.Errorf("CheckOptions.Paths = %v, expected [/custom/path]", opts.Paths)
	}
	if !opts.WarnDefaults {
		t.Error("CheckOptions.WarnDefaults = false, expected true")
	}
}

func TestCheckDocument_NoFrontmatter(t *testing.T) {
	// Create a temp file without frontmatter
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")

	content := "# Hello World\n\nThis is a test document."
	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	result, err := CheckDocument(tmpFile, CheckOptions{})
	if err != nil {
		t.Fatalf("CheckDocument() error = %v", err)
	}

	if result.HasFontSettings() {
		t.Error("Expected no font settings for document without frontmatter")
	}
}

func TestCheckDocument_WithFrontmatter(t *testing.T) {
	// Create a temp file with frontmatter but invalid font
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")

	content := `---
title: Test
mainfont: "NonExistentFont12345"
---

# Hello World
`
	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	result, err := CheckDocument(tmpFile, CheckOptions{})
	if err != nil {
		t.Fatalf("CheckDocument() error = %v", err)
	}

	if !result.HasFontSettings() {
		t.Error("Expected font settings for document with frontmatter")
	}

	if result.FontSettings["mainfont"] != "NonExistentFont12345" {
		t.Errorf("FontSettings[mainfont] = %q, expected %q",
			result.FontSettings["mainfont"], "NonExistentFont12345")
	}

	// Should have an error for non-existent font
	if len(result.Errors) == 0 {
		t.Error("Expected error for non-existent font")
	}

	if !result.HasErrors {
		t.Error("HasErrors should be true for non-existent font")
	}
}

func TestCheckDocument_FileNotFound(t *testing.T) {
	_, err := CheckDocument("/nonexistent/file.md", CheckOptions{})
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
