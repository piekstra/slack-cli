// Package version provides build-time version information.
// Variables are set via ldflags during build (see Makefile).
package version

// Build-time variables set via ldflags
var (
	// Version is the semantic version (from git tag or "dev")
	Version = "dev"

	// Commit is the git commit hash
	Commit = "unknown"

	// Date is the build timestamp
	Date = "unknown"
)

// Info returns a formatted version string
func Info() string {
	return Version + " (commit: " + Commit + ", built: " + Date + ")"
}
