# OIDC Login UI — Swapping Between HTMX and WASM

The OIDC login UI is pluggable via dependency injection. By default, the login flow is served by HTMX-based server-rendered handlers. You can swap this out for the WASM SPA implementation (or any custom UI) by changing a single registration call in your `startup.go`.

## Architecture

The `/oidc-login` POST handler (`pkg/services/echo/handlers/oidclogin/oidc-login.go`) acts as the entry point for all OIDC authorization flows. When an authorization request comes in, it checks `config.OIDCUIConfig.URIEntryPath`:

- If `URIEntryPath` matches the default `/oidc-login`, the handler processes the request directly.
- If `URIEntryPath` differs (e.g., `/oidc-login/`), it redirects to that path — where your chosen UI implementation takes over.

Both HTMX and WASM implementations register handlers at the `/oidc-login/` path prefix.

## Option 1: HTMX (Default)

The default uses server-rendered HTMX handlers via gomponents. This is registered with a single call in your `MyConfigServices` function:

```go
import (
    services_htmx "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/htmx"
)

func (s *startup) MyConfigServices(ctx context.Context, config *rage_contracts_config.Config, builder di.ContainerBuilder) {
    // ... other registrations ...

    // HTMX OIDC Login Handlers (default UI implementation)
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

No additional configuration is needed. The HTMX handlers use the CSS from the WASM build at `/static/go-app/oidc-login/static_output/web/styles.css`.

## Option 2: WASM SPA

To use the WASM go-app SPA instead:

### Step 1: Remove the HTMX registration

In your `MyConfigServices`, **remove or comment out** the HTMX call:

```go
// services_htmx.AddOIDCLoginHandlers(builder)  // <-- remove this
```

### Step 2: Register the CacheBustingHTMLConfig for WASM

Add a `CacheBustingHTMLConfig` registration that serves the WASM SPA. This follows the same pattern used for the management app:

```go
import (
    rage_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
    services_handlers_cache_busting_static_html "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/cache_busting_static_html"
)

func (s *startup) MyConfigServices(ctx context.Context, config *rage_contracts_config.Config, builder di.ContainerBuilder) {
    // ... other registrations ...

    guid := xid.New().String()
    if example_version.Version() != "dev-build" {
        guid = example_version.Version()
    }
    config.CacheBustVersion = guid

    // WASM OIDC Login SPA
    oidcLoginCacheBustingHTMLConfig := &rage_contracts_config.CacheBustingHTMLConfig{
        Version:    guid,
        FilePath:   "./static/go-app/oidc-login/static_output/index_template.html",
        StaticPath: "./static/go-app/oidc-login/static_output/",
        EchoPath:   "/oidc-login/*",
        RootPath:   "/oidc-login/",
        ReplaceParams: []*rage_contracts_config.KeyValuePair{
            {
                Key:   "{title}",
                Value: s.config.OIDCLoginAppConfig.BannerBranding.Title,
            },
            {
                Key:   "{version}",
                Value: guid,
            },
        },
        RoutePatterns: []*rage_contracts_config.RoutePattern{
            {
                Pattern: "/web/app.wasm",
                Handler: func(c *echo.Context, filePath string) (bool, error) {
                    fileInfo, err := os.Stat(filePath)
                    if err != nil {
                        return false, err
                    }
                    c.Response().Header().Set("Content-Type", "application/wasm")
                    c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
                    return true, c.File(filePath)
                },
            },
            {
                Pattern: "web/app.json",
                Handler: func(c *echo.Context, filePath string) (bool, error) {
                    jsonB, err := json.Marshal(s.config.OIDCLoginAppConfig)
                    if err != nil {
                        return false, err
                    }
                    version := c.QueryParam("v")
                    modifiedContent := strings.ReplaceAll(string(jsonB), "{version}", version)
                    return true, c.JSONBlob(http.StatusOK, []byte(modifiedContent))
                },
            },
        },
    }
    services_handlers_cache_busting_static_html.AddScopedIHandler(builder, oidcLoginCacheBustingHTMLConfig)
}
```

### Step 3: Environment / Config

Ensure these environment variables or config values are set (see `launch.json` or `docker-compose.yml` for examples):

| Variable | Description | Example |
|---|---|---|
| `RAGE_oidcUIConfig__uriEntryPath` | Where `/oidc-login` redirects to | `/oidc-login/` |
| `RAGE_oidcUIConfig__cacheBustingConfig__filePath` | Path to `index_template.html` | `./static/go-app/oidc-login/index_template.html` |
| `RAGE_oidcUIConfig__cacheBustingConfig__staticPath` | Static files directory | `./static/go-app/oidc-login/` |
| `oidcLoginAppConfig__rageBaseUrl` | Base URL for the rage identity server | `http://localhost:9044` |

## How DI Swapping Works

Both the HTMX handlers and the WASM `CacheBustingHTMLConfig` handler register against the same `/oidc-login/*` route paths using `AddScopedIHandleWithMetadata`. The DI container uses the path as a lookup key. **The last registration for a given path wins.** This means:

- If you call `services_htmx.AddOIDCLoginHandlers(builder)` — the HTMX handlers serve `/oidc-login/*`.
- If you call `services_handlers_cache_busting_static_html.AddScopedIHandler(builder, oidcLoginCacheBustingHTMLConfig)` — the WASM SPA serves `/oidc-login/*`.
- If you call **both**, the second registration overwrites the first.

Pick one. Don't register both unless you understand the override behavior.

## Custom UI

You can implement your own OIDC login UI by:

1. **Not calling** `services_htmx.AddOIDCLoginHandlers(builder)`
2. Registering your own `IHandler` implementations for the `/oidc-login/*` paths
3. Your handlers must implement the same OIDC flow contract (read session cookies, call the identity service, set auth cookies, etc.)

The HTMX handlers in `pkg/services/echo/handlers/htmx/` serve as reference implementations.
