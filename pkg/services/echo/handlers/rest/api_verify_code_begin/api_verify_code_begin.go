package api_verify_code_begin

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	verify_code "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_code"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config

		oidcSession      contracts_oidc_session.IOIDCSession
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	oidcSession contracts_oidc_session.IOIDCSession,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),

		config:           config,
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
		wellknown_echo.API_VerifyCodeBegin,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API VerifyCodeBegin godoc
// @Summary get the login manifest.
// @Description Validates if we can do a verification code flow
// @Tags root
// @Produce json
// @Success 200 {object} verify_code.VerifyCodeBeginResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/verify-code-begin [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	response := &verify_code.VerifyCodeBeginResponse{
		Valid: false,
	}
	session, err := s.getSession()
	if err != nil {
		log.Error().Err(err).Msg("getSession")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	clearLandingPage := func() {
		landingPageI, err := session.Get("landingPage")
		if err == nil && landingPageI != nil {
			session.Set("landingPage", nil)
			session.Save()
		}
	}
	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		// lets clear out the landing page cookie
		clearLandingPage()
		// not a valid verification code cookie
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	response.Email = getVerificationCodeCookieResponse.VerificationCode.Email
	if s.config.SystemConfig.DeveloperMode {
		response.Code = getVerificationCodeCookieResponse.VerificationCode.PlainCode
	}
	response.Valid = true
	clearLandingPage()
	return c.JSONPretty(http.StatusOK, response, "  ")
}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}
