package shell

import (
	"context"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	oidc_login_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/contracts/config"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/components"
	example_version "github.com/fluffy-bunny/fluffycore-rage-identity/example/version"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config    *contracts_config.Config
		appConfig *oidc_login_config.AppConfig
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	appConfig *oidc_login_config.AppConfig,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
		appConfig:   appConfig,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.HTMXOIDCLoginPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Validate that the session's authorization request state still exists in the
	// backing store. After a server restart with an in-memory store, the cookie
	// session survives but the server-side state is gone. Detect this early and
	// redirect back to the client for a fresh OIDC flow.
	if redirect, ok := s.validateSessionState(ctx, c, &log); ok {
		return redirect
	}

	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	rc.CacheBustVersion = s.config.CacheBustVersion

	// Determine initial page based on session landing page directive
	initialPage := rc.Paths.HTMXHome
	session, err := s.OIDCSession().GetSession()
	if err == nil {
		landingPageI, err := session.Get("landingPage")
		if err == nil && landingPageI != nil {
			if lp, ok := landingPageI.(*models_api_manifest.LandingPage); ok && lp != nil {
				switch lp.Page {
				case models_api_manifest.PageVerifyCode:
					initialPage = rc.Paths.HTMXVerifyCode
				case models_api_manifest.PageKeepSignedIn:
					initialPage = rc.Paths.HTMXKeepSignedIn
				case models_api_manifest.PagePasswordEntry:
					initialPage = rc.Paths.HTMXPassword
				}
				if initialPage != rc.Paths.HTMXHome {
					log.Info().Str("landingPage", string(lp.Page)).Str("initialPage", initialPage).Msg("shell: routing to session landing page")
				}
				// Consume the landing page so it doesn't fire again on refresh
				session.Set("landingPage", nil)
				session.Save()
			}
		}
	}

	// Derive brand title from config
	brandTitle := "Identity"
	if s.appConfig != nil && s.appConfig.BannerBranding.Title != "" {
		brandTitle = s.appConfig.BannerBranding.Title
	}

	return components.RenderNode(c, http.StatusOK, components.ShellPage(components.ShellData{
		RenderContext: rc,
		BrandTitle:    brandTitle,
		InitialPage:   initialPage,
		AppVersion:    example_version.Version(),
		ShowVersion:   s.appConfig.BannerBranding.ShowBannerVersion,
	}))
}

// validateSessionState checks whether the session's authorization request state
// still exists in the backing store. Returns (redirect response, true) if the
// session is stale and the caller should return the redirect. Returns (nil, false)
// if everything is fine and the caller should continue normally.
// Wrapped in recover() so that deserialization issues with stale cookies never
// surface as 500 Internal Server Error.
func (s *service) validateSessionState(ctx context.Context, c *echo.Context, log *zerolog.Logger) (err error, stale bool) {
	defer func() {
		if r := recover(); r != nil {
			fallbackURL := s.GetFallbackURL()
			log.Error().Interface("panic", r).Str("fallbackURL", fallbackURL).Msg("shell: panic during session state validation, clearing session and redirecting to fallback URL")
			// Best-effort clear of the session
			if session, sErr := s.OIDCSession().GetSession(); sErr == nil {
				session.Set("request", nil)
				session.Set("session_id", nil)
				session.Set("landing_page", nil)
				session.Set("landingPage", nil)
				session.Save()
			}
			err = c.Redirect(http.StatusFound, fallbackURL)
			stale = true
		}
	}()

	session, sErr := s.OIDCSession().GetSession()
	if sErr != nil {
		return nil, false
	}
	requestI, gErr := session.Get("request")
	if gErr != nil || requestI == nil {
		return nil, false
	}
	authReq, ok := requestI.(*proto_oidc_models.AuthorizationRequest)
	if !ok || authReq == nil {
		return nil, false
	}

	_, gErr = s.AuthorizationRequestStateStore().GetAuthorizationRequestState(ctx,
		&proto_oidc_flows.GetAuthorizationRequestStateRequest{
			State: authReq.State,
		})
	if gErr != nil {
		clientReturnURL := s.GetClientReturnURL(ctx, authReq.ClientId, authReq.RedirectUri)
		log.Warn().Err(gErr).Str("state", authReq.State).Str("clientReturnURL", clientReturnURL).
			Msg("shell: authorization request state not found in store (server may have restarted), redirecting to client for fresh OIDC flow")
		// Clear the stale session so the next authorization flow starts clean
		session.Set("request", nil)
		session.Set("session_id", nil)
		session.Set("landing_page", nil)
		session.Set("landingPage", nil)
		session.Save()
		return c.Redirect(http.StatusFound, clientReturnURL), true
	}
	return nil, false
}
