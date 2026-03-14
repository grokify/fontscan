package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/font"
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

	files, err := s.Scan()
	if err != nil {
		return fmt.Errorf("failed to scan fonts: %w", err)
	}

	// Compile filter regex if provided
	var filterRegex *regexp.Regexp
	if listFilter != "" {
		filterRegex, err = regexp.Compile("(?i)" + listFilter)
		if err != nil {
			return fmt.Errorf("invalid filter pattern: %w", err)
		}
	}

	var entries []format.FontListEntry
	for _, f := range files {
		loaded, err := font.LoadFile(f.Path)
		if err != nil {
			// Include fonts that fail to load with basic info
			entry := format.FontListEntry{
				Name: f.Name,
				Path: f.Path,
			}
			// Apply filter if set
			if filterRegex != nil {
				if !filterRegex.MatchString(f.Name) {
					continue
				}
			}
			entries = append(entries, entry)
			continue
		}

		// Apply filter if set - match against family or name
		if filterRegex != nil {
			if !filterRegex.MatchString(loaded.Family) && !filterRegex.MatchString(loaded.Name) {
				continue
			}
		}

		entries = append(entries, format.FontListEntry{
			Name:   loaded.Name,
			Path:   loaded.Path,
			Family: loaded.Family,
			Style:  loaded.Style,
		})
	}

	// Sort by name (or family for family-only mode)
	if listFamilyOnly {
		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Family) < strings.ToLower(entries[j].Family)
		})
	} else {
		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
		})
	}

	formatter := format.NewFormatter(os.Stdout, format.ParseFormat(outputFormat))

	// Handle family-only mode
	if listFamilyOnly {
		families := extractUniqueFamilies(entries)
		if listLimit > 0 && len(families) > listLimit {
			families = families[:listLimit]
		}
		return formatter.WriteFamilyList(families)
	}

	// Apply limit
	if listLimit > 0 && len(entries) > listLimit {
		entries = entries[:listLimit]
	}

	return formatter.WriteFontList(entries)
}

// extractUniqueFamilies extracts unique family names from font entries, sorted alphabetically.
func extractUniqueFamilies(entries []format.FontListEntry) []string {
	seen := make(map[string]bool)
	var families []string

	for _, e := range entries {
		family := e.Family
		if family == "" {
			family = e.Name // Fallback to name if family is empty
		}
		familyLower := strings.ToLower(family)
		if !seen[familyLower] {
			seen[familyLower] = true
			families = append(families, family)
		}
	}

	sort.Slice(families, func(i, j int) bool {
		return strings.ToLower(families[i]) < strings.ToLower(families[j])
	})

	return families
}
