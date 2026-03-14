// fontscan is a CLI tool for checking font glyph coverage.
package main

import (
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
