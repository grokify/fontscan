package checker

import (
	"os"
	"sort"

	"github.com/grokify/fontscan/font"
	"github.com/grokify/fontscan/scanner"
)

// Checker performs font compatibility checks on text.
type Checker struct {
	Scanner *scanner.Scanner
}

// NewChecker creates a new Checker with system font paths.
func NewChecker() *Checker {
	return &Checker{
		Scanner: scanner.NewScanner(),
	}
}

// NewCheckerWithPaths creates a Checker with custom font paths.
func NewCheckerWithPaths(paths []string) *Checker {
	return &Checker{
		Scanner: scanner.NewScannerWithPaths(paths),
	}
}

// CheckText checks if all characters in the text are supported by the font.
func CheckText(text string, f *font.Font) *CheckResult {
	charCounts := make(map[rune]int)
	totalChars := 0

	for _, r := range text {
		if !IsSignificant(r) {
			continue
		}
		charCounts[r]++
		totalChars++
	}

	var missing []MissingChar
	for r, count := range charCounts {
		if !f.HasGlyph(r) {
			missing = append(missing, NewMissingChar(r, count))
		}
	}

	// Sort missing characters by codepoint for consistent output
	sort.Slice(missing, func(i, j int) bool {
		return missing[i].Codepoint < missing[j].Codepoint
	})

	return &CheckResult{
		Font:         f.Name,
		FontPath:     f.Path,
		Supported:    len(missing) == 0,
		TotalChars:   totalChars,
		UniqueChars:  len(charCounts),
		MissingCount: len(missing),
		Missing:      missing,
	}
}

// CheckFile checks if all characters in a file are supported by the font.
func CheckFile(path string, f *font.Font) (*CheckResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := CheckText(string(data), f)
	result.File = path
	return result, nil
}

// Check checks a file against a font specified by name.
func (c *Checker) Check(filePath, fontName string) (*CheckResult, error) {
	f, err := c.Scanner.FindByName(fontName)
	if err != nil {
		return nil, err
	}

	return CheckFile(filePath, f)
}

// CheckWithFontPath checks a file against a font at a specific path.
func (c *Checker) CheckWithFontPath(filePath, fontPath string) (*CheckResult, error) {
	f, err := font.LoadFile(fontPath)
	if err != nil {
		return nil, err
	}

	return CheckFile(filePath, f)
}

// Recommend finds fonts that support all characters in a file.
func (c *Checker) Recommend(filePath string, limit int) (*RecommendResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	text := string(data)
	charCounts := make(map[rune]int)
	totalChars := 0

	for _, r := range text {
		if !IsSignificant(r) {
			continue
		}
		charCounts[r]++
		totalChars++
	}

	uniqueChars := len(charCounts)

	result := &RecommendResult{
		File:        filePath,
		TotalChars:  totalChars,
		UniqueChars: uniqueChars,
	}

	fonts, err := c.Scanner.LoadAll()
	if err != nil {
		return nil, err
	}

	var partials []PartialMatch

	for _, f := range fonts {
		supported := 0
		for r := range charCounts {
			if f.HasGlyph(r) {
				supported++
			}
		}

		if supported == uniqueChars {
			result.FullSupport = append(result.FullSupport, FontMatch{
				Name:   f.Name,
				Path:   f.Path,
				Family: f.Family,
				Style:  f.Style,
			})
		} else if supported > 0 {
			percentage := float64(supported) / float64(uniqueChars) * 100
			partials = append(partials, PartialMatch{
				FontMatch: FontMatch{
					Name:   f.Name,
					Path:   f.Path,
					Family: f.Family,
					Style:  f.Style,
				},
				Supported:    supported,
				Total:        uniqueChars,
				Percentage:   percentage,
				MissingCount: uniqueChars - supported,
			})
		}
	}

	// Sort partial matches by support percentage descending
	sort.Slice(partials, func(i, j int) bool {
		return partials[i].Percentage > partials[j].Percentage
	})

	// Limit partial matches if requested
	if limit > 0 && len(partials) > limit {
		partials = partials[:limit]
	}

	result.PartialSupport = partials

	return result, nil
}
