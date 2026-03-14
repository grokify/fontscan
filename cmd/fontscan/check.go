package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/checker"
	"github.com/grokify/fontscan/font"
	"github.com/grokify/fontscan/format"
	"github.com/grokify/fontscan/scanner"
)

var (
	checkFontName string
	checkFontPath string
	checkPaths    []string
)

var checkCmd = &cobra.Command{
	Use:   "check <file>",
	Short: "Check if a font supports all characters in a file",
	Long: `Check if all characters in a file are supported by the specified font.

Specify a font by name (searches system fonts) or by path.

Examples:
  fontscan check document.md --font "Helvetica Neue"
  fontscan check document.md --font-path /Library/Fonts/Arial.ttf
  fontscan check document.md --font "DejaVu Sans" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runCheck,
}

func init() {
	checkCmd.Flags().StringVar(&checkFontName, "font", "", "Font name to check against")
	checkCmd.Flags().StringVar(&checkFontPath, "font-path", "", "Path to font file")
	checkCmd.Flags().StringSliceVar(&checkPaths, "path", nil, "Additional font search paths")
}

func runCheck(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Validate that either font name or path is provided
	if checkFontName == "" && checkFontPath == "" {
		return fmt.Errorf("either --font or --font-path must be specified")
	}

	var f *font.Font
	var err error

	if checkFontPath != "" {
		// Load font from specific path
		f, err = font.LoadFile(checkFontPath)
		if err != nil {
			return fmt.Errorf("failed to load font: %w", err)
		}
	} else {
		// Search for font by name
		s := scanner.NewScanner()
		if len(checkPaths) > 0 {
			s.Paths = append(s.Paths, checkPaths...)
		}

		f, err = s.FindByName(checkFontName)
		if err != nil {
			// Check if it's a load error (font found but couldn't be loaded)
			var loadErr *scanner.FontLoadError
			if errors.As(err, &loadErr) {
				return fmt.Errorf("font '%s' found at %s but failed to load: %v",
					checkFontName, loadErr.Path, loadErr.Err)
			}
			return fmt.Errorf("font not found: %s", checkFontName)
		}
	}

	result, err := checker.CheckFile(filePath, f)
	if err != nil {
		return fmt.Errorf("failed to check file: %w", err)
	}

	formatter := format.NewFormatter(os.Stdout, format.ParseFormat(outputFormat))
	if err := formatter.WriteCheckResult(result); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Exit with non-zero code if there are missing characters
	if !result.Supported {
		os.Exit(1)
	}

	return nil
}
