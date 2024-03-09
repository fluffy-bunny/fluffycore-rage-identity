package oidcloginpassword

import (
	"fmt"
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
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

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_OIDCLoginPassword_001 = "rg-oidclogin-password-001"
	InternalError_OIDCLoginPassword_002 = "rg-oidclogin-password-002"
	InternalError_OIDCLoginPassword_003 = "rg-oidclogin-password-003"
	InternalError_OIDCLoginPassword_004 = "rg-oidclogin-password-004"
	InternalError_OIDCLoginPassword_005 = "rg-oidclogin-password-005"
	InternalError_OIDCLoginPassword_006 = "rg-oidclogin-password-006"
	InternalError_OIDCLoginPassword_007 = "rg-oidclogin-password-007"
	InternalError_OIDCLoginPassword_008 = "rg-oidclogin-password-008"
	InternalError_OIDCLoginPassword_009 = "rg-oidclogin-password-009"
	InternalError_OIDCLoginPassword_010 = "rg-oidclogin-password-010"
	InternalError_OIDCLoginPassword_011 = "rg-oidclogin-password-011"

	InternalError_OIDCLoginPassword_099 = "rg-oidclogin-password-099"
)

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			// Using only auto post here so that our arguments are present in the URL
			//	contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.OIDCLoginPasswordPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LoginGetRequest struct {
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Error     string `param:"error" query:"error" form:"error" json:"error" xml:"error"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
}

type LoginPasswordPostRequest struct {
	UserName string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
}

type row struct {
	Key   string
	Value string
}

func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}
func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginPassword_099)
	}
	log.Info().Interface("model", model).Msg("model")
	var rows []row

	session, err := s.getSession()
	if err != nil {
		rows = append(rows, row{Key: "error", Value: err.Error()})
	}
	dd, err := session.Get("request")
	if err != nil {
		rows = append(rows, row{Key: "error", Value: err.Error()})
	}
	dd2 := dd.(*proto_oidc_models.AuthorizationRequest)

	switch model.Directive {
	case models.IdentityFound:
		return s.handleIdentityFound(c, dd2.State)

	}
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		rows = append(rows, row{Key: "error", Value: err.Error()})
	}

	return s.Render(c, http.StatusOK, "oidc/oidcloginpassword/index",
		map[string]interface{}{
			"errors":    rows,
			"idps":      idps,
			"email":     model.Email,
			"directive": models.LoginDirective,
		})
}

func (s *service) DoPost(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	r := c.Request()
	// is the request get or post?
	rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginPasswordPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginPassword_099)
	}
	log.Info().Interface("model", model).Msg("model")
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		return s.DoGet(c)
	}

	var errors []string
	session, err := s.getSession()
	if err != nil {
		errors = append(errors, err.Error())
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		errors = append(errors, err.Error())
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	model.UserName = strings.ToLower(model.UserName)

	renderError := func(errors []string) error {
		return s.Render(c, http.StatusBadRequest,
			"oidc/oidcloginpassword/index",
			map[string]interface{}{
				"email":     model.UserName,
				"errors":    errors,
				"directive": models.LoginDirective,
			})
	}
	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: model.UserName,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.NotFound {
			log.Warn().Err(err).Msg("GetRageUser")
			errors = append(errors, err.Error())
			return renderError(errors)
		}
		err = nil
	}

	user := getRageUserResponse.User
	if s.config.EmailVerificationRequired && !user.RootIdentity.EmailVerified {

		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:   model.UserName,
					Code:    verificationCode,
					Subject: user.RootIdentity.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return s.TeleportBackToLogin(c, InternalError_OIDCLoginPassword_001)
		}
		s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.UserName,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification..message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})
		formParams := []models.FormParam{

			{
				Name:  "email",
				Value: model.UserName,
			},
			{
				Name:  "directive",
				Value: models.VerifyEmailDirective,
			},
			{
				Name:  "type",
				Value: "GET",
			},
		}

		if s.config.SystemConfig.DeveloperMode {
			formParams = append(formParams, models.FormParam{
				Name:  "code",
				Value: verificationCode,
			})
		}
		return s.RenderAutoPost(c, wellknown_echo.VerifyCodePath, formParams)
	}

	if user.Password == nil {
		ee := utils.LocalizeWithInterperlate(localizer, "username.does.not.have.password", map[string]string{
			"username": model.UserName,
		})

		errors = append(errors, ee)
		return renderError(errors)
	}

	err = s.passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
		Password:       model.Password,
		HashedPassword: user.Password.Hash,
	})
	if err != nil {
		log.Warn().Err(err).Msg("ComparePasswordHash")
		ee := utils.LocalizeWithInterperlate(localizer, "password.is.invalid", nil)
		errors = append(errors, ee)
		return renderError(errors)
	}
	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &proto_oidc_models.Identity{
				Subject:       user.RootIdentity.Subject,
				Email:         user.RootIdentity.Email,
				EmailVerified: user.RootIdentity.EmailVerified,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		errors = append(errors, err.Error())
		// redirect to error page
		return renderError(errors)
	}

	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
		State: authorizationRequest.State,
	})
	if err != nil {
		log.Warn().Err(err).Msg("GetAuthorizationRequestState")
		errors = append(errors, err.Error())
		return renderError(errors)
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
		Subject: user.RootIdentity.Subject,
		Email:   user.RootIdentity.Email,
		Acr: []string{
			models.ACRPassword,
			models.ACRIdpRoot,
		},
		Amr: []string{
			models.AMRPassword,
			// always true, as we are the root idp
			models.AMRIdp,
		},
	}

	// "urn:mastodon:idp:google", "urn:mastodon:idp:spacex", "urn:mastodon:idp:github-enterprise", etc.
	// "urn:mastodon:password", "urn:mastodon:2fa", "urn:mastodon:email", etc.
	// we are done with the state now.  Lets map it to the code so it can be looked up by the client.
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
		State:                     authorizationFinal.Request.Code,
		AuthorizationRequestState: authorizationFinal,
	})
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationRequestState")
		// redirect to error page
		errors = append(errors, err.Error())
		return renderError(errors)
	}
	s.AuthorizationRequestStateStore().DeleteAuthorizationRequestState(ctx, &proto_oidc_flows.DeleteAuthorizationRequestStateRequest{
		State: authorizationRequest.State,
	})
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
		State:                     authorizationRequest.State,
		AuthorizationRequestState: authorizationFinal,
	})
	if err != nil {
		// redirect to error page
		errors = append(errors, err.Error())
		return renderError(errors)
	}
	// redirect to the client with the code.
	redirectUri := authorizationFinal.Request.RedirectUri +
		"?code=" + authorizationFinal.Request.Code +
		"&state=" + authorizationFinal.Request.State +
		"&iss=" + rootPath
	return c.Redirect(http.StatusFound, redirectUri)

}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       code            		query     string  true  "code"
// @Success 200 {object} string
// @Router /login [get,post]
func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
func (s *service) handleIdentityFound(c echo.Context, state string) error {
	r := c.Request()
	// is the request get or post?
	rootPath := s.config.OIDCConfig.BaseUrl
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
		State: state,
	})
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		// redirect to error page
		redirectUrl := fmt.Sprintf("%s?state=%s&error=%s", wellknown_echo.OIDCLoginPath, state, models.InternalError)
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	if authorizationFinal.Identity == nil {
		redirectUrl := fmt.Sprintf("%s?state=%s&error=%s", wellknown_echo.OIDCLoginPath, state, models.InternalError)
		return c.Redirect(http.StatusFound, redirectUrl)
	}

	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &proto_oidc_models.Identity{
				Subject:       authorizationFinal.Identity.Subject,
				Email:         authorizationFinal.Identity.Email,
				IdpSlug:       authorizationFinal.Identity.IdpSlug,
				EmailVerified: authorizationFinal.Identity.EmailVerified,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		// redirect to error page
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginPassword_002)
	}
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
		State:                     authorizationFinal.Request.Code,
		AuthorizationRequestState: authorizationFinal,
	})
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationRequestState")
		// redirect to error page
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginPassword_003)
	}
	s.AuthorizationRequestStateStore().DeleteAuthorizationRequestState(ctx, &proto_oidc_flows.DeleteAuthorizationRequestStateRequest{
		State: state,
	})

	// redirect to the client with the code.
	redirectUri := authorizationFinal.Request.RedirectUri +
		"?code=" + authorizationFinal.Request.Code +
		"&state=" + authorizationFinal.Request.State +
		"&iss=" + rootPath
	return c.Redirect(http.StatusFound, redirectUri)

}
