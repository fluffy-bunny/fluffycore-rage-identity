package callback

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_selfoauth2provider "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/selfoauth2provider"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
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
		wellknown_echo.AccountCallbackPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type CallbackRequest struct {
	Code  string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	State string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
}

var validReturnUrlPaths = map[string]bool{
	wellknown_echo.ProfilePath:             true,
	wellknown_echo.HomePath:                true,
	wellknown_echo.PersonalInformationPath: true,
}

func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &CallbackRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("model", model).Msg("model")

	mm, err := s.wellknownCookies.GetInsecureCookie(c, contracts_cookies.LoginRequest)
	if err != nil {
		log.Error().Err(err).Msg("LoginRequest cookie not found")
		return c.Redirect(http.StatusFound, "/error")
	}
	var loginRequest models.LoginGetRequest
	err = models.ConvertFromInterface[models.LoginGetRequest](mm, &loginRequest)
	if err != nil {
		log.Error().Err(err).Msg("Could not convert LoginRequest")
		return c.Redirect(http.StatusFound, "/error")
	}

	isValidPath := func(path string) bool {
		_, ok := validReturnUrlPaths[path]
		return ok
	}
	if !isValidPath(loginRequest.ReturnUrl) {
		loginRequest.ReturnUrl = wellknown_echo.HomePath
	}
	log.Info().Interface("loginRequest", loginRequest).Msg("loginRequest")

	getAccountStateCookieResponse, err := s.wellknownCookies.GetAccountStateCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetAccountStateCookie")
		return c.Redirect(http.StatusFound, "/error")
	}
	state := getAccountStateCookieResponse.AccountStateCookie.State
	nonce := getAccountStateCookieResponse.AccountStateCookie.Nonce

	s.wellknownCookies.DeleteAccountStateCookie(c)

	if model.State != state {
		log.Error().Msg("State did not match")
		return c.Redirect(http.StatusFound, "/error")
	}

	getConfigResponse, err := s.selfOAuth2Provider.GetConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("GetConfig")
		return c.Redirect(http.StatusFound, "/error")
	}
	config := getConfigResponse.Config
	oauth2Token, err := config.Exchange(ctx, model.Code)
	if err != nil {
		log.Error().Err(err).Msg("Exchange")
		return c.Redirect(http.StatusFound, "/error")
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Error().Msg("id_token not found")
		return c.Redirect(http.StatusFound, "/error")
	}
	verifier := getConfigResponse.Verifier
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Error().Err(err).Msg("Verify")
		return c.Redirect(http.StatusFound, "/error")
	}
	if idToken.Nonce != nonce {
		log.Error().Msg("Nonce did not match")
		return c.Redirect(http.StatusFound, "/error")
	}
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		log.Error().Err(err).Msg("Claims")
		return c.Redirect(http.StatusFound, "/error")
	}

	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &models.Identity{
				Subject:       idToken.Subject,
				Email:         claims.Email,
				EmailVerified: claims.EmailVerified,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	return c.Redirect(http.StatusFound, loginRequest.ReturnUrl)
}
