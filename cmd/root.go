package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

type globalOptions struct {
	outputJSON bool
	cacheTTL   time.Duration
}

var (
	rootOpts = &globalOptions{
		cacheTTL: 30 * time.Second,
	}

	rootCmd = &cobra.Command{
		Use:   "fpl",
		Short: "Interact with the Fantasy Premier League API from the terminal",
		Long: `fpl is a Fantasy Premier League helper that wraps the public FPL API.

It can resolve players either by ID or name (with fuzzy matching) and print their
gameweek-by-gameweek performance in human-friendly tables or machine-friendly JSON.`,
		Version:       "dev",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(
		&rootOpts.outputJSON,
		"json",
		false,
		"print structured JSON output instead of a table",
	)
	rootCmd.PersistentFlags().DurationVar(
		&rootOpts.cacheTTL,
		"cache-ttl",
		rootOpts.cacheTTL,
		"cache duration for bootstrap-static requests (set to 0 to disable caching)",
	)
}
