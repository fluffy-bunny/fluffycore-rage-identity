package htmx

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	mgmt_htmx_home "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/home"
	mgmt_htmx_linked "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/linked"
	mgmt_htmx_passkey "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/passkey"
	mgmt_htmx_password "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/password"
	mgmt_htmx_prefs "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/prefs"
	mgmt_htmx_profile "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/profile"
	mgmt_htmx_shell "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/handlers/shell"
)

// AddManagementHandlers registers all HTMX-based management dashboard handlers.
// This is the default management UI implementation.
// To use a different UI (e.g. WASM SPA), do not call this function and instead
// register your own handler for the /management/* paths.
func AddManagementHandlers(builder di.ContainerBuilder) {
	mgmt_htmx_shell.AddScopedIHandler(builder)
	mgmt_htmx_home.AddScopedIHandler(builder)
	mgmt_htmx_profile.AddScopedIHandler(builder)
	mgmt_htmx_password.AddScopedIHandler(builder)
	mgmt_htmx_passkey.AddScopedIHandler(builder)
	mgmt_htmx_linked.AddScopedIHandler(builder)
	mgmt_htmx_prefs.AddScopedIHandler(builder)
}
