package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/checker"
	"github.com/grokify/fontscan/format"
	"github.com/grokify/fontscan/pandoc"
	"github.com/grokify/fontscan/scanner"
)

var (
	pandocPaths       []string
	pandocCheckAll    bool
	pandocWarnDefault bool
)

var pandocCmd = &cobra.Command{
	Use:   "pandoc <file>",
	Short: "Check fonts specified in Pandoc Markdown frontmatter",
	Long: `Parse Pandoc Markdown frontmatter and check if specified fonts support all characters.

Extracts font variables (mainfont, sansfont, monofont, etc.) from YAML frontmatter
and verifies glyph coverage for each.

Examples:
  fontscan pandoc document.md
  fontscan pandoc document.md --warn-defaults
  fontscan pandoc document.md --check-all`,
	Args: cobra.ExactArgs(1),
	RunE: runPandoc,
}

func init() {
	pandocCmd.Flags().StringSliceVar(&pandocPaths, "path", nil, "Additional font search paths")
	pandocCmd.Flags().BoolVar(&pandocCheckAll, "check-all", false,
		"Check all fonts, not just those in frontmatter")
	pandocCmd.Flags().BoolVar(&pandocWarnDefault, "warn-defaults", false,
		"Warn when using Pandoc/LaTeX default fonts (Latin Modern, etc.)")

	rootCmd.AddCommand(pandocCmd)
}

// PandocResult holds the check results for a Pandoc document.
type PandocResult struct {
	File         string                          `json:"file"`
	FontSettings map[string]string               `json:"fontSettings,omitempty"`
	Results      map[string]*checker.CheckResult `json:"results,omitempty"`
	Warnings     []string                        `json:"warnings,omitempty"`
	Errors       []string                        `json:"errors,omitempty"`
}

func runPandoc(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Parse frontmatter
	settings, err := pandoc.ParseFrontmatter(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	fontMap := settings.FontMap()
	if len(fontMap) == 0 {
		fmt.Println("No font settings found in frontmatter")
		return nil
	}

	// Set up scanner
	s := scanner.NewScanner()
	if len(pandocPaths) > 0 {
		s.Paths = append(s.Paths, pandocPaths...)
	}

	// Read file content for checking
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	text := string(data)

	result := &PandocResult{
		File:         filePath,
		FontSettings: fontMap,
		Results:      make(map[string]*checker.CheckResult),
	}

	hasErrors := false

	for varName, fontName := range fontMap {
		// Check for default fonts
		if pandocWarnDefault && pandoc.IsDefaultFont(fontName) {
			info := pandoc.GetFontInfo(fontName)
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
			hasErrors = true
			continue
		}

		checkResult := checker.CheckText(text, f)
		checkResult.File = filePath
		result.Results[varName] = checkResult

		if !checkResult.Supported {
			hasErrors = true
		}
	}

	// Output results
	if err := writePandocResult(result); err != nil {
		return err
	}

	if hasErrors {
		os.Exit(1)
	}

	return nil
}

func writePandocResult(r *PandocResult) error {
	f := format.ParseFormat(outputFormat)

	if f == format.FormatJSON || f == format.FormatJSONCompact {
		formatter := format.NewFormatter(os.Stdout, f)
		return formatter.WriteJSON(r)
	}

	// TOON format
	fmt.Printf("file: %s\n", r.File)
	fmt.Println("fontSettings:")
	for k, v := range r.FontSettings {
		fmt.Printf("  %s: %s\n", k, v)
	}

	if len(r.Warnings) > 0 {
		fmt.Println("warnings:")
		for _, w := range r.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	if len(r.Errors) > 0 {
		fmt.Println("errors:")
		for _, e := range r.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	if len(r.Results) > 0 {
		fmt.Println("results:")
		for varName, res := range r.Results {
			fmt.Printf("  %s:\n", varName)
			fmt.Printf("    font: %s\n", res.Font)
			fmt.Printf("    supported: %t\n", res.Supported)
			fmt.Printf("    missingCount: %d\n", res.MissingCount)
			if len(res.Missing) > 0 {
				fmt.Println("    missing:")
				for _, m := range res.Missing {
					fmt.Printf("      - char: %s\n", m.Char)
					fmt.Printf("        codepoint: %s\n", m.Codepoint)
					if m.Name != "" {
						fmt.Printf("        name: %s\n", m.Name)
					}
				}
			}
		}
	}

	return nil
}
