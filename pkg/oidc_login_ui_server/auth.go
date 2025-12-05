package oidc_login_ui_server

var requiresNoAuthPaths map[string]bool

// everything requries auth unless otherwise documented here.
// -- this is a list of paths that do not require auth
func RequiresNoAuth() map[string]bool {
	// needs to be a func as some of these are configured in.
	if requiresNoAuthPaths == nil {
		requiresNoAuthPaths = map[string]bool{}
	}
	return requiresNoAuthPaths
}
