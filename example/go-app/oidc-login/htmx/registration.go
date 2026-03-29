package htmx

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	services_htmx_error "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/error"
	services_htmx_forgotpassword "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/forgotpassword"
	services_htmx_home "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/home"
	services_htmx_keepsignedin "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/keepsignedin"
	services_htmx_password "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/password"
	services_htmx_resetpassword "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/resetpassword"
	services_htmx_shell "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/shell"
	services_htmx_signup "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/signup"
	services_htmx_startover "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/startover"
	services_htmx_verifycode "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/verifycode"
)

// AddOIDCLoginHandlers registers all HTMX-based OIDC login handlers.
// This is the default OIDC login UI implementation.
// To use a different UI (e.g. WASM SPA), do not call this function and instead
// register your own handler for the /oidc-login/* paths.
func AddOIDCLoginHandlers(builder di.ContainerBuilder) {
	services_htmx_shell.AddScopedIHandler(builder)
	services_htmx_home.AddScopedIHandler(builder)
	services_htmx_password.AddScopedIHandler(builder)
	services_htmx_verifycode.AddScopedIHandler(builder)
	services_htmx_keepsignedin.AddScopedIHandler(builder)
	services_htmx_signup.AddScopedIHandler(builder)
	services_htmx_forgotpassword.AddScopedIHandler(builder)
	services_htmx_resetpassword.AddScopedIHandler(builder)
	services_htmx_error.AddScopedIHandler(builder)
	services_htmx_startover.AddScopedIHandler(builder)
}
