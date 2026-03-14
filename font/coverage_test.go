package font

import (
	"testing"
)

func TestNewCoverage(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	cov := NewCoverage(f)

	if cov.Font != f {
		t.Error("Coverage.Font should point to the font")
	}

	if len(cov.BlockCoverage) == 0 {
		t.Error("Expected some block coverage data")
	}

	if cov.TotalChecked == 0 {
		t.Error("Expected TotalChecked > 0")
	}
}

func TestCoverage_SortedBlocks(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	cov := NewCoverage(f)
	blocks := cov.SortedBlocks()

	if len(blocks) == 0 {
		t.Error("Expected some blocks")
	}

	// Verify sorted by start codepoint
	for i := 1; i < len(blocks); i++ {
		if blocks[i].Block.Start < blocks[i-1].Block.Start {
			t.Errorf("Blocks not sorted: %s (U+%04X) before %s (U+%04X)",
				blocks[i-1].Block.Name, blocks[i-1].Block.Start,
				blocks[i].Block.Name, blocks[i].Block.Start)
		}
	}
}

func TestCoverage_NonEmptyBlocks(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	cov := NewCoverage(f)
	blocks := cov.NonEmptyBlocks()

	// All blocks should have at least one supported glyph
	for _, bc := range blocks {
		if bc.Supported == 0 {
			t.Errorf("Block %s has 0 supported glyphs but was included in NonEmptyBlocks",
				bc.Block.Name)
		}
	}
}

func TestBlockCoverage_Percentage(t *testing.T) {
	tests := []struct {
		supported int
		total     int
		expected  float64
	}{
		{0, 100, 0},
		{50, 100, 50},
		{100, 100, 100},
		{0, 0, 0}, // Edge case: avoid division by zero
	}

	for _, tt := range tests {
		bc := &BlockCoverage{
			Supported: tt.supported,
			Total:     tt.total,
		}

		got := bc.Percentage()
		if got != tt.expected {
			t.Errorf("Percentage(%d/%d) = %f, want %f",
				tt.supported, tt.total, got, tt.expected)
		}
	}
}

func TestCoverage_OverallPercentage(t *testing.T) {
	path := findTestFont(t)

	f, err := LoadFile(path)
	if err != nil {
		t.Fatalf("Failed to load font: %v", err)
	}

	cov := NewCoverage(f)
	pct := cov.OverallPercentage()

	if pct < 0 || pct > 100 {
		t.Errorf("OverallPercentage() = %f, want 0-100", pct)
	}
}

func TestCommonBlocks(t *testing.T) {
	if len(CommonBlocks) == 0 {
		t.Error("Expected some CommonBlocks defined")
	}

	// Verify no overlapping ranges
	for i, b1 := range CommonBlocks {
		for j, b2 := range CommonBlocks {
			if i >= j {
				continue
			}
			// Check if ranges overlap
			if b1.Start <= b2.End && b2.Start <= b1.End {
				t.Errorf("Blocks %s and %s overlap: [%04X-%04X] and [%04X-%04X]",
					b1.Name, b2.Name, b1.Start, b1.End, b2.Start, b2.End)
			}
		}
	}
}
