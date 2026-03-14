package pandoc

// DefaultFonts lists common Pandoc/LaTeX default fonts.
// These are fonts that Pandoc may use when no explicit font is specified.
var DefaultFonts = []FontInfo{
	// Latin Modern family (default for many LaTeX setups)
	{Name: "Latin Modern Roman", Family: "lmroman", Variants: []string{
		"lmroman10-regular", "lmroman10-bold", "lmroman10-italic", "lmroman10-bolditalic",
		"lmroman12-regular", "lmroman12-bold", "lmroman12-italic",
		"lmroman17-regular",
	}},
	{Name: "Latin Modern Sans", Family: "lmsans", Variants: []string{
		"lmsans10-regular", "lmsans10-bold", "lmsans10-oblique", "lmsans10-boldoblique",
		"lmsans12-regular", "lmsans12-oblique",
	}},
	{Name: "Latin Modern Mono", Family: "lmmono", Variants: []string{
		"lmmono10-regular", "lmmono10-italic",
		"lmmono12-regular",
		"lmmonolt10-regular", "lmmonolt10-bold", "lmmonolt10-oblique", "lmmonolt10-boldoblique",
		"lmmonoltcond10-regular", "lmmonoltcond10-oblique",
		"lmmonoprop10-regular", "lmmonoprop10-oblique",
		"lmmonoproplt10-regular", "lmmonoproplt10-oblique", "lmmonoproplt10-bold", "lmmonoproplt10-boldoblique",
	}},

	// Computer Modern family (TeX classic)
	{Name: "Computer Modern Roman", Family: "cmr", Variants: []string{
		"cmr10", "cmr12", "cmr17",
	}},
	{Name: "Computer Modern Sans", Family: "cmss", Variants: []string{
		"cmss10", "cmss12",
	}},
	{Name: "Computer Modern Typewriter", Family: "cmtt", Variants: []string{
		"cmtt10", "cmtt12",
	}},

	// TeX Gyre family (extended Latin Modern)
	{Name: "TeX Gyre Termes", Family: "texgyretermes", Variants: []string{
		"texgyretermes-regular", "texgyretermes-bold", "texgyretermes-italic", "texgyretermes-bolditalic",
	}},
	{Name: "TeX Gyre Pagella", Family: "texgyrepagella", Variants: []string{
		"texgyrepagella-regular", "texgyrepagella-bold", "texgyrepagella-italic", "texgyrepagella-bolditalic",
	}},
	{Name: "TeX Gyre Heros", Family: "texgyreheros", Variants: []string{
		"texgyreheros-regular", "texgyreheros-bold", "texgyreheros-italic", "texgyreheros-bolditalic",
	}},
	{Name: "TeX Gyre Cursor", Family: "texgyrecursor", Variants: []string{
		"texgyrecursor-regular", "texgyrecursor-bold", "texgyrecursor-italic", "texgyrecursor-bolditalic",
	}},

	// DejaVu family (good Unicode coverage)
	{Name: "DejaVu Serif", Family: "dejavuserif", Variants: []string{
		"DejaVuSerif", "DejaVuSerif-Bold", "DejaVuSerif-Italic", "DejaVuSerif-BoldItalic",
	}},
	{Name: "DejaVu Sans", Family: "dejavusans", Variants: []string{
		"DejaVuSans", "DejaVuSans-Bold", "DejaVuSans-Oblique", "DejaVuSans-BoldOblique",
	}},
	{Name: "DejaVu Sans Mono", Family: "dejavusansmono", Variants: []string{
		"DejaVuSansMono", "DejaVuSansMono-Bold", "DejaVuSansMono-Oblique", "DejaVuSansMono-BoldOblique",
	}},

	// Libertinus family (OpenType successor to Linux Libertine)
	{Name: "Libertinus Serif", Family: "libertinusserif", Variants: []string{
		"LibertinusSerif-Regular", "LibertinusSerif-Bold", "LibertinusSerif-Italic", "LibertinusSerif-BoldItalic",
	}},
	{Name: "Libertinus Sans", Family: "libertinussans", Variants: []string{
		"LibertinusSans-Regular", "LibertinusSans-Bold", "LibertinusSans-Italic",
	}},
	{Name: "Libertinus Mono", Family: "libertinusmono", Variants: []string{
		"LibertinusMono-Regular",
	}},
}

// FontInfo describes a font family and its variants.
type FontInfo struct {
	Name     string   // Human-readable name
	Family   string   // Font family identifier
	Variants []string // List of variant names
}

// AllDefaultFontNames returns all variant names from DefaultFonts.
func AllDefaultFontNames() []string {
	var names []string
	for _, f := range DefaultFonts {
		names = append(names, f.Variants...)
	}
	return names
}

// IsDefaultFont checks if a font name matches a known Pandoc/LaTeX default.
func IsDefaultFont(name string) bool {
	for _, f := range DefaultFonts {
		if f.Family == name || f.Name == name {
			return true
		}
		for _, v := range f.Variants {
			if v == name {
				return true
			}
		}
	}
	return false
}

// GetFontInfo returns info about a default font, or nil if not found.
func GetFontInfo(name string) *FontInfo {
	for i, f := range DefaultFonts {
		if f.Family == name || f.Name == name {
			return &DefaultFonts[i]
		}
		for _, v := range f.Variants {
			if v == name {
				return &DefaultFonts[i]
			}
		}
	}
	return nil
}
