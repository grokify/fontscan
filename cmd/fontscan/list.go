package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/format"
	"github.com/grokify/fontscan/scanner"
)

var (
	listPaths      []string
	listFilter     string
	listFamilyOnly bool
	listLimit      int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available fonts",
	Long: `List all available fonts in system directories or specified paths.

Examples:
  fontscan list
  fontscan list --path /Library/Fonts
  fontscan list --format json
  fontscan list --filter "helvetica|arial" --family-only --limit 10
  fontscan list --filter "mono|courier" --family-only`,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringSliceVar(&listPaths, "path", nil, "Font directories to scan (default: system paths)")
	listCmd.Flags().StringVar(&listFilter, "filter", "", "Case-insensitive regex to filter fonts by family or name")
	listCmd.Flags().BoolVar(&listFamilyOnly, "family-only", false, "Output unique font family names only")
	listCmd.Flags().IntVar(&listLimit, "limit", 0, "Limit number of results (0 = no limit)")
}

func runList(cmd *cobra.Command, args []string) error {
	var s *scanner.Scanner

	if len(listPaths) > 0 {
		s = scanner.NewScannerWithPaths(listPaths)
	} else {
		s = scanner.NewScanner()
	}

	result, err := s.List(scanner.ListOptions{
		Filter:     listFilter,
		FamilyOnly: listFamilyOnly,
		Limit:      listLimit,
	})
	if err != nil {
		return fmt.Errorf("failed to list fonts: %w", err)
	}

	formatter := format.NewFormatter(os.Stdout, format.ParseFormat(outputFormat))
	return formatter.WriteListResult(result)
}
