package api_keep_signed_in

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
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

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_KeepSignedIn,
	)
}

const (
	InternalError_KeepSignedIn_001 = "rg-keepsignedin-001"
	InternalError_KeepSignedIn_002 = "rg-keepsignedin-002"
	InternalError_KeepSignedIn_003 = "rg-keepsignedin-003"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateKeepSignedInRequest(model *login_models.KeepSignedInRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	return nil
}

// API KeepSignedIn godoc
// @Summary keep signed in.
// @Description Set or clear SSO cookie for keep me signed in functionality
// @Tags root
// @Accept json
// @Produce json
// @Param		request body		login_models.KeepSignedInRequest	true	"KeepSignedInRequest"
// @Success 200 {object} login_models.KeepSignedInResponse
// @Failure 500 {object} login_models.KeepSignedInErrorResponse
// @Router /api/keep-signed-in [post]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &login_models.KeepSignedInRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: "bind error",
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	if err := s.validateKeepSignedInRequest(model); err != nil {
		log.Error().Err(err).Msg("validateKeepSignedInRequest")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: "invalid request",
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	// Verify that authentication was completed (to prevent direct browsing to this endpoint)
	getAuthCompletedResponse, err := s.WellknownCookies().GetAuthCompletedCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthCompletedCookie - authentication not completed")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: "authentication not completed",
		}
		return c.JSONPretty(http.StatusUnauthorized, response, "  ")
	}
	authCompleted := getAuthCompletedResponse.AuthCompleted

	// Delete the auth completed cookie (one-time use)
	s.WellknownCookies().DeleteAuthCompletedCookie(c)

	// Get the auth cookie to verify user is authenticated
	getAuthCookieResponse, err := s.WellknownCookies().GetAuthCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthCookie")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: InternalError_KeepSignedIn_001,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	authCookie := getAuthCookieResponse.AuthCookie

	// Verify that the auth completed subject matches the auth subject
	if authCompleted.Subject != authCookie.Identity.Subject {
		log.Error().Str("authCompletedSubject", authCompleted.Subject).Str("authSubject", authCookie.Identity.Subject).Msg("Subject mismatch")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: "subject mismatch",
		}
		return c.JSONPretty(http.StatusUnauthorized, response, "  ")
	}

	// Handle SSO cookie based on keepSignedIn flag
	if model.KeepSignedIn {
		// Set SSO cookie with configured duration - copy all auth context
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
			response := &login_models.KeepSignedInErrorResponse{
				Reason: InternalError_KeepSignedIn_002,
			}
			return c.JSONPretty(http.StatusInternalServerError, response, "  ")
		}
		log.Info().Str("subject", authCookie.Identity.Subject).Msg("SSO cookie set for keep signed in")
	} else {
		// Delete SSO cookie if user opted out
		s.WellknownCookies().DeleteSSOCookie(c)
		log.Info().Str("subject", authCookie.Identity.Subject).Msg("SSO cookie deleted")
	}

	// Handle DoNotShowAgain preference
	if model.DoNotShowAgain {
		err = s.WellknownCookies().SetKeepSigninPreferencesCookie(c,
			&contracts_cookies.SetKeepSigninPreferencesCookieRequest{
				Subject: authCookie.Identity.Subject,
				KeepSigninPreferencesCookie: &contracts_cookies.KeepSigninPreferencesCookie{
					PreferenceValue: true,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetKeepSigninPreferencesCookie")
		}
		log.Info().Str("subject", authCookie.Identity.Subject).Msg("KeepSigninPreferences cookie set")
	} else {
		// Delete the cookie if user wants to see the page again
		s.WellknownCookies().DeleteKeepSigninPreferencesCookie(c,
			&contracts_cookies.DeleteKeepSigninPreferencesCookieRequest{
				Subject: authCookie.Identity.Subject,
			})
	}

	// Get session and authorization request
	session, err := s.getSession()
	if err != nil {
		log.Error().Err(err).Msg("getSession")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: InternalError_KeepSignedIn_003,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	sessionRequest, err := session.Get("request")
	if err != nil {
		log.Error().Err(err).Msg("session.Get")
		response := &login_models.KeepSignedInErrorResponse{
			Reason: InternalError_KeepSignedIn_003,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	rootPath := echo_utils.GetMyRootPath(c)

	// Process the final authentication state and redirect
	finalStateResponse, err := s.ProcessFinalAuthenticationState(ctx, c,
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
		response := &login_models.KeepSignedInErrorResponse{
			Reason: InternalError_KeepSignedIn_003,
		}
		return c.JSONPretty(http.StatusInternalServerError, response, "  ")
	}

	response := &login_models.KeepSignedInResponse{
		Directive: login_models.DIRECTIVE_Redirect,
		DirectiveRedirect: &login_models.DirectiveRedirect{
			RedirectURI: finalStateResponse.RedirectURI,
		},
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
