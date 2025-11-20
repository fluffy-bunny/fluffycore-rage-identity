package oidcloginpasskey

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		oidcSession      contracts_oidc_session.IOIDCSession
		signinResponse   *contracts_cookies.GetSigninUserNameCookieResponse
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_OIDCLoginPasskey_001 = "rg-oidclogin-passkey-001"
	InternalError_OIDCLoginPasskey_002 = "rg-oidclogin-passkey-002"
	InternalError_OIDCLoginPasskey_003 = "rg-oidclogin-passkey-003"
	InternalError_OIDCLoginPasskey_004 = "rg-oidclogin-passkey-004"
	InternalError_OIDCLoginPasskey_005 = "rg-oidclogin-passkey-005"
	InternalError_OIDCLoginPasskey_006 = "rg-oidclogin-passkey-006"
	InternalError_OIDCLoginPasskey_007 = "rg-oidclogin-passkey-007"
	InternalError_OIDCLoginPasskey_008 = "rg-oidclogin-passkey-008"
	InternalError_OIDCLoginPasskey_009 = "rg-oidclogin-passkey-009"
	InternalError_OIDCLoginPasskey_010 = "rg-oidclogin-passkey-010"
	InternalError_OIDCLoginPasskey_011 = "rg-oidclogin-passkey-011"

	InternalError_OIDCLoginPasskey_099 = "rg-oidclogin-passkey-099"
)

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
			contracts_handler.GET,
		},
		wellknown_echo.OIDCLoginPasskeyPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	log.Debug().Msg("OIDCLoginPasskey")

	return s.Render(c, http.StatusOK,
		"oidc/oidcloginpasskey/index",
		map[string]interface{}{
			"returnFailedUrl": wellknown_echo.OIDCLoginPasswordPath,
		})
}

func (s *service) Do(c echo.Context) error {

	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	signinResponse, err := s.wellknownCookies.GetSigninUserNameCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetSigninUserNameCookie")
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginPasskey_004)
	}
	s.signinResponse = signinResponse

	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoPost(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
