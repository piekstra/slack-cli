package validate

import (
	"testing"
)

func TestChannelID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid public channel", "C01234ABCDE", false},
		{"valid private channel", "G01234ABCDE", false},
		{"valid short channel", "C123", false},
		{"invalid prefix", "X01234ABCDE", true},
		{"invalid lowercase", "c01234abcde", true},
		{"empty string", "", true},
		{"user ID", "U01234ABCDE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ChannelID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChannelID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestUserID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid user", "U01234ABCDE", false},
		{"valid enterprise user", "W01234ABCDE", false},
		{"valid short user", "U123", false},
		{"invalid prefix", "X01234ABCDE", true},
		{"invalid lowercase", "u01234abcde", true},
		{"empty string", "", true},
		{"channel ID", "C01234ABCDE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		ts      string
		wantErr bool
	}{
		{"valid timestamp", "1234567890.123456", false},
		{"valid short", "123.456", false},
		{"missing decimal", "1234567890", true},
		{"empty string", "", true},
		{"letters", "abc.def", true},
		{"no digits after decimal", "123.", true},
		{"no digits before decimal", ".123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Timestamp(tt.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp(%q) error = %v, wantErr %v", tt.ts, err, tt.wantErr)
			}
		})
	}
}

func TestEmoji(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"thumbsup", "thumbsup"},
		{":thumbsup:", "thumbsup"},
		{":thumbsup", "thumbsup"},
		{"thumbsup:", "thumbsup"},
		{"::thumbsup::", "thumbsup"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Emoji(tt.input)
			if got != tt.want {
				t.Errorf("Emoji(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		wantErr bool
	}{
		{"valid small", 1, false},
		{"valid medium", 100, false},
		{"valid max", 1000, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"too large", 1001, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Limit(tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Limit(%d) error = %v, wantErr %v", tt.limit, err, tt.wantErr)
			}
		})
	}
}
