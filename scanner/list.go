package scanner

import (
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/fontscan/font"
)

// ListOptions configures font listing behavior.
type ListOptions struct {
	Filter     string // Case-insensitive regex pattern to filter by family or name
	FamilyOnly bool   // Return unique family names only
	Limit      int    // Maximum number of results (0 = no limit)
}

// FontEntry represents a font in a listing.
type FontEntry struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Family string `json:"family,omitempty"`
	Style  string `json:"style,omitempty"`
}

// ListResult contains the results of a font listing operation.
type ListResult struct {
	Entries  []FontEntry // Font entries (when FamilyOnly is false)
	Families []string    // Family names (when FamilyOnly is true)
}

// List returns fonts matching the given options.
func (s *Scanner) List(opts ListOptions) (*ListResult, error) {
	files, err := s.Scan()
	if err != nil {
		return nil, err
	}

	// Compile filter regex if provided
	var filterRegex *regexp.Regexp
	if opts.Filter != "" {
		filterRegex, err = regexp.Compile("(?i)" + opts.Filter)
		if err != nil {
			return nil, err
		}
	}

	var entries []FontEntry
	for _, f := range files {
		loaded, err := font.LoadFile(f.Path)
		if err != nil {
			// Include fonts that fail to load with basic info
			entry := FontEntry{
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

		entries = append(entries, FontEntry{
			Name:   loaded.Name,
			Path:   loaded.Path,
			Family: loaded.Family,
			Style:  loaded.Style,
		})
	}

	// Sort entries
	if opts.FamilyOnly {
		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Family) < strings.ToLower(entries[j].Family)
		})
	} else {
		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
		})
	}

	result := &ListResult{}

	if opts.FamilyOnly {
		result.Families = ExtractUniqueFamilies(entries)
		if opts.Limit > 0 && len(result.Families) > opts.Limit {
			result.Families = result.Families[:opts.Limit]
		}
	} else {
		result.Entries = entries
		if opts.Limit > 0 && len(result.Entries) > opts.Limit {
			result.Entries = result.Entries[:opts.Limit]
		}
	}

	return result, nil
}

// ExtractUniqueFamilies extracts unique family names from font entries, sorted alphabetically.
func ExtractUniqueFamilies(entries []FontEntry) []string {
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
