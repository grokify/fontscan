package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/font"
	"github.com/grokify/fontscan/format"
	"github.com/grokify/fontscan/scanner"
)

var (
	infoCoverage  bool
	infoFontPath  string
	infoFontPaths []string
)

var infoCmd = &cobra.Command{
	Use:   "info <font>",
	Short: "Show font information",
	Long: `Display information about a font, optionally including Unicode block coverage.

The font can be specified by name (searches system fonts) or by path.

Examples:
  fontscan info "DejaVu Sans Mono"
  fontscan info "Helvetica" --coverage
  fontscan info --font-path /Library/Fonts/Arial.ttf
  fontscan info "Courier" --format json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInfo,
}

func init() {
	infoCmd.Flags().BoolVar(&infoCoverage, "coverage", false, "Show Unicode block coverage")
	infoCmd.Flags().StringVar(&infoFontPath, "font-path", "", "Path to font file")
	infoCmd.Flags().StringSliceVar(&infoFontPaths, "path", nil, "Additional font search paths")
}

func runInfo(cmd *cobra.Command, args []string) error {
	var f *font.Font
	var err error

	if infoFontPath != "" {
		// Load font from specific path
		f, err = font.LoadFile(infoFontPath)
		if err != nil {
			return fmt.Errorf("failed to load font: %w", err)
		}
	} else if len(args) == 1 {
		// Search for font by name
		s := scanner.NewScanner()
		if len(infoFontPaths) > 0 {
			s.Paths = append(s.Paths, infoFontPaths...)
		}

		f, err = s.FindByName(args[0])
		if err != nil {
			return fmt.Errorf("font not found: %s", args[0])
		}
	} else {
		return fmt.Errorf("either font name or --font-path must be specified")
	}

	formatter := format.NewFormatter(os.Stdout, format.ParseFormat(outputFormat))
	return formatter.WriteFontInfo(f, infoCoverage)
}
