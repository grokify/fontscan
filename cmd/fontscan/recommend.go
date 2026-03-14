package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/grokify/fontscan/checker"
	"github.com/grokify/fontscan/format"
)

var (
	recommendLimit int
	recommendPaths []string
)

var recommendCmd = &cobra.Command{
	Use:   "recommend <file>",
	Short: "Recommend fonts that support all characters in a file",
	Long: `Find installed fonts that support all characters in the specified file.

Lists fonts with full support first, followed by fonts with partial support
sorted by coverage percentage.

Examples:
  fontscan recommend document.md
  fontscan recommend document.md --limit 10
  fontscan recommend document.md --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runRecommend,
}

func init() {
	recommendCmd.Flags().IntVar(&recommendLimit, "limit", 10,
		"Maximum number of partial-support fonts to show (0 for unlimited)")
	recommendCmd.Flags().StringSliceVar(&recommendPaths, "path", nil,
		"Additional font search paths")
}

func runRecommend(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	var c *checker.Checker
	if len(recommendPaths) > 0 {
		c = checker.NewCheckerWithPaths(recommendPaths)
	} else {
		c = checker.NewChecker()
	}

	result, err := c.Recommend(filePath, recommendLimit)
	if err != nil {
		return fmt.Errorf("failed to analyze file: %w", err)
	}

	formatter := format.NewFormatter(os.Stdout, format.ParseFormat(outputFormat))
	return formatter.WriteRecommendResult(result)
}
