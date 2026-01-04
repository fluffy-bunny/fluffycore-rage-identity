package api_verify_code

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
		oidcSession      contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

// 	API_VerifyCode             = "/api/verify-code"

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_VerifyCode,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateVerifyCodeRequest(model *login_models.VerifyCodeRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Code) {
		return status.Error(codes.InvalidArgument, "model.Code is nil")
	}

	return nil
}

// API VerifyCode godoc
// @Summary verify code.
// @Description verify code
// @Tags root
// @Accept json
// @Produce json
// @Param		request body		login_models.VerifyCodeRequest	true	"VerifyCodeRequest"
// @Success 200 {object} login_models.VerifyCodeResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Router /api/verify-code [post]
func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &login_models.VerifyCodeRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if err := s.validateVerifyCodeRequest(model); err != nil {
		log.Error().Err(err).Msg("validateVerifyCodeRequest")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}

	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetVerificationCodeCookie")
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	verificationCode := getVerificationCodeCookieResponse.VerificationCode
	verificationCodeHashed := verificationCode.CodeHash

	err = echo_utils.VerifyVerificationCode(ctx, s.passwordHasher, model.Code, verificationCodeHashed)
	if err != nil {
		err := status.Error(codes.NotFound, "code does not match")
		return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	userService := s.RageUserService()
	getRageUserResponse, err := userService.GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: verificationCode.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("GetRageUser")
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			err := status.Error(codes.NotFound, "User not found")
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	rageUser := getRageUserResponse.User
	_, err = userService.UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: rageUser.RootIdentity.Subject,
				EmailVerified: &wrapperspb.BoolValue{
					Value: true,
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	// one time only
	s.wellknownCookies.DeleteVerificationCodeCookie(c)

	switch verificationCode.VerifyCodePurpose {
	case contracts_cookies.VerifyCode_EmailVerification:
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_LoginPhaseOne_DisplayPhaseOnePage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	case contracts_cookies.VerifyCode_PasswordReset:
		err = s.wellknownCookies.SetPasswordResetCookie(c,
			&contracts_cookies.SetPasswordResetCookieRequest{
				PasswordReset: &contracts_cookies.PasswordReset{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetPasswordResetCookie")
			return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_PasswordReset_DisplayPasswordResetPage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")
	case contracts_cookies.VerifyCode_Challenge:
		// Set auth cookie with user identity
		err = s.wellknownCookies.SetAuthCookie(c,
			&contracts_cookies.SetAuthCookieRequest{
				AuthCookie: &contracts_cookies.AuthCookie{
					Identity: &proto_oidc_models.Identity{
						Subject:       rageUser.RootIdentity.Subject,
						Email:         rageUser.RootIdentity.Email,
						EmailVerified: rageUser.RootIdentity.EmailVerified,
					},
					Acr: []string{
						models.ACRPassword,
						models.ACRIdpRoot,
					},
					Amr: []string{
						models.AMRPassword,
						models.AMRIdp,
						models.AMRMFA,
						models.AMREmailCode,
					},
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetAuthCookie")
			return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}

		// Set auth completed cookie to track that authentication was successful
		err = s.wellknownCookies.SetAuthCompletedCookie(c,
			&contracts_cookies.SetAuthCompletedCookieRequest{
				AuthCompleted: &contracts_cookies.AuthCompleted{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetAuthCompletedCookie")
			return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}

		// Check if user has keep-signed-in preferences set
		getKeepSigninPreferencesCookieResponse, err := s.wellknownCookies.GetKeepSigninPreferencesCookie(c, &contracts_cookies.GetKeepSigninPreferencesCookieRequest{
			Subject: verificationCode.Subject,
		})
		if err != nil {
			log.Error().Err(err).Msg("GetKeepSigninPreferencesCookie")
			// If we can't read the cookie, continue with default flow
		}

		// If user has opted to skip keep-signed-in page, complete OAuth flow directly
		if getKeepSigninPreferencesCookieResponse != nil && getKeepSigninPreferencesCookieResponse.KeepSigninPreferencesCookie != nil {
			log.Info().
				Str("subject", verificationCode.Subject).
				Msg("Skipping keep-signed-in page due to KeepSigninPreferences cookie in verify code")

			// Set SSO cookie since we're auto-keeping them signed in
			err = s.wellknownCookies.SetSSOCookie(c,
				&contracts_cookies.SetSSOCookieRequest{
					SSOCookie: &contracts_cookies.SSOCookie{
						Identity: &proto_oidc_models.Identity{
							Subject:       rageUser.RootIdentity.Subject,
							Email:         rageUser.RootIdentity.Email,
							EmailVerified: rageUser.RootIdentity.EmailVerified,
							IdpSlug:       models.RootIdp,
						},
						Acr: []string{
							models.ACRPassword,
							models.ACRIdpRoot,
						},
						Amr: []string{
							models.AMRPassword,
							models.AMRIdp,
							models.AMRMFA,
							models.AMREmailCode,
						},
					},
				})
			if err != nil {
				log.Error().Err(err).Msg("SetSSOCookie during verify code auto-skip")
				return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
			}
			log.Info().
				Str("subject", verificationCode.Subject).
				Msg("Set SSO cookie during verify code auto-skip")

			// Delete the AuthCompleted cookie (one-time use)
			s.wellknownCookies.DeleteAuthCompletedCookie(c)

			// Get session and authorization request
			session, err := s.oidcSession.GetSession()
			if err != nil {
				log.Error().Err(err).Msg("GetSession during verify code auto-skip")
				return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
			}

			sessionRequest, err := session.Get("request")
			if err != nil {
				log.Error().Err(err).Msg("session.Get during verify code auto-skip")
				return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
			}
			authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

			rootPath := echo_utils.GetMyRootPath(c)

			// Complete the OAuth flow directly
			finalStateResponse, err := s.ProcessFinalAuthenticationState(ctx, c,
				&services_echo_handlers_base.ProcessFinalAuthenticationStateRequest{
					AuthorizationRequest: authorizationRequest,
					Identity: &proto_oidc_models.OIDCIdentity{
						Subject:       rageUser.RootIdentity.Subject,
						Email:         rageUser.RootIdentity.Email,
						EmailVerified: rageUser.RootIdentity.EmailVerified,
						Acr: []string{
							models.ACRPassword,
							models.ACRIdpRoot,
						},
						Amr: []string{
							models.AMRPassword,
							models.AMRIdp,
							models.AMRMFA,
							models.AMREmailCode,
						},
					},
					RootPath: rootPath,
				})
			if err != nil {
				log.Error().Err(err).Msg("ProcessFinalAuthenticationState during verify code auto-skip")
				return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
			}

			response := &login_models.VerifyCodeResponse{
				Directive: login_models.DIRECTIVE_Redirect,
				DirectiveRedirect: &login_models.DirectiveRedirect{
					RedirectURI: finalStateResponse.RedirectURI,
				},
			}
			return c.JSONPretty(http.StatusOK, response, "  ")
		}

		response := &login_models.VerifyCodeResponse{
			Directive: login_models.DIRECTIVE_KeepSignedIn_DisplayKeepSignedInPage,
		}
		return c.JSONPretty(http.StatusOK, response, "  ")

	}
	return c.JSONPretty(http.StatusInternalServerError, "Unknown VerifyCodePurpose", "  ")

}
