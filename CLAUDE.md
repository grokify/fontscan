# CLAUDE.md

Instructions for AI agents using fontscan.

## Overview

fontscan is a pure-Go CLI and library for checking font glyph coverage. Use it to verify that fonts support all characters in a document before generating PDFs with Pandoc/LaTeX.

## When to Use fontscan

Use fontscan when:

- Generating PDFs with Pandoc/LaTeX and need to verify font compatibility
- A document contains special characters (arrows, math symbols, emoji, CJK, etc.)
- You need to find fonts that support specific Unicode characters
- Checking Pandoc YAML frontmatter font settings

## Quick Reference

### Find fonts by name/family

```bash
# Find sans-serif fonts
fontscan list --filter "helvetica|arial|dejavu|liberation" --family-only --limit 10

# Find monospace fonts
fontscan list --filter "menlo|courier|mono|consolas" --family-only --limit 10

# Find fonts with good Unicode support
fontscan list --filter "dejavu|liberation|noto" --family-only --limit 10
```

### Check document compatibility

```bash
# Check if font supports all characters in document
fontscan check document.md --font "Helvetica Neue"

# Check with specific font file
fontscan check document.md --font-path /Library/Fonts/Arial.ttf
```

### Get font recommendations

```bash
# Find fonts that support all characters in document
fontscan recommend document.md --limit 5
```

### Check Pandoc frontmatter

```bash
# Parse and validate fonts in YAML frontmatter
fontscan pandoc document.md
```

### Get font info

```bash
# Show font details with Unicode coverage
fontscan info "Menlo" --coverage
```

## Workflow for Document/Font Compliance

When a user needs to generate a PDF and ensure font compatibility:

1. **Find candidate fonts** matching their criteria:
   ```bash
   fontscan list --filter "helvetica|arial" --family-only --limit 10
   ```

2. **Check document compatibility** with chosen font:
   ```bash
   fontscan check document.md --font "Helvetica Neue"
   ```

3. **If missing characters**, get recommendations:
   ```bash
   fontscan recommend document.md --limit 5
   ```

4. **For Pandoc documents**, validate frontmatter fonts:
   ```bash
   fontscan pandoc document.md
   ```

## Output Formats

- `toon` (default) - Human-readable YAML-like format, token-efficient
- `json` - Pretty-printed JSON for structured parsing
- `json-compact` - Single-line JSON

Use `--format json` when you need to parse the output programmatically.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All characters supported / success |
| 1 | Missing characters found or error |

## Common Font Categories

| Category | Filter Pattern |
|----------|---------------|
| Sans-serif | `helvetica\|arial\|dejavu sans\|liberation sans\|noto sans` |
| Serif | `times\|georgia\|dejavu serif\|liberation serif\|noto serif` |
| Monospace | `menlo\|courier\|mono\|consolas\|dejavu.*mono` |
| Unicode-rich | `dejavu\|liberation\|noto\|unifont` |
| CJK | `noto.*cjk\|source han\|heiti\|songti` |

## Replacing fc-list

fontscan eliminates the need for fc-list in most workflows:

| fc-list command | fontscan equivalent |
|-----------------|---------------------|
| `fc-list : family \| grep -i "pattern" \| sort -u` | `fontscan list --filter "pattern" --family-only` |
| `fc-list : family \| grep -i "pattern" \| sort -u \| head -N` | `fontscan list --filter "pattern" --family-only --limit N` |

## Troubleshooting

### Font not found

If `fontscan check --font "Name"` fails to find a font:

1. List available fonts to find exact name:
   ```bash
   fontscan list --filter "partial-name" --limit 20
   ```

2. Use the full font name or path:
   ```bash
   fontscan check doc.md --font "Helvetica Neue"
   fontscan check doc.md --font-path /Library/Fonts/HelveticaNeue.ttc
   ```

### Missing characters

When `fontscan check` reports missing characters:

1. Note the Unicode codepoints (e.g., `U+2192 RIGHTWARDS ARROW`)
2. Find fonts supporting those characters:
   ```bash
   fontscan recommend document.md --limit 10
   ```
3. Or check a Unicode-rich font family:
   ```bash
   fontscan check document.md --font "DejaVu Sans"
   ```

## Development

```bash
# Build
go build ./cmd/fontscan

# Test
go test ./...

# Lint
golangci-lint run
```
