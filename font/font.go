// Package font provides font loading and glyph coverage analysis.
package font

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font/sfnt"
)

// Font represents a loaded font file with glyph coverage information.
type Font struct {
	Path      string
	Name      string
	Family    string
	Style     string
	NumGlyphs int
	sfntFont  *sfnt.Font
}

// LoadFile loads a font from the given file path.
// For font collections (.ttc), it loads the first font in the collection.
// If the first font fails to load, it tries subsequent fonts.
func LoadFile(path string) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Try as collection first, iterate fonts until one works
	collection, collErr := sfnt.ParseCollection(data)
	if collErr == nil {
		var lastErr error
		for i := 0; i < collection.NumFonts(); i++ {
			f, err := collection.Font(i)
			if err != nil {
				lastErr = err
				continue
			}
			font := &Font{
				Path:      path,
				sfntFont:  f,
				NumGlyphs: f.NumGlyphs(),
			}
			font.Name = font.extractName(sfnt.NameIDFull)
			if font.Name == "" {
				font.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			}
			font.Family = font.extractName(sfnt.NameIDFamily)
			font.Style = font.extractName(sfnt.NameIDSubfamily)
			return font, nil
		}
		if lastErr != nil {
			return nil, lastErr
		}
	}

	// Try as single font
	return LoadFileIndex(path, 0)
}

// LoadFileIndex loads a specific font from a file by index.
// For single font files, index should be 0.
// For font collections, index selects which font to load.
func LoadFileIndex(path string, index int) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f *sfnt.Font

	// Try parsing as a font collection first
	collection, err := sfnt.ParseCollection(data)
	if err == nil {
		// It's a collection
		if index >= collection.NumFonts() {
			index = 0
		}
		f, err = collection.Font(index)
		if err != nil {
			return nil, err
		}
	} else {
		// Try parsing as a single font
		f, err = sfnt.Parse(data)
		if err != nil {
			return nil, err
		}
	}

	font := &Font{
		Path:      path,
		sfntFont:  f,
		NumGlyphs: f.NumGlyphs(),
	}

	// Extract font name from metadata
	font.Name = font.extractName(sfnt.NameIDFull)
	if font.Name == "" {
		// Fallback to filename without extension
		font.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	font.Family = font.extractName(sfnt.NameIDFamily)
	font.Style = font.extractName(sfnt.NameIDSubfamily)

	return font, nil
}

// LoadCollection loads all fonts from a font collection file (.ttc).
func LoadCollection(path string) ([]*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	collection, err := sfnt.ParseCollection(data)
	if err != nil {
		// Not a collection, try as single font
		f, err := LoadFile(path)
		if err != nil {
			return nil, err
		}
		return []*Font{f}, nil
	}

	var fonts []*Font
	for i := 0; i < collection.NumFonts(); i++ {
		f, err := LoadFileIndex(path, i)
		if err != nil {
			continue
		}
		fonts = append(fonts, f)
	}

	return fonts, nil
}

// extractName extracts a name string from the font's name table.
func (f *Font) extractName(nameID sfnt.NameID) string {
	var buf sfnt.Buffer
	name, err := f.sfntFont.Name(&buf, nameID)
	if err != nil {
		return ""
	}
	return name
}

// HasGlyph checks if the font contains a glyph for the given rune.
func (f *Font) HasGlyph(r rune) bool {
	var buf sfnt.Buffer
	idx, err := f.sfntFont.GlyphIndex(&buf, r)
	if err != nil {
		return false
	}
	// GlyphIndex returns 0 for the .notdef glyph (missing)
	return idx != 0
}

// MissingGlyphs returns a slice of runes from the text that are not supported by this font.
func (f *Font) MissingGlyphs(text string) []rune {
	seen := make(map[rune]bool)
	var missing []rune

	for _, r := range text {
		if seen[r] {
			continue
		}
		seen[r] = true

		if !f.HasGlyph(r) {
			missing = append(missing, r)
		}
	}

	return missing
}

// SupportedGlyphs returns a slice of runes from the text that are supported by this font.
func (f *Font) SupportedGlyphs(text string) []rune {
	seen := make(map[rune]bool)
	var supported []rune

	for _, r := range text {
		if seen[r] {
			continue
		}
		seen[r] = true

		if f.HasGlyph(r) {
			supported = append(supported, r)
		}
	}

	return supported
}

// Coverage returns the Unicode coverage information for this font.
func (f *Font) Coverage() *Coverage {
	return NewCoverage(f)
}

// String returns a string representation of the font.
func (f *Font) String() string {
	if f.Name != "" {
		return f.Name
	}
	return f.Path
}
