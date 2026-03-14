# Release Notes - v0.2.0

This release adds library APIs for programmatic access to font listing and Pandoc document checking, improving testability and enabling direct integration without CLI invocation.

## Highlights

- **Library API for font listing** - Use `scanner.List()` with `ListOptions` to filter and list fonts programmatically
- **Library API for Pandoc checking** - Use `pandoc.CheckDocument()` to validate frontmatter fonts without CLI

## New Library Functions

### scanner package

```go
// List fonts with filtering
opts := scanner.ListOptions{
    Filter:     "helvetica|arial",
    FamilyOnly: true,
    Limit:      10,
}
result, err := scanner.NewScanner().List(opts)

// Extract unique families from entries
families := scanner.ExtractUniqueFamilies(result.Entries)
```

**New types:**
- `ListOptions` - Configure filtering, family-only mode, and limits
- `FontEntry` - Represents a font in listings
- `ListResult` - Contains entries or families based on options

**New functions:**
- `(*Scanner).List(ListOptions)` - List fonts with options
- `ExtractUniqueFamilies([]FontEntry)` - Deduplicate font families

### pandoc package

```go
// Check a Pandoc document's fonts
result, err := pandoc.CheckDocument("document.md", pandoc.CheckOptions{
    Paths:        []string{"/custom/fonts"},
    WarnDefaults: true,
})

if result.HasErrors {
    for _, e := range result.Errors {
        fmt.Println(e)
    }
}
```

**New types:**
- `CheckOptions` - Configure paths and default font warnings
- `CheckResult` - Contains font settings, results, warnings, and errors

**New functions:**
- `CheckDocument(path, CheckOptions)` - Check document fonts
- `(*CheckResult).HasFontSettings()` - Check if fonts were found

### format package

```go
// Format scanner results
formatter := format.NewFormatter(os.Stdout, format.FormatTOON)
formatter.WriteListResult(scannerResult)

// Format pandoc results
formatter.WritePandocResult(pandocResult)
```

**New types:**
- `PandocCheckResult` - Pandoc result for formatting

**New functions:**
- `(*Formatter).WriteListResult(*scanner.ListResult)` - Format list results
- `(*Formatter).WritePandocResult(*PandocCheckResult)` - Format pandoc results

## Changed

- CLI commands (`list`, `pandoc`) now delegate to library functions
- Business logic extracted from `cmd/` for better testability

## Tests

- Added unit tests for `scanner.ExtractUniqueFamilies()`
- Added unit tests for `pandoc.CheckDocument()`

## Upgrading

This release is fully backwards compatible. Existing CLI usage remains unchanged. The new library functions are additive.

## Installation

```bash
go install github.com/grokify/fontscan/cmd/fontscan@v0.2.0
```

Or as a library:

```bash
go get github.com/grokify/fontscan@v0.2.0
```
