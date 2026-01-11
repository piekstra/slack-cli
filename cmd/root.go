package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/piekstra/slack-cli/internal/version"
)

var (
	outputJSON bool
)

var rootCmd = &cobra.Command{
	Use:   "slack-cli",
	Short: "A CLI tool for interacting with Slack",
	Long: `slack-cli is a command-line interface for Slack.

It provides commands for managing channels, users, messages,
and other Slack workspace operations.

Configure your API token with:
  slack-cli config set-token <your-token>

Or set the SLACK_API_TOKEN environment variable.`,
	Version: version.Version,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output in JSON format")

	// Set custom version template to include commit and build date
	rootCmd.SetVersionTemplate("slack-cli " + version.Info() + "\n")
}
