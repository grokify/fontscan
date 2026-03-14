package pandoc

import (
	"errors"
	"fmt"
	"os"

	"github.com/grokify/fontscan/checker"
	"github.com/grokify/fontscan/scanner"
)

// CheckOptions configures document checking behavior.
type CheckOptions struct {
	Paths        []string // Additional font search paths
	WarnDefaults bool     // Warn when using Pandoc/LaTeX default fonts
}

// CheckResult holds the check results for a Pandoc document.
type CheckResult struct {
	File         string                          `json:"file"`
	FontSettings map[string]string               `json:"fontSettings,omitempty"`
	Results      map[string]*checker.CheckResult `json:"results,omitempty"`
	Warnings     []string                        `json:"warnings,omitempty"`
	Errors       []string                        `json:"errors,omitempty"`
	HasErrors    bool                            `json:"-"`
}

// CheckDocument parses frontmatter and checks all specified fonts against the document.
func CheckDocument(filePath string, opts CheckOptions) (*CheckResult, error) {
	// Parse frontmatter
	settings, err := ParseFrontmatter(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	fontMap := settings.FontMap()

	result := &CheckResult{
		File:         filePath,
		FontSettings: fontMap,
		Results:      make(map[string]*checker.CheckResult),
	}

	if len(fontMap) == 0 {
		return result, nil
	}

	// Set up scanner
	s := scanner.NewScanner()
	if len(opts.Paths) > 0 {
		s.Paths = append(s.Paths, opts.Paths...)
	}

	// Read file content for checking
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	text := string(data)

	for varName, fontName := range fontMap {
		// Check for default fonts
		if opts.WarnDefaults && IsDefaultFont(fontName) {
			info := GetFontInfo(fontName)
			if info != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("%s: using default font '%s' (%s) - may have limited Unicode coverage",
						varName, fontName, info.Name))
			}
		}

		// Find and check the font
		f, err := s.FindByName(fontName)
		if err != nil {
			var loadErr *scanner.FontLoadError
			if errors.As(err, &loadErr) {
				result.Errors = append(result.Errors,
					fmt.Sprintf("%s: font '%s' found at %s but failed to load: %v",
						varName, fontName, loadErr.Path, loadErr.Err))
			} else {
				result.Errors = append(result.Errors,
					fmt.Sprintf("%s: font '%s' not found", varName, fontName))
			}
			result.HasErrors = true
			continue
		}

		checkResult := checker.CheckText(text, f)
		checkResult.File = filePath
		result.Results[varName] = checkResult

		if !checkResult.Supported {
			result.HasErrors = true
		}
	}

	return result, nil
}

// HasFontSettings returns true if the result contains any font settings.
func (r *CheckResult) HasFontSettings() bool {
	return len(r.FontSettings) > 0
}
