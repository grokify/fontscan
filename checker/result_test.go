package checker

import (
	"testing"
)

func TestNewMissingChar(t *testing.T) {
	tests := []struct {
		r        rune
		count    int
		wantChar string
		wantCode string
	}{
		{'→', 5, "→", "U+2192"},
		{'A', 1, "A", "U+0041"},
		{'─', 3, "─", "U+2500"},
		{'│', 2, "│", "U+2502"},
	}

	for _, tt := range tests {
		mc := NewMissingChar(tt.r, tt.count)

		if mc.Char != tt.wantChar {
			t.Errorf("NewMissingChar(%U).Char = %s, want %s", tt.r, mc.Char, tt.wantChar)
		}

		if mc.Codepoint != tt.wantCode {
			t.Errorf("NewMissingChar(%U).Codepoint = %s, want %s", tt.r, mc.Codepoint, tt.wantCode)
		}

		if mc.Count != tt.count {
			t.Errorf("NewMissingChar(%U).Count = %d, want %d", tt.r, mc.Count, tt.count)
		}

		// Name should be set from Unicode name table
		if mc.Name == "" {
			t.Errorf("NewMissingChar(%U).Name is empty, expected Unicode name", tt.r)
		}
	}
}

func TestCheckResult_Supported(t *testing.T) {
	result := &CheckResult{
		Font:         "Test Font",
		Supported:    true,
		TotalChars:   100,
		UniqueChars:  50,
		MissingCount: 0,
		Missing:      nil,
	}

	if !result.Supported {
		t.Error("Expected Supported = true when no missing chars")
	}

	if result.MissingCount != 0 {
		t.Errorf("Expected MissingCount = 0, got %d", result.MissingCount)
	}
}

func TestCheckResult_NotSupported(t *testing.T) {
	result := &CheckResult{
		Font:         "Test Font",
		Supported:    false,
		TotalChars:   100,
		UniqueChars:  50,
		MissingCount: 3,
		Missing: []MissingChar{
			{Char: "→", Codepoint: "U+2192", Name: "RIGHTWARDS ARROW", Count: 5},
			{Char: "─", Codepoint: "U+2500", Name: "BOX DRAWINGS LIGHT HORIZONTAL", Count: 10},
			{Char: "│", Codepoint: "U+2502", Name: "BOX DRAWINGS LIGHT VERTICAL", Count: 8},
		},
	}

	if result.Supported {
		t.Error("Expected Supported = false when missing chars present")
	}

	if result.MissingCount != 3 {
		t.Errorf("Expected MissingCount = 3, got %d", result.MissingCount)
	}

	if len(result.Missing) != 3 {
		t.Errorf("Expected 3 missing chars, got %d", len(result.Missing))
	}
}

func TestRecommendResult_FullSupport(t *testing.T) {
	result := &RecommendResult{
		File:        "test.md",
		TotalChars:  100,
		UniqueChars: 50,
		FullSupport: []FontMatch{
			{Name: "Font A", Path: "/path/a.ttf"},
			{Name: "Font B", Path: "/path/b.ttf"},
		},
	}

	if len(result.FullSupport) != 2 {
		t.Errorf("Expected 2 full support fonts, got %d", len(result.FullSupport))
	}
}

func TestRecommendResult_PartialSupport(t *testing.T) {
	result := &RecommendResult{
		File:        "test.md",
		TotalChars:  100,
		UniqueChars: 50,
		PartialSupport: []PartialMatch{
			{
				FontMatch:    FontMatch{Name: "Font C", Path: "/path/c.ttf"},
				Supported:    45,
				Total:        50,
				Percentage:   90.0,
				MissingCount: 5,
			},
		},
	}

	if len(result.PartialSupport) != 1 {
		t.Errorf("Expected 1 partial support font, got %d", len(result.PartialSupport))
	}

	if result.PartialSupport[0].Percentage != 90.0 {
		t.Errorf("Expected 90%% support, got %.1f%%", result.PartialSupport[0].Percentage)
	}
}
