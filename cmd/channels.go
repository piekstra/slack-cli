package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/piekstra/slack-cli/internal/client"
)

var channelsCmd = &cobra.Command{
	Use:     "channels",
	Aliases: []string{"ch"},
	Short:   "Manage Slack channels",
}

var listChannelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all channels",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		types, _ := cmd.Flags().GetString("types")
		excludeArchived, _ := cmd.Flags().GetBool("exclude-archived")
		limit, _ := cmd.Flags().GetInt("limit")

		channels, err := c.ListChannels(types, excludeArchived, limit)
		if err != nil {
			return err
		}

		if outputJSON {
			data, _ := json.MarshalIndent(channels, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(channels) == 0 {
			fmt.Println("No channels found")
			return nil
		}

		fmt.Printf("%-12s %-30s %s\n", "ID", "NAME", "MEMBERS")
		fmt.Println(strings.Repeat("-", 60))
		for _, ch := range channels {
			private := ""
			if ch.IsPrivate {
				private = " (private)"
			}
			fmt.Printf("%-12s %-30s %d%s\n", ch.ID, ch.Name, ch.NumMembers, private)
		}

		return nil
	},
}

var getChannelCmd = &cobra.Command{
	Use:   "get <channel-id>",
	Short: "Get channel information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		channel, err := c.GetChannelInfo(args[0])
		if err != nil {
			return err
		}

		if outputJSON {
			data, _ := json.MarshalIndent(channel, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:       %s\n", channel.ID)
		fmt.Printf("Name:     %s\n", channel.Name)
		fmt.Printf("Private:  %t\n", channel.IsPrivate)
		fmt.Printf("Archived: %t\n", channel.IsArchived)
		fmt.Printf("Members:  %d\n", channel.NumMembers)
		if channel.Topic.Value != "" {
			fmt.Printf("Topic:    %s\n", channel.Topic.Value)
		}
		if channel.Purpose.Value != "" {
			fmt.Printf("Purpose:  %s\n", channel.Purpose.Value)
		}

		return nil
	},
}

var createChannelCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		isPrivate, _ := cmd.Flags().GetBool("private")

		channel, err := c.CreateChannel(args[0], isPrivate)
		if err != nil {
			return err
		}

		if outputJSON {
			data, _ := json.MarshalIndent(channel, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Created channel: %s (%s)\n", channel.Name, channel.ID)
		return nil
	},
}

var archiveChannelCmd = &cobra.Command{
	Use:   "archive <channel-id>",
	Short: "Archive a channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		if err := c.ArchiveChannel(args[0]); err != nil {
			return err
		}

		fmt.Printf("Archived channel: %s\n", args[0])
		return nil
	},
}

var unarchiveChannelCmd = &cobra.Command{
	Use:   "unarchive <channel-id>",
	Short: "Unarchive a channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		if err := c.UnarchiveChannel(args[0]); err != nil {
			return err
		}

		fmt.Printf("Unarchived channel: %s\n", args[0])
		return nil
	},
}

var setTopicCmd = &cobra.Command{
	Use:   "set-topic <channel-id> <topic>",
	Short: "Set channel topic",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		if err := c.SetChannelTopic(args[0], args[1]); err != nil {
			return err
		}

		fmt.Printf("Set topic for channel %s\n", args[0])
		return nil
	},
}

var setPurposeCmd = &cobra.Command{
	Use:   "set-purpose <channel-id> <purpose>",
	Short: "Set channel purpose",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		if err := c.SetChannelPurpose(args[0], args[1]); err != nil {
			return err
		}

		fmt.Printf("Set purpose for channel %s\n", args[0])
		return nil
	},
}

var inviteCmd = &cobra.Command{
	Use:   "invite <channel-id> <user-id>...",
	Short: "Invite users to a channel",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		if err := c.InviteToChannel(args[0], args[1:]); err != nil {
			return err
		}

		fmt.Printf("Invited %d user(s) to channel %s\n", len(args)-1, args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(channelsCmd)

	channelsCmd.AddCommand(listChannelsCmd)
	listChannelsCmd.Flags().String("types", "", "Channel types (public_channel,private_channel,mpim,im)")
	listChannelsCmd.Flags().Bool("exclude-archived", true, "Exclude archived channels")
	listChannelsCmd.Flags().Int("limit", 100, "Maximum channels to return")

	channelsCmd.AddCommand(getChannelCmd)

	channelsCmd.AddCommand(createChannelCmd)
	createChannelCmd.Flags().Bool("private", false, "Create as private channel")

	channelsCmd.AddCommand(archiveChannelCmd)
	channelsCmd.AddCommand(unarchiveChannelCmd)
	channelsCmd.AddCommand(setTopicCmd)
	channelsCmd.AddCommand(setPurposeCmd)
	channelsCmd.AddCommand(inviteCmd)
}
