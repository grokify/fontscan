package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags during build)
var (
	version = "0.1.0"
	commit  = "dev"
	date    = "unknown"
)

var (
	outputFormat string
	showVersion  bool
)

var rootCmd = &cobra.Command{
	Use:   "fontscan",
	Short: "Check font glyph coverage for documents",
	Long: `fontscan is a tool for checking if fonts support all characters in a document.

It helps prevent Pandoc/LaTeX missing character errors by verifying glyph coverage
before document generation.

Examples:
  fontscan check document.md --font "Helvetica Neue"
  fontscan list
  fontscan info "DejaVu Sans Mono" --coverage
  fontscan recommend document.md`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("fontscan %s (%s) built %s\n", version, commit, date)
			return
		}
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "toon",
		"Output format: toon, json, json-compact")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(recommendCmd)
}
