// Package format provides output formatting for fontscan results.
package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/grokify/fontscan/checker"
	"github.com/grokify/fontscan/font"
	"github.com/grokify/fontscan/scanner"
)

// Format specifies the output format.
type Format string

const (
	FormatTOON        Format = "toon"
	FormatJSON        Format = "json"
	FormatJSONCompact Format = "json-compact"
)

// ParseFormat parses a format string.
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	case "json-compact":
		return FormatJSONCompact
	default:
		return FormatTOON
	}
}

// Formatter handles output formatting.
type Formatter struct {
	Format Format
	Writer io.Writer
}

// NewFormatter creates a new Formatter.
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{
		Format: format,
		Writer: w,
	}
}

// WriteCheckResult writes a CheckResult in the configured format.
func (f *Formatter) WriteCheckResult(r *checker.CheckResult) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(r, true)
	case FormatJSONCompact:
		return f.writeJSON(r, false)
	default:
		return f.writeCheckResultTOON(r)
	}
}

// WriteRecommendResult writes a RecommendResult in the configured format.
func (f *Formatter) WriteRecommendResult(r *checker.RecommendResult) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(r, true)
	case FormatJSONCompact:
		return f.writeJSON(r, false)
	default:
		return f.writeRecommendResultTOON(r)
	}
}

// WriteFontInfo writes font information in the configured format.
func (f *Formatter) WriteFontInfo(fnt *font.Font, showCoverage bool) error {
	info := struct {
		Name      string              `json:"name"`
		Family    string              `json:"family,omitempty"`
		Style     string              `json:"style,omitempty"`
		Path      string              `json:"path"`
		NumGlyphs int                 `json:"numGlyphs"`
		Coverage  []BlockCoverageInfo `json:"coverage,omitempty"`
	}{
		Name:      fnt.Name,
		Family:    fnt.Family,
		Style:     fnt.Style,
		Path:      fnt.Path,
		NumGlyphs: fnt.NumGlyphs,
	}

	if showCoverage {
		cov := fnt.Coverage()
		for _, bc := range cov.NonEmptyBlocks() {
			info.Coverage = append(info.Coverage, BlockCoverageInfo{
				Block:      bc.Block.Name,
				Supported:  bc.Supported,
				Total:      bc.Total,
				Percentage: bc.Percentage(),
			})
		}
	}

	switch f.Format {
	case FormatJSON:
		return f.writeJSON(info, true)
	case FormatJSONCompact:
		return f.writeJSON(info, false)
	default:
		return f.writeFontInfoTOON(fnt, showCoverage)
	}
}

// BlockCoverageInfo represents coverage info for JSON output.
type BlockCoverageInfo struct {
	Block      string  `json:"block"`
	Supported  int     `json:"supported"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

func (f *Formatter) writeJSON(v any, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	_, err = f.Writer.Write(data)
	if err != nil {
		return err
	}

	_, err = f.Writer.Write([]byte("\n"))
	return err
}

// WriteJSON writes any value as JSON in the configured format.
func (f *Formatter) WriteJSON(v any) error {
	return f.writeJSON(v, f.Format == FormatJSON)
}

func (f *Formatter) writeCheckResultTOON(r *checker.CheckResult) error {
	var sb strings.Builder

	if r.File != "" {
		fmt.Fprintf(&sb, "file: %s\n", r.File)
	}
	fmt.Fprintf(&sb, "font: %s\n", r.Font)
	fmt.Fprintf(&sb, "supported: %t\n", r.Supported)
	fmt.Fprintf(&sb, "totalChars: %d\n", r.TotalChars)
	fmt.Fprintf(&sb, "uniqueChars: %d\n", r.UniqueChars)
	fmt.Fprintf(&sb, "missingCount: %d\n", r.MissingCount)

	if len(r.Missing) > 0 {
		sb.WriteString("missing:\n")
		for _, m := range r.Missing {
			fmt.Fprintf(&sb, "  - char: %s\n", m.Char)
			fmt.Fprintf(&sb, "    codepoint: %s\n", m.Codepoint)
			if m.Name != "" {
				fmt.Fprintf(&sb, "    name: %s\n", m.Name)
			}
			fmt.Fprintf(&sb, "    count: %d\n", m.Count)
		}
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

func (f *Formatter) writeRecommendResultTOON(r *checker.RecommendResult) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "file: %s\n", r.File)
	fmt.Fprintf(&sb, "totalChars: %d\n", r.TotalChars)
	fmt.Fprintf(&sb, "uniqueChars: %d\n", r.UniqueChars)
	fmt.Fprintf(&sb, "fullSupportCount: %d\n", len(r.FullSupport))

	if len(r.FullSupport) > 0 {
		sb.WriteString("fullSupport:\n")
		for _, m := range r.FullSupport {
			fmt.Fprintf(&sb, "  - name: %s\n", m.Name)
			fmt.Fprintf(&sb, "    path: %s\n", m.Path)
		}
	}

	if len(r.PartialSupport) > 0 {
		sb.WriteString("partialSupport:\n")
		for _, m := range r.PartialSupport {
			fmt.Fprintf(&sb, "  - name: %s\n", m.Name)
			fmt.Fprintf(&sb, "    percentage: %.1f%%\n", m.Percentage)
			fmt.Fprintf(&sb, "    missingCount: %d\n", m.MissingCount)
		}
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

func (f *Formatter) writeFontInfoTOON(fnt *font.Font, showCoverage bool) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "name: %s\n", fnt.Name)
	if fnt.Family != "" {
		fmt.Fprintf(&sb, "family: %s\n", fnt.Family)
	}
	if fnt.Style != "" {
		fmt.Fprintf(&sb, "style: %s\n", fnt.Style)
	}
	fmt.Fprintf(&sb, "path: %s\n", fnt.Path)
	fmt.Fprintf(&sb, "numGlyphs: %d\n", fnt.NumGlyphs)

	if showCoverage {
		cov := fnt.Coverage()
		sb.WriteString("coverage:\n")
		for _, bc := range cov.NonEmptyBlocks() {
			fmt.Fprintf(&sb, "  - block: %s\n", bc.Block.Name)
			fmt.Fprintf(&sb, "    supported: %d/%d (%.1f%%)\n",
				bc.Supported, bc.Total, bc.Percentage())
		}
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

// WriteFontList writes a list of fonts in the configured format.
func (f *Formatter) WriteFontList(fonts []FontListEntry) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(fonts, true)
	case FormatJSONCompact:
		return f.writeJSON(fonts, false)
	default:
		return f.writeFontListTOON(fonts)
	}
}

// FontListEntry represents a font in a list.
type FontListEntry struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Family string `json:"family,omitempty"`
	Style  string `json:"style,omitempty"`
}

func (f *Formatter) writeFontListTOON(fonts []FontListEntry) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "count: %d\n", len(fonts))
	sb.WriteString("fonts:\n")

	for _, fnt := range fonts {
		fmt.Fprintf(&sb, "  - name: %s\n", fnt.Name)
		fmt.Fprintf(&sb, "    path: %s\n", fnt.Path)
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

// WriteFamilyList writes a list of font family names in the configured format.
func (f *Formatter) WriteFamilyList(families []string) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(families, true)
	case FormatJSONCompact:
		return f.writeJSON(families, false)
	default:
		return f.writeFamilyListTOON(families)
	}
}

func (f *Formatter) writeFamilyListTOON(families []string) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "count: %d\n", len(families))
	sb.WriteString("families:\n")

	for _, family := range families {
		fmt.Fprintf(&sb, "  - %s\n", family)
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}

// WriteListResult writes a scanner.ListResult in the configured format.
func (f *Formatter) WriteListResult(r *scanner.ListResult) error {
	if len(r.Families) > 0 {
		return f.WriteFamilyList(r.Families)
	}

	// Convert scanner.FontEntry to FontListEntry
	entries := make([]FontListEntry, len(r.Entries))
	for i, e := range r.Entries {
		entries[i] = FontListEntry{
			Name:   e.Name,
			Path:   e.Path,
			Family: e.Family,
			Style:  e.Style,
		}
	}
	return f.WriteFontList(entries)
}

// PandocCheckResult mirrors pandoc.CheckResult to avoid import cycles.
// Use this when formatting pandoc check results.
type PandocCheckResult struct {
	File         string                          `json:"file"`
	FontSettings map[string]string               `json:"fontSettings,omitempty"`
	Results      map[string]*checker.CheckResult `json:"results,omitempty"`
	Warnings     []string                        `json:"warnings,omitempty"`
	Errors       []string                        `json:"errors,omitempty"`
}

// WritePandocResult writes a PandocCheckResult in the configured format.
func (f *Formatter) WritePandocResult(r *PandocCheckResult) error {
	switch f.Format {
	case FormatJSON, FormatJSONCompact:
		return f.WriteJSON(r)
	default:
		return f.writePandocResultTOON(r)
	}
}

func (f *Formatter) writePandocResultTOON(r *PandocCheckResult) error {
	var sb strings.Builder

	fmt.Fprintf(&sb, "file: %s\n", r.File)
	sb.WriteString("fontSettings:\n")
	for k, v := range r.FontSettings {
		fmt.Fprintf(&sb, "  %s: %s\n", k, v)
	}

	if len(r.Warnings) > 0 {
		sb.WriteString("warnings:\n")
		for _, w := range r.Warnings {
			fmt.Fprintf(&sb, "  - %s\n", w)
		}
	}

	if len(r.Errors) > 0 {
		sb.WriteString("errors:\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "  - %s\n", e)
		}
	}

	if len(r.Results) > 0 {
		sb.WriteString("results:\n")
		for varName, res := range r.Results {
			fmt.Fprintf(&sb, "  %s:\n", varName)
			fmt.Fprintf(&sb, "    font: %s\n", res.Font)
			fmt.Fprintf(&sb, "    supported: %t\n", res.Supported)
			fmt.Fprintf(&sb, "    missingCount: %d\n", res.MissingCount)
			if len(res.Missing) > 0 {
				sb.WriteString("    missing:\n")
				for _, m := range res.Missing {
					fmt.Fprintf(&sb, "      - char: %s\n", m.Char)
					fmt.Fprintf(&sb, "        codepoint: %s\n", m.Codepoint)
					if m.Name != "" {
						fmt.Fprintf(&sb, "        name: %s\n", m.Name)
					}
				}
			}
		}
	}

	_, err := f.Writer.Write([]byte(sb.String()))
	return err
}
