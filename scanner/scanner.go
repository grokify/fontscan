// Package scanner provides font discovery functionality across different operating systems.
package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/fontscan/font"
)

// SupportedExtensions lists the font file extensions that can be loaded.
var SupportedExtensions = []string{".ttf", ".otf", ".ttc"}

// Scanner discovers and loads fonts from the filesystem.
type Scanner struct {
	Paths []string
}

// NewScanner creates a new Scanner with system-specific font paths.
func NewScanner() *Scanner {
	return &Scanner{
		Paths: SystemFontPaths(),
	}
}

// NewScannerWithPaths creates a Scanner with custom font paths.
func NewScannerWithPaths(paths []string) *Scanner {
	return &Scanner{
		Paths: paths,
	}
}

// FontFile represents a discovered font file.
type FontFile struct {
	Path string
	Name string
}

// Scan discovers all font files in the configured paths.
func (s *Scanner) Scan() ([]FontFile, error) {
	var fonts []FontFile

	for _, dir := range s.Paths {
		// Expand home directory
		if strings.HasPrefix(dir, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				continue
			}
			dir = filepath.Join(home, dir[2:])
		}

		files, err := s.scanDir(dir)
		if err != nil {
			// Skip directories that don't exist or can't be read
			continue
		}
		fonts = append(fonts, files...)
	}

	return fonts, nil
}

// scanDir recursively scans a directory for font files.
func (s *Scanner) scanDir(dir string) ([]FontFile, error) {
	var fonts []FontFile

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files/dirs we can't access
		}

		if d.IsDir() {
			return nil
		}

		if isFontFile(path) {
			fonts = append(fonts, FontFile{
				Path: path,
				Name: strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())),
			})
		}

		return nil
	})

	return fonts, err
}

// isFontFile checks if a file has a supported font extension.
func isFontFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, supported := range SupportedExtensions {
		if ext == supported {
			return true
		}
	}
	return false
}

// LoadAll loads all discovered fonts.
func (s *Scanner) LoadAll() ([]*font.Font, error) {
	files, err := s.Scan()
	if err != nil {
		return nil, err
	}

	var fonts []*font.Font
	for _, f := range files {
		loaded, err := font.LoadFile(f.Path)
		if err != nil {
			// Skip fonts that fail to load
			continue
		}
		fonts = append(fonts, loaded)
	}

	return fonts, nil
}

// normalize removes spaces and converts to lowercase for flexible matching.
func normalize(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}

// FontLoadError indicates a font file was found but couldn't be loaded.
type FontLoadError struct {
	Path string
	Err  error
}

func (e *FontLoadError) Error() string {
	return "font file found but failed to load: " + e.Path + ": " + e.Err.Error()
}

func (e *FontLoadError) Unwrap() error {
	return e.Err
}

// FindByName searches for a font by name (case-insensitive, space-insensitive).
func (s *Scanner) FindByName(name string) (*font.Font, error) {
	files, err := s.Scan()
	if err != nil {
		return nil, err
	}

	nameLower := strings.ToLower(name)
	nameNorm := normalize(name)

	// Track load errors for matched files
	var loadErr error
	var matchedPath string

	// First try exact match on filename (case-insensitive)
	for _, f := range files {
		if strings.ToLower(f.Name) == nameLower {
			loaded, err := font.LoadFile(f.Path)
			if err != nil {
				loadErr = err
				matchedPath = f.Path
				continue
			}
			return loaded, nil
		}
	}

	// Try normalized match on filename (ignoring spaces)
	for _, f := range files {
		if normalize(f.Name) == nameNorm {
			loaded, err := font.LoadFile(f.Path)
			if err != nil {
				loadErr = err
				matchedPath = f.Path
				continue
			}
			return loaded, nil
		}
	}

	// Then try loading fonts and matching on font name/family
	for _, f := range files {
		loaded, err := font.LoadFile(f.Path)
		if err != nil {
			continue
		}

		loadedNameNorm := normalize(loaded.Name)
		loadedFamilyNorm := normalize(loaded.Family)

		if loadedNameNorm == nameNorm || loadedFamilyNorm == nameNorm {
			return loaded, nil
		}
	}

	// Try partial match on filename
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name), nameLower) {
			loaded, err := font.LoadFile(f.Path)
			if err != nil {
				if loadErr == nil {
					loadErr = err
					matchedPath = f.Path
				}
				continue
			}
			return loaded, nil
		}
	}

	// Try partial match with normalized name
	for _, f := range files {
		if strings.Contains(normalize(f.Name), nameNorm) {
			loaded, err := font.LoadFile(f.Path)
			if err != nil {
				if loadErr == nil {
					loadErr = err
					matchedPath = f.Path
				}
				continue
			}
			return loaded, nil
		}
	}

	// If we found a matching file but couldn't load it, return that error
	if loadErr != nil {
		return nil, &FontLoadError{Path: matchedPath, Err: loadErr}
	}

	return nil, os.ErrNotExist
}

// FindByPath loads a font from a specific file path.
func FindByPath(path string) (*font.Font, error) {
	return font.LoadFile(path)
}
