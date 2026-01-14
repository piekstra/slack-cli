package search

import (
	"testing"
)

func TestValidateScope(t *testing.T) {
	tests := []struct {
		name      string
		scope     string
		wantError bool
	}{
		{"empty scope", "", false},
		{"all scope", "all", false},
		{"public scope", "public", false},
		{"private scope", "private", false},
		{"dm scope", "dm", false},
		{"mpim scope", "mpim", false},
		{"invalid scope", "invalid", true},
		{"random scope", "random", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScope(tt.scope)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateScope(%q) error = %v, wantError %v", tt.scope, err, tt.wantError)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		wantError bool
	}{
		{"empty date", "", false},
		{"valid date", "2025-01-15", false},
		{"valid date 2", "2024-12-31", false},
		{"invalid format - no dashes", "20250115", true},
		{"invalid format - wrong separator", "2025/01/15", true},
		{"invalid format - text", "invalid-date", true},
		{"invalid format - partial", "2025-01", true},
		{"invalid format - extra", "2025-01-15-01", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDate(tt.date)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateDate(%q) error = %v, wantError %v", tt.date, err, tt.wantError)
			}
		})
	}
}

func TestBuildQuery_Empty(t *testing.T) {
	result := BuildQuery("test query", nil)
	if result != "test query" {
		t.Errorf("BuildQuery with nil opts = %q, want %q", result, "test query")
	}

	result = BuildQuery("test query", &QueryOptions{})
	if result != "test query" {
		t.Errorf("BuildQuery with empty opts = %q, want %q", result, "test query")
	}
}

func TestBuildQuery_Scope(t *testing.T) {
	tests := []struct {
		name  string
		scope string
		want  string
	}{
		{"all scope", "all", "test query"},
		{"public scope", "public", "is:public test query"},
		{"private scope", "private", "is:private test query"},
		{"dm scope", "dm", "is:dm test query"},
		{"mpim scope", "mpim", "is:mpim test query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &QueryOptions{Scope: tt.scope}
			result := BuildQuery("test query", opts)
			if result != tt.want {
				t.Errorf("BuildQuery with scope %q = %q, want %q", tt.scope, result, tt.want)
			}
		})
	}
}

func TestBuildQuery_InChannel(t *testing.T) {
	tests := []struct {
		name    string
		channel string
		want    string
	}{
		{"channel without hash", "general", "in:#general test query"},
		{"channel with hash", "#general", "in:#general test query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &QueryOptions{InChannel: tt.channel}
			result := BuildQuery("test query", opts)
			if result != tt.want {
				t.Errorf("BuildQuery with channel %q = %q, want %q", tt.channel, result, tt.want)
			}
		})
	}
}

func TestBuildQuery_FromUser(t *testing.T) {
	tests := []struct {
		name string
		user string
		want string
	}{
		{"user without at", "alice", "from:@alice test query"},
		{"user with at", "@alice", "from:@alice test query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &QueryOptions{FromUser: tt.user}
			result := BuildQuery("test query", opts)
			if result != tt.want {
				t.Errorf("BuildQuery with user %q = %q, want %q", tt.user, result, tt.want)
			}
		})
	}
}

func TestBuildQuery_DateRange(t *testing.T) {
	tests := []struct {
		name   string
		after  string
		before string
		want   string
	}{
		{"after only", "2025-01-01", "", "after:2025-01-01 test query"},
		{"before only", "", "2025-12-31", "before:2025-12-31 test query"},
		{"both dates", "2025-01-01", "2025-12-31", "after:2025-01-01 before:2025-12-31 test query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &QueryOptions{After: tt.after, Before: tt.before}
			result := BuildQuery("test query", opts)
			if result != tt.want {
				t.Errorf("BuildQuery with dates = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestBuildQuery_HasFlags(t *testing.T) {
	tests := []struct {
		name        string
		hasLink     bool
		hasReaction bool
		hasPin      bool
		want        string
	}{
		{"has link", true, false, false, "has:link test query"},
		{"has reaction", false, true, false, "has:reaction test query"},
		{"has pin", false, false, true, "has:pin test query"},
		{"has link and reaction", true, true, false, "has:link has:reaction test query"},
		{"all has flags", true, true, true, "has:link has:reaction has:pin test query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &QueryOptions{
				HasLink:     tt.hasLink,
				HasReaction: tt.hasReaction,
				HasPin:      tt.hasPin,
			}
			result := BuildQuery("test query", opts)
			if result != tt.want {
				t.Errorf("BuildQuery with has flags = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestBuildQuery_FileType(t *testing.T) {
	opts := &QueryOptions{FileType: "pdf"}
	result := BuildQuery("test query", opts)
	want := "type:pdf test query"
	if result != want {
		t.Errorf("BuildQuery with file type = %q, want %q", result, want)
	}
}

func TestBuildQuery_Combined(t *testing.T) {
	opts := &QueryOptions{
		Scope:       "public",
		InChannel:   "#general",
		FromUser:    "@alice",
		After:       "2025-01-01",
		Before:      "2025-12-31",
		HasLink:     true,
		HasReaction: true,
	}
	result := BuildQuery("meeting notes", opts)
	want := "is:public in:#general from:@alice after:2025-01-01 before:2025-12-31 has:link has:reaction meeting notes"
	if result != want {
		t.Errorf("BuildQuery combined = %q, want %q", result, want)
	}
}

func TestValidateQueryOptions(t *testing.T) {
	tests := []struct {
		name      string
		opts      *QueryOptions
		wantError bool
	}{
		{"nil opts", nil, false},
		{"empty opts", &QueryOptions{}, false},
		{"valid scope", &QueryOptions{Scope: "public"}, false},
		{"invalid scope", &QueryOptions{Scope: "invalid"}, true},
		{"valid dates", &QueryOptions{After: "2025-01-01", Before: "2025-12-31"}, false},
		{"invalid after date", &QueryOptions{After: "invalid"}, true},
		{"invalid before date", &QueryOptions{Before: "invalid"}, true},
		{"all valid", &QueryOptions{Scope: "public", After: "2025-01-01", Before: "2025-12-31"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQueryOptions(tt.opts)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateQueryOptions() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
