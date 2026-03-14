package format

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/grokify/fontscan/checker"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"toon", FormatTOON},
		{"TOON", FormatTOON},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"json-compact", FormatJSONCompact},
		{"JSON-COMPACT", FormatJSONCompact},
		{"", FormatTOON},
		{"unknown", FormatTOON},
	}

	for _, tt := range tests {
		got := ParseFormat(tt.input)
		if got != tt.expected {
			t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestNewFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	if f.Writer != &buf {
		t.Error("Writer not set correctly")
	}

	if f.Format != FormatJSON {
		t.Errorf("Format = %q, want %q", f.Format, FormatJSON)
	}
}

func TestWriteCheckResult_TOON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatTOON)

	result := &checker.CheckResult{
		File:         "test.md",
		Font:         "Test Font",
		Supported:    false,
		TotalChars:   100,
		UniqueChars:  50,
		MissingCount: 1,
		Missing: []checker.MissingChar{
			{Char: "→", Codepoint: "U+2192", Name: "RIGHTWARDS ARROW", Count: 5},
		},
	}

	if err := f.WriteCheckResult(result); err != nil {
		t.Fatalf("WriteCheckResult failed: %v", err)
	}

	output := buf.String()

	expectedLines := []string{
		"file: test.md",
		"font: Test Font",
		"supported: false",
		"totalChars: 100",
		"uniqueChars: 50",
		"missingCount: 1",
		"missing:",
		"char: →",
		"codepoint: U+2192",
		"name: RIGHTWARDS ARROW",
		"count: 5",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("Output missing %q", line)
		}
	}
}

func TestWriteCheckResult_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	result := &checker.CheckResult{
		File:         "test.md",
		Font:         "Test Font",
		Supported:    true,
		TotalChars:   100,
		UniqueChars:  50,
		MissingCount: 0,
	}

	if err := f.WriteCheckResult(result); err != nil {
		t.Fatalf("WriteCheckResult failed: %v", err)
	}

	// Verify valid JSON
	var parsed checker.CheckResult
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if parsed.File != "test.md" {
		t.Errorf("File = %q, want %q", parsed.File, "test.md")
	}
}

func TestWriteCheckResult_JSONCompact(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSONCompact)

	result := &checker.CheckResult{
		Font:      "Test Font",
		Supported: true,
	}

	if err := f.WriteCheckResult(result); err != nil {
		t.Fatalf("WriteCheckResult failed: %v", err)
	}

	output := buf.String()

	// Should be single line (plus newline)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("JSON-compact should be single line, got %d lines", len(lines))
	}
}

func TestWriteRecommendResult_TOON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatTOON)

	result := &checker.RecommendResult{
		File:        "test.md",
		TotalChars:  100,
		UniqueChars: 50,
		FullSupport: []checker.FontMatch{
			{Name: "Font A", Path: "/path/a.ttf"},
		},
		PartialSupport: []checker.PartialMatch{
			{
				FontMatch:    checker.FontMatch{Name: "Font B", Path: "/path/b.ttf"},
				Percentage:   90.0,
				MissingCount: 5,
			},
		},
	}

	if err := f.WriteRecommendResult(result); err != nil {
		t.Fatalf("WriteRecommendResult failed: %v", err)
	}

	output := buf.String()

	expectedLines := []string{
		"file: test.md",
		"totalChars: 100",
		"fullSupportCount: 1",
		"fullSupport:",
		"name: Font A",
		"partialSupport:",
		"name: Font B",
		"percentage: 90.0%",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("Output missing %q", line)
		}
	}
}

func TestWriteFontList_TOON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatTOON)

	fonts := []FontListEntry{
		{Name: "Font A", Path: "/path/a.ttf"},
		{Name: "Font B", Path: "/path/b.ttf"},
	}

	if err := f.WriteFontList(fonts); err != nil {
		t.Fatalf("WriteFontList failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "count: 2") {
		t.Error("Output missing count")
	}

	if !strings.Contains(output, "Font A") {
		t.Error("Output missing Font A")
	}

	if !strings.Contains(output, "Font B") {
		t.Error("Output missing Font B")
	}
}

func TestWriteFontList_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	fonts := []FontListEntry{
		{Name: "Font A", Path: "/path/a.ttf"},
	}

	if err := f.WriteFontList(fonts); err != nil {
		t.Fatalf("WriteFontList failed: %v", err)
	}

	var parsed []FontListEntry
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Errorf("Expected 1 font, got %d", len(parsed))
	}
}
