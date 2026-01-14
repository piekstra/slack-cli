package search

import (
	"fmt"
	"regexp"
	"strings"
)

// QueryOptions contains options for building search queries
type QueryOptions struct {
	Scope       string
	InChannel   string
	FromUser    string
	After       string
	Before      string
	HasLink     bool
	HasReaction bool
	HasPin      bool
	FileType    string
}

// ValidScopes contains the allowed scope values
var ValidScopes = []string{"all", "public", "private", "dm", "mpim"}

// ValidateScope checks if the scope value is valid
func ValidateScope(scope string) error {
	if scope == "" {
		return nil
	}
	for _, s := range ValidScopes {
		if scope == s {
			return nil
		}
	}
	return fmt.Errorf("invalid scope: %q (must be one of: %s)", scope, strings.Join(ValidScopes, ", "))
}

// datePattern matches YYYY-MM-DD format
var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// ValidateDate checks if the date string is in YYYY-MM-DD format
func ValidateDate(date string) error {
	if date == "" {
		return nil
	}
	if !datePattern.MatchString(date) {
		return fmt.Errorf("invalid date format: %q (must be YYYY-MM-DD)", date)
	}
	return nil
}

// BuildQuery constructs a Slack search query from base query and options
func BuildQuery(baseQuery string, opts *QueryOptions) string {
	if opts == nil {
		return baseQuery
	}

	var parts []string

	// Add scope modifier
	if opts.Scope != "" && opts.Scope != "all" {
		parts = append(parts, "is:"+opts.Scope)
	}

	// Add channel filter
	if opts.InChannel != "" {
		channel := strings.TrimPrefix(opts.InChannel, "#")
		parts = append(parts, "in:#"+channel)
	}

	// Add user filter
	if opts.FromUser != "" {
		user := strings.TrimPrefix(opts.FromUser, "@")
		parts = append(parts, "from:@"+user)
	}

	// Add date filters
	if opts.After != "" {
		parts = append(parts, "after:"+opts.After)
	}
	if opts.Before != "" {
		parts = append(parts, "before:"+opts.Before)
	}

	// Add has filters
	if opts.HasLink {
		parts = append(parts, "has:link")
	}
	if opts.HasReaction {
		parts = append(parts, "has:reaction")
	}
	if opts.HasPin {
		parts = append(parts, "has:pin")
	}

	// Add file type filter
	if opts.FileType != "" {
		parts = append(parts, "type:"+opts.FileType)
	}

	// Add the base query at the end
	parts = append(parts, baseQuery)

	return strings.Join(parts, " ")
}

// ValidateQueryOptions validates all query options
func ValidateQueryOptions(opts *QueryOptions) error {
	if opts == nil {
		return nil
	}

	if err := ValidateScope(opts.Scope); err != nil {
		return err
	}
	if err := ValidateDate(opts.After); err != nil {
		return err
	}
	if err := ValidateDate(opts.Before); err != nil {
		return err
	}

	return nil
}
