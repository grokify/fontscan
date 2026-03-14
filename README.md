# fontscan

[![Go Reference](https://pkg.go.dev/badge/github.com/grokify/fontscan.svg)](https://pkg.go.dev/github.com/grokify/fontscan)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/fontscan)](https://goreportcard.com/report/github.com/grokify/fontscan)

A Go library and CLI for checking font glyph coverage to prevent Pandoc/LaTeX missing character errors.

## Overview

When generating PDFs with Pandoc and LaTeX, missing font glyphs result in errors or blank characters. `fontscan` helps you:

- **Check document compliance** - Verify all characters in a document are supported by specified fonts
- **Find missing characters** - Identify which characters are not supported by a font
- **Recommend fonts** - Find installed fonts that support all characters in a document
- **Parse Pandoc frontmatter** - Extract and check fonts specified in YAML frontmatter

## Installation

```bash
go install github.com/grokify/fontscan/cmd/fontscan@latest
```

## CLI Usage

### Check a document against a font

```bash
# Check if a font supports all characters in a file
fontscan check document.md --font "Helvetica Neue"

# Use a specific font file
fontscan check document.md --font-path /Library/Fonts/Arial.ttf

# JSON output
fontscan check document.md --font "Arial" --format json
```

Output:

```
file: document.md
font: Helvetica Neue Italic
supported: false
totalChars: 157
uniqueChars: 41
missingCount: 1
missing:
  - char: →
    codepoint: U+2192
    name: RIGHTWARDS ARROW
    count: 1
```

### Check Pandoc Markdown frontmatter

```bash
# Parse YAML frontmatter and check all specified fonts
fontscan pandoc document.md

# Warn about default LaTeX fonts with limited Unicode coverage
fontscan pandoc document.md --warn-defaults
```

Example document:

```markdown
---
title: My Document
mainfont: "Helvetica Neue"
monofont: "Menlo"
sansfont: "Arial"
---

# Content here...
```

Output:

```
file: document.md
fontSettings:
  mainfont: Helvetica Neue
  sansfont: Arial
  monofont: Menlo
results:
  mainfont:
    font: Helvetica Neue Italic
    supported: false
    missingCount: 1
    missing:
      - char: →
        codepoint: U+2192
        name: RIGHTWARDS ARROW
  sansfont:
    font: Arial
    supported: true
    missingCount: 0
  monofont:
    font: Menlo Regular
    supported: true
    missingCount: 0
```

### List available fonts

```bash
# List system fonts
fontscan list

# List fonts in a specific directory
fontscan list --path /Library/Fonts

# Filter fonts by name/family (case-insensitive regex)
fontscan list --filter "helvetica|arial"

# List unique font families only
fontscan list --filter "mono|courier" --family-only

# Limit results
fontscan list --filter "noto" --family-only --limit 10

# JSON output
fontscan list --format json
```

The `--filter` flag accepts a case-insensitive regular expression that matches against font family and name. Combined with `--family-only`, this replaces common `fc-list` workflows:

| fc-list command | fontscan equivalent |
|-----------------|---------------------|
| `fc-list : family \| grep -i "helvetica" \| sort -u` | `fontscan list --filter "helvetica" --family-only` |
| `fc-list : family \| grep -i "mono" \| sort -u \| head -10` | `fontscan list --filter "mono" --family-only --limit 10` |

### Show font information

```bash
# Basic info
fontscan info "Helvetica"

# With Unicode block coverage
fontscan info "Menlo" --coverage
```

Output with `--coverage`:

```
name: Menlo Regular
family: Menlo
style: Regular
path: /System/Library/Fonts/Menlo.ttc
numGlyphs: 5765
coverage:
  - block: Basic Latin
    supported: 95/128 (74.2%)
  - block: Latin-1 Supplement
    supported: 96/128 (75.0%)
  - block: Box Drawing
    supported: 128/128 (100.0%)
  ...
```

### Recommend fonts for a document

```bash
# Find fonts that support all characters in a document
fontscan recommend document.md

# Limit partial matches
fontscan recommend document.md --limit 5
```

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/grokify/fontscan/checker"
    "github.com/grokify/fontscan/font"
    "github.com/grokify/fontscan/pandoc"
)

func main() {
    // Load a font
    f, err := font.LoadFile("/System/Library/Fonts/Menlo.ttc")
    if err != nil {
        panic(err)
    }

    // Check if a character is supported
    if f.HasGlyph('→') {
        fmt.Println("Arrow is supported")
    }

    // Check text against a font
    result := checker.CheckText("Hello → World", f)
    fmt.Printf("Supported: %v, Missing: %d\n", result.Supported, result.MissingCount)

    // Parse Pandoc frontmatter
    settings, _ := pandoc.ParseFrontmatter("document.md")
    for variable, fontName := range settings.FontMap() {
        fmt.Printf("%s: %s\n", variable, fontName)
    }
}
```

## Supported Font Formats

| Format | Extension | Notes |
|--------|-----------|-------|
| TrueType | `.ttf` | Full support |
| OpenType | `.otf` | Full support |
| TrueType Collection | `.ttc` | Loads first working font in collection |

## Output Formats

- `toon` (default) - Human-readable YAML-like format
- `json` - Pretty-printed JSON
- `json-compact` - Single-line JSON

Use `--format` or `-f` to specify:

```bash
fontscan check doc.md --font Arial --format json
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All characters supported |
| 1 | Missing characters found or error occurred |

## Pandoc Font Variables

The `pandoc` command recognizes these standard font variables:

- `mainfont` - Main document font
- `sansfont` - Sans-serif font
- `monofont` - Monospace font (code blocks)
- `mathfont` - Math font
- `CJKmainfont` - CJK main font
- `CJKsansfont` - CJK sans font
- `CJKmonofont` - CJK mono font
- `titlefont` - Title font
- `headerfont` - Header font
- `footerfont` - Footer font

## Known Limitations

- Some fonts use cmap encodings not supported by Go's `sfnt` library. When loading a TTC collection, fontscan iterates through fonts until one loads successfully.
- Coverage analysis checks common Unicode blocks but not all possible code points.

## License

MIT License - see [LICENSE](LICENSE) for details.
