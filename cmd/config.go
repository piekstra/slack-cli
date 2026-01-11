package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piekstra/slack-cli/internal/keychain"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure slack-cli",
}

var setTokenCmd = &cobra.Command{
	Use:   "set-token [token]",
	Short: "Set the Slack API token",
	Long: `Set the Slack API token for authentication.

On macOS: Token is stored securely in the system Keychain.
On Linux: Token is stored in ~/.config/slack-cli/credentials (file permissions 0600).

If no token is provided as an argument, you will be prompted to enter it.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Warn Linux users about file-based storage
		if !keychain.IsSecureStorage() {
			fmt.Println("Warning: On Linux, your token will be stored in a config file")
			fmt.Println("         (~/.config/slack-cli/credentials) with restricted permissions (0600).")
			fmt.Println("         This is less secure than macOS Keychain storage.")
			fmt.Println()
		}

		var token string

		if len(args) > 0 {
			token = args[0]
		} else {
			fmt.Print("Enter Slack API token: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			token = strings.TrimSpace(input)
		}

		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}

		if err := keychain.SetAPIToken(token); err != nil {
			return fmt.Errorf("failed to store token: %w", err)
		}

		if keychain.IsSecureStorage() {
			fmt.Println("API token stored securely in Keychain")
		} else {
			fmt.Println("API token stored in ~/.config/slack-cli/credentials")
		}
		return nil
	},
}

var deleteTokenCmd = &cobra.Command{
	Use:   "delete-token",
	Short: "Delete the stored Slack API token",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := keychain.DeleteAPIToken(); err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}

		if keychain.IsSecureStorage() {
			fmt.Println("API token deleted from Keychain")
		} else {
			fmt.Println("API token deleted from config file")
		}
		return nil
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := keychain.GetAPIToken()
		if err != nil {
			fmt.Println("API Token: Not configured")
			fmt.Println("\nRun 'slack-cli config set-token' to configure")
			return nil
		}

		// Mask the token for display
		masked := token[:8] + strings.Repeat("*", len(token)-12) + token[len(token)-4:]

		if keychain.IsSecureStorage() {
			fmt.Printf("API Token: %s (from Keychain)\n", masked)
		} else {
			fmt.Printf("API Token: %s (from config file)\n", masked)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setTokenCmd)
	configCmd.AddCommand(deleteTokenCmd)
	configCmd.AddCommand(showConfigCmd)
}
