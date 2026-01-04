package oidclogin

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
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
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

var _ contracts_handler.IHandler = stemService

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_OIDCLogin_001 = "rg-oidclogin-001"
	InternalError_OIDCLogin_002 = "rg-oidclogin-002"
	InternalError_OIDCLogin_003 = "rg-oidclogin-003"
	InternalError_OIDCLogin_004 = "rg-oidclogin-004"
	InternalError_OIDCLogin_005 = "rg-oidclogin-005"
	InternalError_OIDCLogin_006 = "rg-oidclogin-006"
	InternalError_OIDCLogin_007 = "rg-oidclogin-007"
	InternalError_OIDCLogin_008 = "rg-oidclogin-008"
	InternalError_OIDCLogin_009 = "rg-oidclogin-009"
	InternalError_OIDCLogin_010 = "rg-oidclogin-010"
	InternalError_OIDCLogin_011 = "rg-oidclogin-011"
	InternalError_OIDCLogin_099 = "rg-oidclogin-099"
)

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	sessionFactory contracts_sessions.ISessionFactory,
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

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			// Using only auto post here so that our arguments are present in the URL
			//	contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.OIDCLoginPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LoginGetRequest struct {
	Email string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	//Error            string            `param:"error" query:"error" form:"error" json:"error" xml:"error"`
	//ErrorCode        string            `param:"error_code" query:"error_code" form:"error_code" json:"error_code" xml:"error_code"`
	Directive        string            `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
	AdditionalParams map[string]string `param:"additional_params" query:"additional_params" form:"additional_params" json:"additional_params" xml:"additional_params"`
}
type ExternalIDPAuthRequest struct {
	IDPHint string `param:"idp_hint" query:"idp_hint" form:"idp_hint" json:"idp_hint" xml:"idp_hint"`
}
type LoginPostRequest struct {
	UserName  string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
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
		return s.TeleportBackToLoginWithError(c, InternalError_OIDCLogin_099, InternalError_OIDCLogin_099)
	}
	log.Debug().Interface("model", model).Msg("model")
	var errors []string

	/*
		// not going to write errors here anymore.
		// the originator of the error should drop their own cookie.
			 		// Check for errors from query parameters
					if fluffycore_utils.IsNotEmptyOrNil(model.Error) {
						code := InternalError_OIDCLogin_099
						if fluffycore_utils.IsNotEmptyOrNil(model.ErrorCode) {
							code = model.ErrorCode
						}
						s.wellknownCookies.SetErrorCookie(c, &contracts_cookies.SetErrorCookieRequest{
							Value: &contracts_cookies.ErrorCookie{
								Code:  code,
								Error: model.Error,
							},
						})
						errors = append(errors, model.Error)
					}

	*/

	// Check for errors from cookie (e.g., from OAuth2 callback redirects)
	// Note: Don't delete the cookie here - let the WASM app read and delete it
	errorCookieResponse, err := s.wellknownCookies.GetErrorCookie(c)
	if err == nil && errorCookieResponse != nil && errorCookieResponse.Value != nil {
		if fluffycore_utils.IsNotEmptyOrNil(errorCookieResponse.Value.Error) {
			log.Info().
				Str("error", errorCookieResponse.Value.Error).
				Str("code", errorCookieResponse.Value.Code).
				Msg("Found error cookie (will be displayed by WASM app)")
			errors = append(errors, errorCookieResponse.Value.Error)
			// Don't delete the cookie here - the WASM app will read and delete it
		}
	} else if err != nil {
		log.Debug().Err(err).Msg("No error cookie found")
	}
	session, err := s.getSession()
	if err != nil {
		errors = append(errors, err.Error())
	}
	requestSession, err := session.Get("request")
	if err != nil {
		errors = append(errors, err.Error())
	}
	authorizationRequest := requestSession.(*proto_oidc_models.AuthorizationRequest)

	log.Debug().Interface("requestSession", requestSession).Msg("requestSession")

	switch model.Directive {
	case models.IdentityFound:
		return s.handleIdentityFound(c, authorizationRequest.State)

	}
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		errors = append(errors, err.Error())
	}
	if s.config.OIDCUIConfig.URIEntryPath != wellknown_echo.OIDCLoginPath {

		// we redirect over to URIEntryPath
		return c.Redirect(http.StatusFound, s.config.OIDCUIConfig.URIEntryPath)
	}
	return s.Render(c, http.StatusOK, "oidc/oidclogin/index",
		map[string]interface{}{
			"errors":    errors,
			"idps":      idps,
			"email":     model.Email,
			"directive": models.LoginDirective,
		})
}

func (s *service) DoPost(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	var err error
	model := &LoginPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLoginWithError(c, InternalError_OIDCLogin_099, InternalError_OIDCLogin_099)
	}
	log.Debug().Interface("model", model).Msg("model")

	doVerifyCode := model.Directive == models.MFA_VerifyEmailDirective || model.Directive == models.VerifyEmailDirective
	doKeepSignedIn := model.Directive == models.KeepSignedInDirective

	if !doVerifyCode && !doKeepSignedIn && fluffycore_utils.IsEmptyOrNil(model.UserName) {
		return s.DoGet(c)
	}

	var errors []string

	session, err := s.getSession()
	if err != nil {
		return s.DoGet(c)
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		return s.DoGet(c)
	}
	if doVerifyCode {
		landingPage := &models_api_manifest.LandingPage{
			Page: models_api_manifest.PageVerifyCode,
		}
		session.Set("landingPage", landingPage)
		session.Save()
		return s.DoGet(c)
	}
	if doKeepSignedIn {

		// Check if user has KeepSigninPreferences cookie ("don't show again" preference)
		// Get the auth cookie to retrieve the subject
		getAuthCookieResponse, err := s.wellknownCookies.GetAuthCookie(c)
		if err == nil && getAuthCookieResponse.AuthCookie != nil {
			getPreferencesCookieResponse, err := s.wellknownCookies.GetKeepSigninPreferencesCookie(c,
				&contracts_cookies.GetKeepSigninPreferencesCookieRequest{
					Subject: getAuthCookieResponse.AuthCookie.Identity.Subject,
				})
			if err == nil && getPreferencesCookieResponse.KeepSigninPreferencesCookie != nil && getPreferencesCookieResponse.KeepSigninPreferencesCookie.PreferenceValue {
				log.Info().Msg("Skipping keep-signed-in page due to KeepSigninPreferences cookie")

				// Set SSO cookie since we're auto-keeping them signed in - copy all auth context
				err = s.wellknownCookies.SetSSOCookie(c,
					&contracts_cookies.SetSSOCookieRequest{
						SSOCookie: &contracts_cookies.SSOCookie{
							Identity: getAuthCookieResponse.AuthCookie.Identity,
							Acr:      getAuthCookieResponse.AuthCookie.Acr,
							Amr:      getAuthCookieResponse.AuthCookie.Amr,
						},
					})
				if err != nil {
					log.Error().Err(err).Msg("SetSSOCookie failed when auto-skipping keep-signed-in")
				} else {
					log.Info().Str("subject", getAuthCookieResponse.AuthCookie.Identity.Subject).Msg("SSO cookie set for auto keep signed in")
				}

				// Delete the AuthCompleted cookie (one-time use)
				s.wellknownCookies.DeleteAuthCompletedCookie(c)

				// Go directly to handleIdentityFound to complete the OAuth flow
				authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)
				return s.handleIdentityFound(c, authorizationRequest.State)
			}
		}

		landingPage := &models_api_manifest.LandingPage{
			Page: models_api_manifest.PageKeepSignedIn,
		}
		session.Set("landingPage", landingPage)
		session.Save()
		return s.DoGet(c)
	}
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		errors = append(errors, err.Error())
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	log.Debug().Interface("sessionRequest", sessionRequest).Msg("sessionRequest")

	model.UserName = strings.ToLower(model.UserName)

	email, ok := echo_utils.IsValidEmailAddress(model.UserName)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "username.not.valid", map[string]string{"username": model.UserName})

		errors = append(errors, msg)
		if s.config.OIDCUIConfig.URIEntryPath != wellknown_echo.OIDCLoginPath {
			// we redirect over to URIEntryPath
			return c.Redirect(http.StatusFound, s.config.OIDCUIConfig.URIEntryPath)
		}
		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"idps":      idps,
				"errors":    errors,
				"directive": models.LoginDirective,
			})
	}
	// get the domain from the email
	parts := strings.Split(email, "@")
	domainPart := parts[1]

	// first lets see if this domain has been claimed.
	listIDPRequest, err := s.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			Enabled: &proto_types.BoolFilterExpression{
				Eq: true,
			},
			ClaimedDomains: &proto_types.StringArrayFilterExpression{
				Eq: domainPart,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("ListIDP")
		errors = append(errors, err.Error())
		if s.config.OIDCUIConfig.URIEntryPath != wellknown_echo.OIDCLoginPath {
			// we redirect over to URIEntryPath
			return c.Redirect(http.StatusFound, s.config.OIDCUIConfig.URIEntryPath)
		}
		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     authorizationRequest.State,
				"idps":      idps,
				"errors":    errors,
				"directive": models.LoginDirective,
			})
	}
	if len(listIDPRequest.IDPs) > 0 {
		// an idp has claimed this domain.
		// post to the externalIDP

		return s.RenderAutoPost(c, wellknown_echo.ExternalIDPPath,
			[]models.FormParam{
				{
					Name:  "state",
					Value: authorizationRequest.State,
				},
				{
					Name:  "idp_hint",
					Value: listIDPRequest.IDPs[0].Slug,
				},
				{
					Name:  "directive",
					Value: models.LoginDirective,
				},
			})

	}
	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: strings.ToLower(model.UserName),
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.NotFound {
			return err
		}
		err = nil
	}
	if getRageUserResponse == nil {
		msg := utils.LocalizeWithInterperlate(localizer, "username.not.found", map[string]string{"username": model.UserName})

		errors = append(errors, msg)
		if s.config.OIDCUIConfig.URIEntryPath != wellknown_echo.OIDCLoginPath {
			// we redirect over to URIEntryPath
			return c.Redirect(http.StatusFound, s.config.OIDCUIConfig.URIEntryPath)
		}
		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     authorizationRequest.State,
				"idps":      idps,
				"errors":    errors,
				"directive": models.LoginDirective,
			})

	}
	if err != nil {
		log.Warn().Err(err).Msg("ListUser")
		errors = append(errors, err.Error())
		if s.config.OIDCUIConfig.URIEntryPath != wellknown_echo.OIDCLoginPath {
			// we redirect over to URIEntryPath
			return c.Redirect(http.StatusFound, s.config.OIDCUIConfig.URIEntryPath)
		}
		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     authorizationRequest.State,
				"idps":      idps,
				"errors":    errors,
				"directive": models.LoginDirective,
			})
	}
	user := getRageUserResponse.User
	if s.config.EmailVerificationRequired && !user.RootIdentity.EmailVerified {
		codeResult, err := echo_utils.GenerateHashedVerificationCode(ctx, s.passwordHasher, 6)
		if err != nil {
			log.Error().Err(err).Msg("GenerateHashedVerificationCode")
			return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:             model.UserName,
					PlainCode:         codeResult.PlainCode,
					CodeHash:          codeResult.HashedCode,
					Subject:           user.RootIdentity.Subject,
					VerifyCodePurpose: contracts_cookies.VerifyCode_EmailVerification,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return s.TeleportBackToLoginWithError(c, InternalError_OIDCLogin_001, InternalError_OIDCLogin_001)
		}
		s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.UserName,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification.message",
				Data: map[string]string{
					"code": codeResult.PlainCode,
				},
			})
		formParams := []models.FormParam{
			{
				Name:  "state",
				Value: authorizationRequest.State,
			},
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
				Value: codeResult.PlainCode,
			})
		}
		return s.RenderAutoPost(c, wellknown_echo.VerifyCodePath, formParams)
	}
	hasPasskey := false
	if user.WebAuthN != nil && fluffycore_utils.IsNotEmptyOrNil(user.WebAuthN.Credentials) {
		hasPasskey = true
	}
	err = s.wellknownCookies.SetSigninUserNameCookie(c, &contracts_cookies.SetSigninUserNameCookieRequest{
		Value: &contracts_cookies.SigninUserNameCookie{
			Email:      model.UserName,
			HasPasskey: hasPasskey,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetSigninUserNameCookie")
		return s.TeleportBackToLoginWithError(c, InternalError_OIDCLogin_004, InternalError_OIDCLogin_004)
	}
	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPasswordPath,
		[]models.FormParam{
			{
				Name:  "state",
				Value: authorizationRequest.State,
			},
			{
				Name:  "directive",
				Value: models.LoginDirective,
			},
		})

}

// OIDC Login godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       code            		query     string  true  "code"
// @Success 200 {object} string
// @Router /oidc-login [get]
// @Router /oidc-login [post]
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
				EmailVerified: authorizationFinal.Identity.EmailVerified,
				IdpSlug:       authorizationFinal.Identity.IdpSlug,
			},
			Acr: authorizationFinal.Identity.Acr,
			Amr: authorizationFinal.Identity.Amr,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		// redirect to error page
		return s.TeleportBackToLoginWithError(c, InternalError_OIDCLogin_002, InternalError_OIDCLogin_002)
	}

	// Set AuthCompleted cookie to mark successful authentication
	err = s.wellknownCookies.SetAuthCompletedCookie(c,
		&contracts_cookies.SetAuthCompletedCookieRequest{
			AuthCompleted: &contracts_cookies.AuthCompleted{
				Subject: authorizationFinal.Identity.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCompletedCookie")
		return s.TeleportBackToLoginWithError(c, InternalError_OIDCLogin_003, InternalError_OIDCLogin_003)
	}

	// Return directive to navigate to keep-signed-in page instead of redirecting to client
	formParams := []models.FormParam{
		{
			Name:  "state",
			Value: state,
		},
		{
			Name:  "directive",
			Value: models.KeepSignedInDirective,
		},
	}
	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath, formParams)

}
