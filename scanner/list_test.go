package scanner

import (
	"testing"
)

func TestExtractUniqueFamilies(t *testing.T) {
	tests := []struct {
		name     string
		entries  []FontEntry
		expected []string
	}{
		{
			name:     "empty entries",
			entries:  []FontEntry{},
			expected: nil,
		},
		{
			name: "single entry",
			entries: []FontEntry{
				{Name: "Arial Bold", Family: "Arial"},
			},
			expected: []string{"Arial"},
		},
		{
			name: "duplicate families",
			entries: []FontEntry{
				{Name: "Arial Bold", Family: "Arial"},
				{Name: "Arial Regular", Family: "Arial"},
				{Name: "Arial Italic", Family: "Arial"},
			},
			expected: []string{"Arial"},
		},
		{
			name: "multiple families",
			entries: []FontEntry{
				{Name: "Arial Bold", Family: "Arial"},
				{Name: "Helvetica Regular", Family: "Helvetica"},
				{Name: "Times New Roman", Family: "Times New Roman"},
			},
			expected: []string{"Arial", "Helvetica", "Times New Roman"},
		},
		{
			name: "case insensitive dedup",
			entries: []FontEntry{
				{Name: "Arial Bold", Family: "Arial"},
				{Name: "arial regular", Family: "arial"},
				{Name: "ARIAL Italic", Family: "ARIAL"},
			},
			expected: []string{"Arial"},
		},
		{
			name: "fallback to name when family empty",
			entries: []FontEntry{
				{Name: "MyFont", Family: ""},
				{Name: "AnotherFont", Family: ""},
			},
			expected: []string{"AnotherFont", "MyFont"},
		},
		{
			name: "sorted alphabetically",
			entries: []FontEntry{
				{Name: "Zebra", Family: "Zebra"},
				{Name: "Alpha", Family: "Alpha"},
				{Name: "Middle", Family: "Middle"},
			},
			expected: []string{"Alpha", "Middle", "Zebra"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractUniqueFamilies(tt.entries)

			if len(result) != len(tt.expected) {
				t.Errorf("ExtractUniqueFamilies() returned %d families, expected %d", len(result), len(tt.expected))
				return
			}

			for i, family := range result {
				if family != tt.expected[i] {
					t.Errorf("ExtractUniqueFamilies()[%d] = %q, expected %q", i, family, tt.expected[i])
				}
			}
		})
	}
}

func TestListOptions(t *testing.T) {
	// Test that ListOptions struct works correctly
	opts := ListOptions{
		Filter:     "arial",
		FamilyOnly: true,
		Limit:      10,
	}

	if opts.Filter != "arial" {
		t.Errorf("ListOptions.Filter = %q, expected %q", opts.Filter, "arial")
	}
	if !opts.FamilyOnly {
		t.Error("ListOptions.FamilyOnly = false, expected true")
	}
	if opts.Limit != 10 {
		t.Errorf("ListOptions.Limit = %d, expected %d", opts.Limit, 10)
	}
}

func TestScanner_List(t *testing.T) {
	// Test with system fonts (integration test)
	s := NewScanner()

	t.Run("no options", func(t *testing.T) {
		result, err := s.List(ListOptions{})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if result == nil {
			t.Fatal("List() returned nil")
		}
		// Should have entries, not families
		if len(result.Families) > 0 {
			t.Error("Expected no families when FamilyOnly is false")
		}
	})

	t.Run("family only", func(t *testing.T) {
		result, err := s.List(ListOptions{FamilyOnly: true})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if result == nil {
			t.Fatal("List() returned nil")
		}
		// Should have families, not entries
		if len(result.Entries) > 0 {
			t.Error("Expected no entries when FamilyOnly is true")
		}
	})

	t.Run("with limit", func(t *testing.T) {
		result, err := s.List(ListOptions{Limit: 5})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(result.Entries) > 5 {
			t.Errorf("List() returned %d entries, expected <= 5", len(result.Entries))
		}
	})

	t.Run("family only with limit", func(t *testing.T) {
		result, err := s.List(ListOptions{FamilyOnly: true, Limit: 3})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(result.Families) > 3 {
			t.Errorf("List() returned %d families, expected <= 3", len(result.Families))
		}
	})

	t.Run("with filter", func(t *testing.T) {
		result, err := s.List(ListOptions{Filter: "arial", FamilyOnly: true})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		// All results should match the filter
		for _, family := range result.Families {
			if !containsIgnoreCase(family, "arial") {
				t.Errorf("Family %q does not match filter 'arial'", family)
			}
		}
	})

	t.Run("invalid filter regex", func(t *testing.T) {
		_, err := s.List(ListOptions{Filter: "[invalid"})
		if err == nil {
			t.Error("Expected error for invalid regex")
		}
	})
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(substr) == 0 ||
			(len(s) > 0 && containsIgnoreCaseHelper(s, substr)))
}

func containsIgnoreCaseHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if equalIgnoreCase(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

func equalIgnoreCase(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
