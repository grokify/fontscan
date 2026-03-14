// Package pandoc provides Pandoc Markdown-specific functionality.
package pandoc

import (
	"bufio"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// StandardFontVariables lists the standard Pandoc font variables.
var StandardFontVariables = []string{
	"mainfont",
	"sansfont",
	"monofont",
	"mathfont",
	"CJKmainfont",
	"CJKsansfont",
	"CJKmonofont",
	"titlefont",
	"footerfont",
	"headerfont",
}

// FontSettings holds font configuration extracted from Pandoc frontmatter.
type FontSettings struct {
	MainFont    string            `yaml:"mainfont"`
	SansFont    string            `yaml:"sansfont"`
	MonoFont    string            `yaml:"monofont"`
	MathFont    string            `yaml:"mathfont"`
	CJKMainFont string            `yaml:"CJKmainfont"`
	CJKSansFont string            `yaml:"CJKsansfont"`
	CJKMonoFont string            `yaml:"CJKmonofont"`
	TitleFont   string            `yaml:"titlefont"`
	FooterFont  string            `yaml:"footerfont"`
	HeaderFont  string            `yaml:"headerfont"`
	Custom      map[string]string `yaml:"-"` // For future custom variable support
}

// Fonts returns a list of all specified fonts (non-empty values).
func (fs *FontSettings) Fonts() []string {
	var fonts []string
	seen := make(map[string]bool)

	add := func(f string) {
		if f != "" && !seen[f] {
			seen[f] = true
			fonts = append(fonts, f)
		}
	}

	add(fs.MainFont)
	add(fs.SansFont)
	add(fs.MonoFont)
	add(fs.MathFont)
	add(fs.CJKMainFont)
	add(fs.CJKSansFont)
	add(fs.CJKMonoFont)
	add(fs.TitleFont)
	add(fs.FooterFont)
	add(fs.HeaderFont)

	for _, f := range fs.Custom {
		add(f)
	}

	return fonts
}

// FontMap returns a map of variable name to font name.
func (fs *FontSettings) FontMap() map[string]string {
	m := make(map[string]string)

	if fs.MainFont != "" {
		m["mainfont"] = fs.MainFont
	}
	if fs.SansFont != "" {
		m["sansfont"] = fs.SansFont
	}
	if fs.MonoFont != "" {
		m["monofont"] = fs.MonoFont
	}
	if fs.MathFont != "" {
		m["mathfont"] = fs.MathFont
	}
	if fs.CJKMainFont != "" {
		m["CJKmainfont"] = fs.CJKMainFont
	}
	if fs.CJKSansFont != "" {
		m["CJKsansfont"] = fs.CJKSansFont
	}
	if fs.CJKMonoFont != "" {
		m["CJKmonofont"] = fs.CJKMonoFont
	}
	if fs.TitleFont != "" {
		m["titlefont"] = fs.TitleFont
	}
	if fs.FooterFont != "" {
		m["footerfont"] = fs.FooterFont
	}
	if fs.HeaderFont != "" {
		m["headerfont"] = fs.HeaderFont
	}

	for k, v := range fs.Custom {
		m[k] = v
	}

	return m
}

// ParseFrontmatter extracts YAML frontmatter from a Pandoc Markdown file.
func ParseFrontmatter(path string) (*FontSettings, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)

	// Check for opening ---
	if !scanner.Scan() {
		return &FontSettings{}, nil
	}
	if strings.TrimSpace(scanner.Text()) != "---" {
		return &FontSettings{}, nil
	}

	// Collect YAML content until closing ---
	var yamlLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "..." {
			break
		}
		yamlLines = append(yamlLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(yamlLines) == 0 {
		return &FontSettings{}, nil
	}

	yamlContent := strings.Join(yamlLines, "\n")

	var settings FontSettings
	if err := yaml.Unmarshal([]byte(yamlContent), &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// ParseFrontmatterString extracts font settings from YAML content.
func ParseFrontmatterString(content string) (*FontSettings, error) {
	var settings FontSettings
	if err := yaml.Unmarshal([]byte(content), &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}
