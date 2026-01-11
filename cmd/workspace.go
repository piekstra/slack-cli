package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piekstra/slack-cli/internal/client"
)

var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws", "team"},
	Short:   "Get workspace information",
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get workspace/team information",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		team, err := c.GetTeamInfo()
		if err != nil {
			return err
		}

		if outputJSON {
			data, _ := json.MarshalIndent(team, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:     %s\n", team.ID)
		fmt.Printf("Name:   %s\n", team.Name)
		fmt.Printf("Domain: %s.slack.com\n", team.Domain)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(infoCmd)
}
