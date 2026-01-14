package users

import (
	"github.com/spf13/cobra"
)

// NewCmd creates the users command with all subcommands
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users",
		Aliases: []string{"u"},
		Short:   "Manage Slack users",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newSearchCmd())

	return cmd
}
