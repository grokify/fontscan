package pandoc

import (
	"os"
	"testing"
)

func TestParseFrontmatterString(t *testing.T) {
	yaml := `
mainfont: "Helvetica Neue"
sansfont: Arial
monofont: "Menlo"
CJKmainfont: "Hiragino Sans"
`
	settings, err := ParseFrontmatterString(yaml)
	if err != nil {
		t.Fatalf("ParseFrontmatterString failed: %v", err)
	}

	if settings.MainFont != "Helvetica Neue" {
		t.Errorf("MainFont = %q, want %q", settings.MainFont, "Helvetica Neue")
	}

	if settings.SansFont != "Arial" {
		t.Errorf("SansFont = %q, want %q", settings.SansFont, "Arial")
	}

	if settings.MonoFont != "Menlo" {
		t.Errorf("MonoFont = %q, want %q", settings.MonoFont, "Menlo")
	}

	if settings.CJKMainFont != "Hiragino Sans" {
		t.Errorf("CJKMainFont = %q, want %q", settings.CJKMainFont, "Hiragino Sans")
	}
}

func TestFontSettings_Fonts(t *testing.T) {
	settings := &FontSettings{
		MainFont: "Helvetica",
		SansFont: "Arial",
		MonoFont: "Menlo",
	}

	fonts := settings.Fonts()

	if len(fonts) != 3 {
		t.Errorf("Expected 3 fonts, got %d", len(fonts))
	}

	// Check deduplication
	settings.TitleFont = "Helvetica" // Same as MainFont
	fonts = settings.Fonts()

	if len(fonts) != 3 {
		t.Errorf("Expected 3 fonts (deduplicated), got %d", len(fonts))
	}
}

func TestFontSettings_FontMap(t *testing.T) {
	settings := &FontSettings{
		MainFont:  "Helvetica",
		MonoFont:  "Menlo",
		TitleFont: "Arial",
	}

	m := settings.FontMap()

	if m["mainfont"] != "Helvetica" {
		t.Errorf("mainfont = %q, want %q", m["mainfont"], "Helvetica")
	}

	if m["monofont"] != "Menlo" {
		t.Errorf("monofont = %q, want %q", m["monofont"], "Menlo")
	}

	if m["titlefont"] != "Arial" {
		t.Errorf("titlefont = %q, want %q", m["titlefont"], "Arial")
	}

	// sansfont should not be in map (empty)
	if _, ok := m["sansfont"]; ok {
		t.Error("sansfont should not be in map when empty")
	}
}

func TestParseFrontmatter(t *testing.T) {
	content := `---
title: Test Document
mainfont: "Times New Roman"
monofont: "Courier New"
---

# Hello World

This is the document body.
`
	// Create temp file
	tmpfile, err := os.CreateTemp("", "pandoc-test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	settings, err := ParseFrontmatter(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}

	if settings.MainFont != "Times New Roman" {
		t.Errorf("MainFont = %q, want %q", settings.MainFont, "Times New Roman")
	}

	if settings.MonoFont != "Courier New" {
		t.Errorf("MonoFont = %q, want %q", settings.MonoFont, "Courier New")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := `# Hello World

No frontmatter here.
`
	tmpfile, err := os.CreateTemp("", "pandoc-test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	settings, err := ParseFrontmatter(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}

	// Should return empty settings, not error
	if len(settings.Fonts()) != 0 {
		t.Errorf("Expected 0 fonts for file without frontmatter, got %d", len(settings.Fonts()))
	}
}

func TestParseFrontmatter_TripleDotEnding(t *testing.T) {
	content := `---
mainfont: Helvetica
...

Body text.
`
	tmpfile, err := os.CreateTemp("", "pandoc-test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	settings, err := ParseFrontmatter(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}

	if settings.MainFont != "Helvetica" {
		t.Errorf("MainFont = %q, want %q", settings.MainFont, "Helvetica")
	}
}
