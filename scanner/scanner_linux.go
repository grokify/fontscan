//go:build linux

package scanner

// SystemFontPaths returns the standard font directories on Linux.
func SystemFontPaths() []string {
	return []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		"~/.fonts",
		"~/.local/share/fonts",
	}
}
