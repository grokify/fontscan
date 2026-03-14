package font

import (
	"sort"
)

// UnicodeBlock represents a range of Unicode code points.
type UnicodeBlock struct {
	Name  string
	Start rune
	End   rune
}

// Coverage represents the Unicode coverage of a font.
type Coverage struct {
	Font           *Font
	BlockCoverage  map[string]*BlockCoverage
	TotalSupported int
	TotalChecked   int
}

// BlockCoverage represents coverage statistics for a Unicode block.
type BlockCoverage struct {
	Block     UnicodeBlock
	Supported int
	Total     int
}

// Percentage returns the coverage percentage for a block.
func (bc *BlockCoverage) Percentage() float64 {
	if bc.Total == 0 {
		return 0
	}
	return float64(bc.Supported) / float64(bc.Total) * 100
}

// CommonBlocks defines commonly used Unicode blocks for analysis.
var CommonBlocks = []UnicodeBlock{
	{"Basic Latin", 0x0000, 0x007F},
	{"Latin-1 Supplement", 0x0080, 0x00FF},
	{"Latin Extended-A", 0x0100, 0x017F},
	{"Latin Extended-B", 0x0180, 0x024F},
	{"IPA Extensions", 0x0250, 0x02AF},
	{"Spacing Modifier Letters", 0x02B0, 0x02FF},
	{"Combining Diacritical Marks", 0x0300, 0x036F},
	{"Greek and Coptic", 0x0370, 0x03FF},
	{"Cyrillic", 0x0400, 0x04FF},
	{"Hebrew", 0x0590, 0x05FF},
	{"Arabic", 0x0600, 0x06FF},
	{"Thai", 0x0E00, 0x0E7F},
	{"General Punctuation", 0x2000, 0x206F},
	{"Superscripts and Subscripts", 0x2070, 0x209F},
	{"Currency Symbols", 0x20A0, 0x20CF},
	{"Letterlike Symbols", 0x2100, 0x214F},
	{"Number Forms", 0x2150, 0x218F},
	{"Arrows", 0x2190, 0x21FF},
	{"Mathematical Operators", 0x2200, 0x22FF},
	{"Miscellaneous Technical", 0x2300, 0x23FF},
	{"Box Drawing", 0x2500, 0x257F},
	{"Block Elements", 0x2580, 0x259F},
	{"Geometric Shapes", 0x25A0, 0x25FF},
	{"Miscellaneous Symbols", 0x2600, 0x26FF},
	{"Dingbats", 0x2700, 0x27BF},
	{"Braille Patterns", 0x2800, 0x28FF},
	{"CJK Symbols and Punctuation", 0x3000, 0x303F},
	{"Hiragana", 0x3040, 0x309F},
	{"Katakana", 0x30A0, 0x30FF},
	{"CJK Unified Ideographs", 0x4E00, 0x9FFF},
	{"Private Use Area", 0xE000, 0xF8FF},
	{"Alphabetic Presentation Forms", 0xFB00, 0xFB4F},
	{"Arabic Presentation Forms-A", 0xFB50, 0xFDFF},
	{"Halfwidth and Fullwidth Forms", 0xFF00, 0xFFEF},
	{"Emoticons", 0x1F600, 0x1F64F},
	{"Miscellaneous Symbols and Pictographs", 0x1F300, 0x1F5FF},
}

// NewCoverage creates a new Coverage analysis for the given font.
func NewCoverage(f *Font) *Coverage {
	cov := &Coverage{
		Font:          f,
		BlockCoverage: make(map[string]*BlockCoverage),
	}

	for _, block := range CommonBlocks {
		bc := &BlockCoverage{
			Block: block,
			Total: int(block.End - block.Start + 1),
		}

		for r := block.Start; r <= block.End; r++ {
			cov.TotalChecked++
			if f.HasGlyph(r) {
				bc.Supported++
				cov.TotalSupported++
			}
		}

		cov.BlockCoverage[block.Name] = bc
	}

	return cov
}

// SortedBlocks returns block coverage sorted by name.
func (c *Coverage) SortedBlocks() []*BlockCoverage {
	blocks := make([]*BlockCoverage, 0, len(c.BlockCoverage))
	for _, bc := range c.BlockCoverage {
		blocks = append(blocks, bc)
	}
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Block.Start < blocks[j].Block.Start
	})
	return blocks
}

// NonEmptyBlocks returns blocks with at least one supported glyph.
func (c *Coverage) NonEmptyBlocks() []*BlockCoverage {
	var blocks []*BlockCoverage
	for _, bc := range c.SortedBlocks() {
		if bc.Supported > 0 {
			blocks = append(blocks, bc)
		}
	}
	return blocks
}

// OverallPercentage returns the overall coverage percentage.
func (c *Coverage) OverallPercentage() float64 {
	if c.TotalChecked == 0 {
		return 0
	}
	return float64(c.TotalSupported) / float64(c.TotalChecked) * 100
}
