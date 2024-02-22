package about

import (
	"net/http"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_selfoauth2provider "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/selfoauth2provider"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		wellknownCookies   contracts_cookies.IWellknownCookies
		selfOAuth2Provider contracts_selfoauth2provider.ISelfOAuth2Provider
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	selfOAuth2Provider contracts_selfoauth2provider.ISelfOAuth2Provider,
) (*service, error) {
	return &service{
		BaseHandler:        services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies:   wellknownCookies,
		selfOAuth2Provider: selfOAuth2Provider,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.LoginPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	state := xid.New().String()
	nonce := xid.New().String()
	err := s.wellknownCookies.SetAccountStateCookie(c, &contracts_cookies.SetAccountStateCookieRequest{
		AccountStateCookie: &contracts_cookies.AccountStateCookie{
			State: state,
			Nonce: nonce,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAccountStateCookie")
		return c.Redirect(http.StatusFound, "/error")
	}

	getConfigResponse, err := s.selfOAuth2Provider.GetConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("GetConfig")
		return c.Redirect(http.StatusFound, "/error")
	}
	config := getConfigResponse.Config
	authUrl := config.AuthCodeURL(state, oidc.Nonce(nonce))
	log.Info().Str("authUrl", authUrl).Msg("authUrl")

	return c.Redirect(http.StatusFound, authUrl)

}
