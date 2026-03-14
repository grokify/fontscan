package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/format"
	"github.com/grokify/fontscan/pandoc"
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

func runPandoc(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	result, err := pandoc.CheckDocument(filePath, pandoc.CheckOptions{
		Paths:        pandocPaths,
		WarnDefaults: pandocWarnDefault,
	})
	if err != nil {
		return err
	}

	if !result.HasFontSettings() {
		fmt.Println("No font settings found in frontmatter")
		return nil
	}

	// Convert to format.PandocCheckResult
	formatResult := &format.PandocCheckResult{
		File:         result.File,
		FontSettings: result.FontSettings,
		Results:      result.Results,
		Warnings:     result.Warnings,
		Errors:       result.Errors,
	}

	formatter := format.NewFormatter(os.Stdout, format.ParseFormat(outputFormat))
	if err := formatter.WritePandocResult(formatResult); err != nil {
		return err
	}

	if result.HasErrors {
		os.Exit(1)
	}

	return nil
}
