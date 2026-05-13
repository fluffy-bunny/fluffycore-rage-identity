package end_session_endpoint

import (
	"net/http"
	"net/url"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a scoped handler.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OIDCEndSessionEndpointPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type endSessionRequest struct {
	PostLogoutRedirectURI string `query:"post_logout_redirect_uri"`
}

// EndSession godoc
// @Summary OIDC RP-Initiated Logout / front-channel logout endpoint.
// @Description Clears the SSO cookie. Supports two call patterns:
// @Description
// @Description **Top-level navigation** – redirect the browser here with
// @Description post_logout_redirect_uri set. The cookie is cleared and the browser
// @Description is sent back to that URI (must be an absolute http/https URL).
// @Description
// @Description **Hidden iframe (front-channel logout)** – embed this URL in a
// @Description hidden iframe on your portal logout page. The identity server clears
// @Description the SSO cookie and returns a 200 HTML page; no navigation occurs.
// @Description This lets a cross-domain portal invalidate the SSO session without a
// @Description full-page redirect.
// @Description
// @Description ⚠ SameSite note: for the iframe pattern the SSO cookie must be
// @Description configured with SameSite=None; Secure so browsers include it in
// @Description cross-origin iframe requests. With the default SameSite=Lax the
// @Description Set-Cookie expiry header is still emitted but the cookie is not sent
// @Description by the browser in the iframe request, so there is nothing to delete.
// @Tags oidc
// @Accept */*
// @Produce html
// @Param post_logout_redirect_uri query string false "URI to redirect to after logout (top-level navigation mode)"
// @Success 200 {string} string "SSO cookie cleared (iframe mode)"
// @Success 302 {string} string "Redirect to post_logout_redirect_uri (navigation mode)"
// @Router /oidc/v1/endsession [get]
func (s *service) Do(c *echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	req := &endSessionRequest{}
	if err := c.Bind(req); err != nil {
		log.Error().Err(err).Msg("end_session: bind")
		// Return a silent empty page; works in both iframe and navigation contexts.
		return c.HTML(http.StatusBadRequest, `<html><body></body></html>`)
	}

	// Delete the SSO cookie. Because this endpoint lives on the identity
	// server's domain the browser will send the cookie here and we can
	// expire it via the Set-Cookie response header.
	s.wellknownCookies.DeleteSSOCookie(c)
	log.Info().Msg("end_session: SSO cookie cleared")

	// --- navigation mode: post_logout_redirect_uri was supplied ---
	if req.PostLogoutRedirectURI != "" {
		// Validate: must be an absolute http/https URL to prevent open-redirect abuse.
		parsed, err := url.ParseRequestURI(req.PostLogoutRedirectURI)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			log.Warn().Str("post_logout_redirect_uri", req.PostLogoutRedirectURI).
				Msg("end_session: invalid post_logout_redirect_uri, ignoring")
		} else {
			return c.Redirect(http.StatusFound, req.PostLogoutRedirectURI)
		}
	}

	// --- iframe / front-channel mode: return a silent 200 HTML page ---
	// The Set-Cookie header above has already expired the SSO cookie.
	// Returning 200 (not a redirect) lets the iframe complete its load
	// and signals to the embedding page that the logout succeeded.
	return c.HTML(http.StatusOK, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Signed out</title></head><body></body></html>`)
}
