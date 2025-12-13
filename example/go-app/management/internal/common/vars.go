package common

var (
	// Those vars are used during the build --- linker (ld) bakes the values in.
	AppVersion = "dev"     // Application version
	BuildTime  = "unknown" // Build timestamp
	GitCommit  = "unknown" // Git commit hash
	GitBranch  = "unknown" // Git branch name
)
