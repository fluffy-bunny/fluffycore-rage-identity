package shell

import (
	"fmt"
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.HTMXSupportPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	cp := s.ClaimsPrincipal()
	userEmail := ""
	userName := ""
	if cp != nil {
		emailClaims := cp.GetClaimsByType("email")
		if len(emailClaims) > 0 {
			userEmail = emailClaims[0].Value
		}
		nameClaims := cp.GetClaimsByType("name")
		if len(nameClaims) > 0 {
			userName = nameClaims[0].Value
		}
	}

	deepLink := wellknown_echo.HTMXSupportAuditsPath
	if redirect := c.QueryParam("redirect"); redirect != "" && strings.HasPrefix(redirect, wellknown_echo.SupportPath) {
		deepLink = redirect
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Support Portal</title>
  <script src="https://unpkg.com/htmx.org@1.9.12"></script>
  <style>
    body { font-family: Segoe UI, sans-serif; margin: 0; background: #0f172a; color: #e2e8f0; }
    header { padding: 16px 24px; background: #111827; border-bottom: 1px solid #334155; }
    .meta { color: #94a3b8; font-size: 12px; }
    .layout { display: grid; grid-template-columns: 240px 1fr; min-height: calc(100vh - 68px); }
    nav { border-right: 1px solid #334155; padding: 16px; background: #111827; }
    nav a { display: block; color: #cbd5e1; text-decoration: none; padding: 10px 8px; border-radius: 6px; margin-bottom: 8px; }
    nav a:hover { background: #1f2937; }
    main { padding: 20px; }
    .card { background: #111827; border: 1px solid #334155; border-radius: 8px; padding: 16px; }
  </style>
</head>
<body>
  <header>
    <div><strong>Backend Support Portal</strong></div>
    <div class="meta">Support Admin: %s (%s)</div>
  </header>
  <div class="layout">
    <nav>
      <a href="%s" hx-get="%s" hx-target="#content" hx-push-url="true">Audit Logs</a>
    </nav>
    <main>
			<div class="card" style="margin-bottom:16px;">
				<h3 style="margin-top:0;">Boundary</h3>
				<p style="margin:0; color:#cbd5e1;">
					This portal is scoped to identity-engine operations (audit and auth context tooling).
					User directory management is owned by integrator systems.
				</p>
			</div>
      <div id="content" class="card" hx-get="%s" hx-trigger="load" hx-push-url="true">Loading...</div>
    </main>
  </div>
</body>
</html>`, userName, userEmail, wellknown_echo.HTMXSupportAuditsPath, wellknown_echo.HTMXSupportAuditsPath, deepLink)

	return c.HTML(http.StatusOK, html)
}
