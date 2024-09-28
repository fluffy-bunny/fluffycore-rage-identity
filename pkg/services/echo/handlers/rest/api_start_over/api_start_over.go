package api_start_over

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config

		webAuthNConfig   *contracts_webauthn.WebAuthNConfig
		oidcSession      contracts_oidc_session.IOIDCSession
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	webAuthNConfig *contracts_webauthn.WebAuthNConfig,
	oidcSession contracts_oidc_session.IOIDCSession,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),

		config:           config,
		webAuthNConfig:   webAuthNConfig,
		oidcSession:      oidcSession,
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_StartOver,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API StartOver godoc
// @Summary start over and return the original manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Produce json
// @Success 200 {object} manifest.Manifest
// @Router /api/start-over [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()

	idps, err := s.GetIDPs(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	response := &manifest.Manifest{
		DevelopmentMode: s.config.SystemConfig.DeveloperMode,
	}
	for _, idp := range idps {
		if idp.Enabled && !idp.Hidden {
			response.SocialIdps = append(response.SocialIdps, manifest.IDP{
				Slug: idp.Slug,
			})
		}
	}
	response.PasskeyEnabled = false
	if s.webAuthNConfig != nil {
		response.PasskeyEnabled = s.webAuthNConfig.Enabled
	}
	// we may have a session in flight and got redirect back here.
	// we may have an external OIDC callback that requires a verify code to continue.
	session, err := s.getSession()
	if err == nil {
		sessionRequest, err := session.Get("request")
		if err == nil {
			authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)
			// we
			if authorizationRequest != nil {
				landingPageI, err := session.Get("landingPage")
				if err == nil && landingPageI != nil {

					// get rid of it.
					session.Set("landingPage", nil)
					session.Save()
				}
			}
		}

	}
	return c.JSONPretty(http.StatusOK, response, "  ")
}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}
