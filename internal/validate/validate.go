package validate

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	channelIDRegex = regexp.MustCompile(`^[CG][A-Z0-9]+$`)
	userIDRegex    = regexp.MustCompile(`^[UW][A-Z0-9]+$`)
	timestampRegex = regexp.MustCompile(`^\d+\.\d+$`)
)

// ChannelID validates that the given string is a valid Slack channel ID.
// Channel IDs start with C (public) or G (private/group).
func ChannelID(id string) error {
	if !channelIDRegex.MatchString(id) {
		return fmt.Errorf("invalid channel ID %q: must start with C or G (e.g., C01234ABCDE)", id)
	}
	return nil
}

// UserID validates that the given string is a valid Slack user ID.
// User IDs start with U (regular user) or W (enterprise user).
func UserID(id string) error {
	if !userIDRegex.MatchString(id) {
		return fmt.Errorf("invalid user ID %q: must start with U or W (e.g., U01234ABCDE)", id)
	}
	return nil
}

// Timestamp validates that the given string is a valid Slack message timestamp.
// Timestamps are in the format "1234567890.123456".
func Timestamp(ts string) error {
	if !timestampRegex.MatchString(ts) {
		return fmt.Errorf("invalid timestamp %q: must be format 1234567890.123456", ts)
	}
	return nil
}

// Emoji normalizes an emoji name by stripping surrounding colons.
// Returns the cleaned emoji name.
func Emoji(emoji string) string {
	return strings.Trim(emoji, ":")
}

// Limit validates that the given limit is within acceptable bounds.
func Limit(limit int) error {
	if limit < 1 {
		return fmt.Errorf("invalid limit %d: must be at least 1", limit)
	}
	if limit > 1000 {
		return fmt.Errorf("invalid limit %d: must be at most 1000", limit)
	}
	return nil
}
