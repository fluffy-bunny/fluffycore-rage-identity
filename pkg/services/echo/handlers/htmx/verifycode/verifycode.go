package verifycode

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
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

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXVerifyCodePath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type VerifyCodePostRequest struct {
	Code      string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
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

func (s *service) renderVerifyCode(c *echo.Context, code int, errors []string, email, directive, codeVal string) error {
	return s.Render(c, code, "oidc/htmx/_partials/verify-code", map[string]interface{}{
		"errors":    errors,
		"email":     email,
		"directive": directive,
		"code":      codeVal,
	})
}

func (s *service) renderError(c *echo.Context, errorCode, errorMessage string) error {
	return s.Render(c, http.StatusOK, "oidc/htmx/_partials/error", map[string]interface{}{
		"errorCode":    errorCode,
		"errorMessage": errorMessage,
	})
}

func (s *service) DoGet(c *echo.Context) error {
	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		return s.renderError(c, "htmx-verify-001", "No verification code in progress. Please start over.")
	}
	vc := getVerificationCodeCookieResponse.VerificationCode
	code := ""
	if s.config.SystemConfig.DeveloperMode {
		code = vc.PlainCode
	}
	return s.renderVerifyCode(c, http.StatusOK, nil, vc.Email, string(vc.VerifyCodePurpose), code)
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &VerifyCodePostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderError(c, "htmx-verify-099", "Invalid request")
	}

	if fluffycore_utils.IsEmptyOrNil(model.Code) {
		return s.renderVerifyCode(c, http.StatusOK, []string{"Code is required"}, model.Email, model.Directive, "")
	}

	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetVerificationCodeCookie")
		// Session expired, start over
		return s.Render(c, http.StatusOK, "oidc/htmx/_partials/home", map[string]interface{}{
			"errors": []string{"Verification session expired. Please start over."},
			"email":  model.Email,
		})
	}
	verificationCode := getVerificationCodeCookieResponse.VerificationCode

	err = echo_utils.VerifyVerificationCode(ctx, s.passwordHasher, model.Code, verificationCode.CodeHash)
	if err != nil {
		return s.renderVerifyCode(c, http.StatusOK, []string{"Invalid verification code"}, model.Email, model.Directive, "")
	}

	// Get user
	userService := s.RageUserService()
	getRageUserResponse, err := userService.GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: verificationCode.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("GetRageUser")
		return s.renderError(c, "htmx-verify-002", "User not found")
	}
	rageUser := getRageUserResponse.User

	// Mark email as verified
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
		return s.renderError(c, "htmx-verify-003", err.Error())
	}

	// Delete the one-time code
	s.wellknownCookies.DeleteVerificationCodeCookie(c)

	switch verificationCode.VerifyCodePurpose {
	case contracts_cookies.VerifyCode_EmailVerification:
		// Go back to login
		return s.Render(c, http.StatusOK, "oidc/htmx/_partials/home", map[string]interface{}{
			"errors": []string{},
			"email":  verificationCode.Email,
		})

	case contracts_cookies.VerifyCode_PasswordReset:
		err = s.wellknownCookies.SetPasswordResetCookie(c,
			&contracts_cookies.SetPasswordResetCookieRequest{
				PasswordReset: &contracts_cookies.PasswordReset{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetPasswordResetCookie")
			return s.renderError(c, "htmx-verify-004", err.Error())
		}
		return s.Render(c, http.StatusOK, "oidc/htmx/_partials/reset-password", map[string]interface{}{
			"email":  verificationCode.Email,
			"errors": []string{},
		})

	case contracts_cookies.VerifyCode_Challenge:
		// MFA challenge passed - set auth cookies
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
			return s.renderError(c, "htmx-verify-005", err.Error())
		}

		err = s.wellknownCookies.SetAuthCompletedCookie(c,
			&contracts_cookies.SetAuthCompletedCookieRequest{
				AuthCompleted: &contracts_cookies.AuthCompleted{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetAuthCompletedCookie")
		}

		// Check if user has keep-signed-in preferences
		getKeepSigninPreferencesCookieResponse, err := s.wellknownCookies.GetKeepSigninPreferencesCookie(c, &contracts_cookies.GetKeepSigninPreferencesCookieRequest{
			Subject: verificationCode.Subject,
		})
		if err == nil && getKeepSigninPreferencesCookieResponse != nil && getKeepSigninPreferencesCookieResponse.KeepSigninPreferencesCookie != nil {
			// Auto-skip keep-signed-in page
			s.wellknownCookies.DeleteAuthCompletedCookie(c)
			err = s.wellknownCookies.SetSSOCookie(c,
				&contracts_cookies.SetSSOCookieRequest{
					SSOCookie: &contracts_cookies.SSOCookie{
						Identity: &proto_oidc_models.Identity{
							Subject:       rageUser.RootIdentity.Subject,
							Email:         rageUser.RootIdentity.Email,
							EmailVerified: rageUser.RootIdentity.EmailVerified,
							IdpSlug:       models.RootIdp,
						},
						Acr: []string{models.ACRPassword, models.ACRIdpRoot},
						Amr: []string{models.AMRPassword, models.AMRIdp, models.AMRMFA, models.AMREmailCode},
					},
				})
			if err != nil {
				log.Error().Err(err).Msg("SetSSOCookie")
				return s.renderError(c, "htmx-verify-006", err.Error())
			}

			return s.completeOAuthFlow(c, rageUser)
		}

		// Show keep-signed-in page
		return s.Render(c, http.StatusOK, "oidc/htmx/_partials/keep-signed-in", map[string]interface{}{})
	}

	return s.renderError(c, "htmx-verify-099", "Unknown verification purpose")
}

func (s *service) completeOAuthFlow(c *echo.Context, rageUser *proto_oidc_models.RageUser) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	session, err := s.oidcSession.GetSession()
	if err != nil {
		log.Error().Err(err).Msg("GetSession")
		return s.renderError(c, "htmx-verify-007", "Session error")
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		log.Error().Err(err).Msg("session.Get request")
		return s.renderError(c, "htmx-verify-008", "Session error")
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)
	rootPath := echo_utils.GetMyRootPath(c)

	result, err := s.ProcessFinalAuthenticationState(ctx, c,
		&services_echo_handlers_base.ProcessFinalAuthenticationStateRequest{
			AuthorizationRequest: authorizationRequest,
			Identity: &proto_oidc_models.OIDCIdentity{
				Subject:       rageUser.RootIdentity.Subject,
				Email:         rageUser.RootIdentity.Email,
				EmailVerified: rageUser.RootIdentity.EmailVerified,
				Acr:           []string{models.ACRPassword, models.ACRIdpRoot},
				Amr:           []string{models.AMRPassword, models.AMRIdp, models.AMRMFA, models.AMREmailCode},
			},
			RootPath: rootPath,
		})
	if err != nil {
		log.Error().Err(err).Msg("ProcessFinalAuthenticationState")
		return s.renderError(c, "htmx-verify-009", err.Error())
	}

	c.Response().Header().Set("HX-Redirect", result.RedirectURI)
	return c.NoContent(http.StatusOK)
}
