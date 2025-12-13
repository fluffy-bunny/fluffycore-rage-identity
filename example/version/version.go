package version

// Values for these are injected by the build
var (
	version = "dev-build"
	commit  = ""
	date    = ""
)

// Version returns the gateways version.
func Version() string {
	return version
}

// Commit returns the git commit SHA for the code that gateways micro was built from.
func Commit() string {
	return commit
}

// Date returns the date for the code that gateways micro was built from.
func Date() string {
	return date
}
