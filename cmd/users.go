package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piekstra/slack-cli/internal/client"
)

var usersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"u"},
	Short:   "Manage Slack users",
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")

		users, err := c.ListUsers(limit)
		if err != nil {
			return err
		}

		if outputJSON {
			data, _ := json.MarshalIndent(users, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(users) == 0 {
			fmt.Println("No users found")
			return nil
		}

		fmt.Printf("%-12s %-20s %-30s %s\n", "ID", "USERNAME", "REAL NAME", "EMAIL")
		fmt.Println(strings.Repeat("-", 90))
		for _, u := range users {
			if u.IsBot {
				continue
			}
			fmt.Printf("%-12s %-20s %-30s %s\n", u.ID, u.Name, u.RealName, u.Profile.Email)
		}

		return nil
	},
}

var getUserCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get user information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		user, err := c.GetUserInfo(args[0])
		if err != nil {
			return err
		}

		if outputJSON {
			data, _ := json.MarshalIndent(user, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:           %s\n", user.ID)
		fmt.Printf("Username:     %s\n", user.Name)
		fmt.Printf("Real Name:    %s\n", user.RealName)
		fmt.Printf("Display Name: %s\n", user.Profile.DisplayName)
		fmt.Printf("Email:        %s\n", user.Profile.Email)
		fmt.Printf("Admin:        %t\n", user.IsAdmin)
		fmt.Printf("Bot:          %t\n", user.IsBot)
		if user.Profile.StatusText != "" {
			fmt.Printf("Status:       %s %s\n", user.Profile.StatusEmoji, user.Profile.StatusText)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().Int("limit", 100, "Maximum users to return")

	usersCmd.AddCommand(getUserCmd)
}
