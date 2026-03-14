package pandoc

import (
	"testing"
)

func TestIsDefaultFont(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"lmmono10-regular", true},
		{"lmroman10-regular", true},
		{"Latin Modern Mono", true},
		{"texgyrecursor-regular", true},
		{"DejaVuSansMono", true},
		{"Arial", false},
		{"Helvetica", false},
		{"RandomFont", false},
	}

	for _, tt := range tests {
		got := IsDefaultFont(tt.name)
		if got != tt.expected {
			t.Errorf("IsDefaultFont(%q) = %v, want %v", tt.name, got, tt.expected)
		}
	}
}

func TestGetFontInfo(t *testing.T) {
	info := GetFontInfo("lmmono10-regular")
	if info == nil {
		t.Fatal("Expected to find info for lmmono10-regular")
	}

	if info.Family != "lmmono" {
		t.Errorf("Family = %q, want %q", info.Family, "lmmono")
	}

	if info.Name != "Latin Modern Mono" {
		t.Errorf("Name = %q, want %q", info.Name, "Latin Modern Mono")
	}

	// Not found
	info = GetFontInfo("NonExistentFont")
	if info != nil {
		t.Error("Expected nil for unknown font")
	}
}

func TestAllDefaultFontNames(t *testing.T) {
	names := AllDefaultFontNames()

	if len(names) == 0 {
		t.Error("Expected some default font names")
	}

	// Check that lmmono10-regular is in the list
	found := false
	for _, n := range names {
		if n == "lmmono10-regular" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected lmmono10-regular in default font names")
	}
}

func TestDefaultFonts_NotEmpty(t *testing.T) {
	if len(DefaultFonts) == 0 {
		t.Error("DefaultFonts should not be empty")
	}

	for _, f := range DefaultFonts {
		if f.Name == "" {
			t.Error("Font Name should not be empty")
		}
		if f.Family == "" {
			t.Error("Font Family should not be empty")
		}
		if len(f.Variants) == 0 {
			t.Errorf("Font %s should have variants", f.Name)
		}
	}
}
