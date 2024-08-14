package api_login_username_phase_one

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
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

		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
		oidcSession      contracts_oidc_session.IOIDCSession
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		wellknownCookies: wellknownCookies,
		passwordHasher:   passwordHasher,
		oidcSession:      oidcSession,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_LoginPhaseOne,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateLoginPhaseOneRequest(model *login_models.LoginPhaseOneRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		return status.Error(codes.InvalidArgument, "model.Email is nil")
	}
	_, ok := echo_utils.IsValidEmailAddress(model.Email)
	if !ok {
		return status.Error(codes.InvalidArgument, "model.Email is not a valid email address")

	}
	return nil
}

// API LoginUserPhaseOne godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept json
// @Produce json
// @Param		request body		login_models.LoginPhaseOneRequest	true	"LoginPhaseOneRequest"
// @Success 200 {object} login_models.LoginPhaseOneResponse
// @Router /api/login-phase-one [post]
func (s *service) Do(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &login_models.LoginPhaseOneRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	if err := s.validateLoginPhaseOneRequest(model); err != nil {
		log.Error().Err(err).Msg("validateLoginPhaseOneRequest")
		return c.JSONPretty(http.StatusBadRequest, err.Error(), "  ")
	}
	response := &login_models.LoginPhaseOneResponse{
		Email: model.Email,
	}

	session, err := s.getSession()
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")

	}

	log.Debug().Interface("sessionRequest", sessionRequest).Msg("sessionRequest")

	model.Email = strings.ToLower(model.Email)

	email, ok := echo_utils.IsValidEmailAddress(model.Email)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "username.not.valid", map[string]string{"username": model.Email})
		return c.JSONPretty(http.StatusBadRequest, msg, "  ")
	}
	// get the domain from the email
	parts := strings.Split(email, "@")
	domainPart := parts[1]
	// first lets see if this domain has been claimed.
	listIDPRequest, err := s.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			ClaimedDomain: &proto_types.StringArrayFilterExpression{
				Eq: domainPart,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("ListIDP")
		return c.JSONPretty(http.StatusBadRequest, err.Error(), "  ")
	}
	if len(listIDPRequest.Idps) > 0 {
		// an idp has claimed this domain.
		// lets start that session and return the redirect URI to the externalIDP
		// post to the externalIDP
		response.Directive = login_models.DIRECTIVE_StartExternalLogin
		response.DirectiveStartExternalLogin = &login_models.DirectiveStartExternalLogin{
			Slug: listIDPRequest.Idps[0].Slug,
		}

		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: strings.ToLower(model.Email),
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, err.Error(), "  ")
		}
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	if getRageUserResponse == nil {
		return c.JSONPretty(http.StatusNotFound, "User not found", "  ")
	}
	user := getRageUserResponse.User
	if s.config.EmailVerificationRequired && !user.RootIdentity.EmailVerified {
		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:             model.Email,
					Code:              verificationCode,
					Subject:           user.RootIdentity.Subject,
					VerifyCodePurpose: contracts_cookies.VerifyCode_EmailVerification,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		_, err = s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.Email,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification.message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendSimpleEmail")
			return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
		}
		if s.config.SystemConfig.DeveloperMode {
			response.DirectiveEmailCodeChallenge = &login_models.DirectiveEmailCodeChallenge{
				Code: verificationCode,
			}
		}
		response.Directive = login_models.DIRECTIVE_VerifyCode_DisplayVerifyCodePage
		return c.JSONPretty(http.StatusOK, response, "  ")
	}
	hasPasskey := false
	if user.WebAuthN != nil && fluffycore_utils.IsNotEmptyOrNil(user.WebAuthN.Credentials) {
		hasPasskey = true
	}
	err = s.wellknownCookies.SetSigninUserNameCookie(c,
		&contracts_cookies.SetSigninUserNameCookieRequest{
			Value: &contracts_cookies.SigninUserNameCookie{
				Email:      model.Email,
				HasPasskey: hasPasskey,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetSigninUserNameCookie")
		return c.JSONPretty(http.StatusInternalServerError, err.Error(), "  ")
	}
	response.Directive = login_models.DIRECTIVE_LoginPhaseOne_DisplayPasswordPage
	response.DirectiveDisplayPasswordPage = &login_models.DirectiveDisplayPasswordPage{
		Email:      model.Email,
		HasPasskey: hasPasskey,
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
