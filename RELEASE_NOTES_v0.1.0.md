# Release Notes - v0.1.0

Initial release of fontscan, a Go library and CLI for checking font glyph coverage.

## Features

### CLI Commands

- **`fontscan check`** - Check if a font supports all characters in a document
  - `--font` - Specify font by name (searches system fonts)
  - `--font-path` - Specify font by file path
  - `--path` - Additional font search directories
  - `--format` - Output format (toon, json, json-compact)

- **`fontscan pandoc`** - Parse Pandoc Markdown frontmatter and check fonts
  - Extracts font variables from YAML frontmatter (mainfont, sansfont, monofont, etc.)
  - Checks each font against document content
  - `--warn-defaults` - Warn about Latin Modern/Computer Modern default fonts

- **`fontscan list`** - List available system fonts
  - `--filter` - Case-insensitive regex to filter by family or name
  - `--family-only` - Output unique font family names only
  - `--limit` - Limit number of results
  - `--path` - Scan specific directories
  - Replaces common `fc-list | grep | sort -u | head` workflows

- **`fontscan info`** - Show font information
  - `--coverage` - Display Unicode block coverage statistics

- **`fontscan recommend`** - Find fonts that support all characters in a document
  - `--limit` - Limit number of partial-match results

### Library Packages

- **`font`** - Font loading and glyph checking
  - Supports TTF, OTF, and TTC (TrueType Collection) formats
  - Unicode coverage analysis by block

- **`scanner`** - System font discovery
  - Cross-platform support (macOS, Linux, Windows)
  - Flexible name matching (case-insensitive, space-insensitive)

- **`checker`** - Document-to-font compatibility checking
  - Returns detailed missing character information with Unicode names

- **`pandoc`** - Pandoc Markdown support
  - YAML frontmatter parsing
  - Default font database (Latin Modern, Computer Modern, TeX Gyre, DejaVu, Libertinus)

- **`format`** - Output formatting
  - TOON format (default, human-readable)
  - JSON and JSON-compact formats

### Key Implementation Details

- **TTC Collection Handling** - When loading a TrueType Collection, fontscan iterates through all fonts until one loads successfully, avoiding failures when the first font has unsupported cmap encodings.

- **Unicode Character Names** - Missing characters are reported with their Unicode name (e.g., "RIGHTWARDS ARROW" for U+2192).

- **Exit Codes** - Returns exit code 1 when missing characters are found, useful for CI/CD pipelines.

## Supported Platforms

- macOS (tested)
- Linux
- Windows

## Dependencies

- `golang.org/x/image` - TTF/OTF parsing via `font/sfnt`
- `golang.org/x/text` - Unicode character names
- `github.com/spf13/cobra` - CLI framework
- `gopkg.in/yaml.v3` - YAML frontmatter parsing

## Known Limitations

- Some fonts use cmap encodings not supported by Go's `sfnt` library
- Coverage analysis checks common Unicode blocks, not all possible code points

## Installation

```bash
go install github.com/grokify/fontscan/cmd/fontscan@latest
```

## What's Next

Planned for future releases:

- Custom font variable support in Pandoc frontmatter
- Font fallback chain analysis
- Integration with Pandoc defaults files
- Coverage report export (HTML, CSV)
