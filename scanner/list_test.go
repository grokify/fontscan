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
