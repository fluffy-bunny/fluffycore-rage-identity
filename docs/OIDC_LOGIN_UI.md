# OIDC Login UI

The OIDC login UI is served by HTMX-based server-rendered handlers using gomponents.

## Architecture

The `/oidc-login` POST handler (`pkg/services/echo/handlers/oidclogin/oidc-login.go`) acts as the entry point for all OIDC authorization flows. When an authorization request comes in, it checks `config.OIDCUIConfig.URIEntryPath`:

- If `URIEntryPath` matches the default `/oidc-login`, the handler processes the request directly.
- If `URIEntryPath` differs (e.g., `/oidc-login/`), it redirects to that path — where your chosen UI implementation takes over.

## HTMX Handlers

The HTMX handlers are registered with a single call in your `MyConfigServices` function:

```go
import (
    services_htmx "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/htmx"
)

func (s *startup) MyConfigServices(ctx context.Context, config *rage_contracts_config.Config, builder di.ContainerBuilder) {
    // ... other registrations ...

    // HTMX OIDC Login Handlers
    services_htmx.AddOIDCLoginHandlers(builder)
}
```

This registers 10 handlers that serve the login flow:
- Shell (`/oidc-login/`) — full HTML page with HTMX loaded
- Home (`/oidc-login/home`) — email entry + social IdP buttons
- Password (`/oidc-login/password`) — password entry
- Signup (`/oidc-login/signup`) — account creation
- Verify Code (`/oidc-login/verify-code`) — email verification
- Keep Signed In (`/oidc-login/keep-signed-in`) — session persistence prompt
- Forgot Password (`/oidc-login/forgot-password`)
- Reset Password (`/oidc-login/reset-password`)
- Error (`/oidc-login/error`)
- Start Over (`/oidc-login/start-over`)

No additional configuration is needed. The HTMX handlers use the CSS at `/static/go-app/oidc-login/htmx/styles.css`.

## Custom UI

You can implement your own OIDC login UI by:

1. **Not calling** `services_htmx.AddOIDCLoginHandlers(builder)`
2. Registering your own `IHandler` implementations for the `/oidc-login/*` paths
3. Your handlers must implement the same OIDC flow contract (read session cookies, call the identity service, set auth cookies, etc.)

The HTMX handlers in `pkg/services/echo/handlers/htmx/` serve as reference implementations.
