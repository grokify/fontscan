//go:build darwin

package scanner

// SystemFontPaths returns the standard font directories on macOS.
func SystemFontPaths() []string {
	return []string{
		"/System/Library/Fonts",
		"/Library/Fonts",
		"~/Library/Fonts",
	}
}
