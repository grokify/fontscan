//go:build windows

package scanner

import (
	"os"
	"path/filepath"
)

// SystemFontPaths returns the standard font directories on Windows.
func SystemFontPaths() []string {
	paths := []string{
		filepath.Join(os.Getenv("WINDIR"), "Fonts"),
	}

	// Add user fonts directory (Windows 10+)
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		paths = append(paths, filepath.Join(localAppData, "Microsoft", "Windows", "Fonts"))
	}

	return paths
}
