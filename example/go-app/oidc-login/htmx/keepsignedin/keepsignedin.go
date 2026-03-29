package keepsignedin

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/htmx/components"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config      *contracts_config.Config
		oidcSession contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
		oidcSession: oidcSession,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXKeepSignedInPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type KeepSignedInPostRequest struct {
	KeepSignedIn string `param:"keepSignedIn" query:"keepSignedIn" form:"keepSignedIn" json:"keepSignedIn" xml:"keepSignedIn"`
}

func (s *service) Do(c *echo.Context) error {
	r := c.Request()
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	return c.NoContent(http.StatusNotFound)
}

func (s *service) renderError(c *echo.Context, errorCode, errorMessage string) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.ErrorPartial(components.ErrorData{
		RenderContext: rc,
		ErrorCode:     errorCode,
		ErrorMessage:  errorMessage,
	}))
}

func (s *service) renderErrorWithReturn(c *echo.Context, errorCode, errorMessage, returnURL string) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.ErrorPartial(components.ErrorData{
		RenderContext: rc,
		ErrorCode:     errorCode,
		ErrorMessage:  errorMessage,
		ReturnURL:     returnURL,
	}))
}

func (s *service) DoGet(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.KeepSignedInPartial(components.KeepSignedInData{
		RenderContext: rc,
	}))
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &KeepSignedInPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderError(c, "htmx-ksi-099", "Invalid request")
	}

	keepSignedIn := model.KeepSignedIn == "true"

	// Verify authentication was completed
	getAuthCompletedResponse, err := s.WellknownCookies().GetAuthCompletedCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthCompletedCookie")
		return s.renderError(c, "htmx-ksi-001", "Authentication not completed")
	}
	authCompleted := getAuthCompletedResponse.AuthCompleted
	s.WellknownCookies().DeleteAuthCompletedCookie(c)

	// Get the auth cookie
	getAuthCookieResponse, err := s.WellknownCookies().GetAuthCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthCookie")
		return s.renderError(c, "htmx-ksi-002", "Auth cookie not found")
	}
	authCookie := getAuthCookieResponse.AuthCookie

	if authCompleted.Subject != authCookie.Identity.Subject {
		return s.renderError(c, "htmx-ksi-003", "Subject mismatch")
	}

	if keepSignedIn {
		err = s.WellknownCookies().SetSSOCookie(c,
			&contracts_cookies.SetSSOCookieRequest{
				SSOCookie: &contracts_cookies.SSOCookie{
					Identity: authCookie.Identity,
					Acr:      authCookie.Acr,
					Amr:      authCookie.Amr,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetSSOCookie")
			return s.renderError(c, "htmx-ksi-004", err.Error())
		}
	} else {
		s.WellknownCookies().DeleteSSOCookie(c)
	}

	// Get session and complete OAuth flow
	session, err := s.oidcSession.GetSession()
	if err != nil {
		log.Error().Err(err).Msg("GetSession")
		return s.renderError(c, "htmx-ksi-005", "Session error")
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		log.Error().Err(err).Msg("session.Get request")
		return s.renderError(c, "htmx-ksi-006", "Session error")
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)
	rootPath := echo_utils.GetMyRootPath(c)

	result, err := s.ProcessFinalAuthenticationState(ctx, c,
		&services_echo_handlers_base.ProcessFinalAuthenticationStateRequest{
			AuthorizationRequest: authorizationRequest,
			Identity: &proto_oidc_models.OIDCIdentity{
				Subject:       authCookie.Identity.Subject,
				Email:         authCookie.Identity.Email,
				EmailVerified: authCookie.Identity.EmailVerified,
				IdpSlug:       authCookie.Identity.IdpSlug,
				Acr:           authCookie.Acr,
				Amr:           authCookie.Amr,
			},
			RootPath: rootPath,
		})
	if err != nil {
		log.Error().Err(err).Msg("ProcessFinalAuthenticationState")
		returnURL := s.GetClientReturnURL(ctx, authorizationRequest.ClientId, authorizationRequest.RedirectUri)
		return s.renderErrorWithReturn(c, "htmx-ksi-007", err.Error(), returnURL)
	}

	c.Response().Header().Set("HX-Redirect", result.RedirectURI)
	return c.NoContent(http.StatusOK)
}
