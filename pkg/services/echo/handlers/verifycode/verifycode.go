package forgotpassword

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
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
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		oidcSession      contracts_oidc_session.IOIDCSession
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_VerifyCode_001 = "rg-verifycode-001"
	InternalError_VerifyCode_002 = "rg-verifycode-002"
	InternalError_VerifyCode_003 = "rg-verifycode-003"
	InternalError_VerifyCode_004 = "rg-verifycode-004"
	InternalError_VerifyCode_005 = "rg-verifycode-005"
	InternalError_VerifyCode_006 = "rg-verifycode-006"
	InternalError_VerifyCode_007 = "rg-verifycode-007"
	InternalError_VerifyCode_008 = "rg-verifycode-008"
	InternalError_VerifyCode_009 = "rg-verifycode-009"
	InternalError_VerifyCode_010 = "rg-verifycode-010"
	InternalError_VerifyCode_011 = "rg-verifycode-011"
	InternalError_VerifyCode_099 = "rg-verifycode-099" // 99 is a bind problem

)

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	oidcSession contracts_oidc_session.IOIDCSession,

) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies: wellknownCookies,
		oidcSession:      oidcSession,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			// do auto post
			//contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.VerifyCodePath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type VerifyCodeGetRequest struct {
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Code      string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
}

type VerifyCodePostRequest struct {
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Code      string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
	Type      string `param:"type" query:"type" form:"type" json:"type" xml:"type"`
	Action    string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
}

func (s *service) validateVerifyCodeGetRequest(model *VerifyCodeGetRequest) error {

	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		return status.Error(codes.InvalidArgument, "Email is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Directive) {
		return status.Error(codes.InvalidArgument, "Directive is empty")
	}
	return nil
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &VerifyCodeGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return s.TeleportBackToLogin(c, InternalError_VerifyCode_099)
	}
	log.Debug().Interface("model", model).Msg("model")
	err := s.validateVerifyCodeGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validateVerifyCodeGetRequest")
		return s.TeleportBackToLogin(c, InternalError_VerifyCode_001)
	}

	err = s.Render(c, http.StatusOK, "oidc/verifycode/index",
		map[string]interface{}{
			"email":     model.Email,
			"code":      model.Code,
			"directive": model.Directive,
			"errors":    make([]string, 0),
		})
	return err
}

func (s *service) validateVerifyCodePostRequest(request *VerifyCodePostRequest) ([]string, error) {
	localizer := s.Localizer().GetLocalizer()
	var err error
	errors := make([]string, 0)

	if fluffycore_utils.IsEmptyOrNil(request.Email) {
		msg := utils.LocalizeSimple(localizer, "email.is.empty")
		errors = append(errors, msg)
	}
	if fluffycore_utils.IsEmptyOrNil(request.Code) {
		msg := utils.LocalizeSimple(localizer, "code.is.empty")
		errors = append(errors, msg)
	}
	_, ok := echo_utils.IsValidEmailAddress(request.Email)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "email.is.not.valid", map[string]string{"email": request.Email})
		errors = append(errors, msg)
	}
	if fluffycore_utils.IsEmptyOrNil(request.Directive) {
		msg := utils.LocalizeSimple(localizer, "directive.is.empty")
		errors = append(errors, msg)
	}
	return errors, err
}

func (s *service) DoPost(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	r := c.Request()
	// is the request get or post?
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &VerifyCodePostRequest{}
	if err := c.Bind(model); err != nil {
		return s.TeleportBackToLogin(c, InternalError_VerifyCode_099)
	}
	log.Debug().Interface("model", model).Msg("model")

	errors, err := s.validateVerifyCodePostRequest(model)
	if err != nil {
		return s.Render(c, http.StatusBadRequest, "oidc/verifycode/index",
			map[string]interface{}{
				"email":     model.Email,
				"code":      model.Code,
				"directive": model.Directive,
				"errors":    errors,
			})
	}
	if model.Type == "GET" {
		return s.DoGet(c)
	}
	if model.Action == "cancel" {
		return s.TeleportToPath(c, wellknown_echo.OIDCLoginPath)
	}
	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetVerificationCodeCookie")
		return s.RenderAutoPost(c, wellknown_echo.ForgotPasswordPath,
			[]models.FormParam{

				{
					Name:  "email",
					Value: model.Email,
				},
				{
					Name:  "type",
					Value: "GET",
				},
			})
	}
	verificationCode := getVerificationCodeCookieResponse.VerificationCode
	code := verificationCode.Code

	if code != model.Code {
		return s.Render(c, http.StatusBadRequest, "oidc/verifycode/index",
			map[string]interface{}{
				"email":     model.Email,
				"code":      model.Code,
				"directive": model.Directive,
				"errors": []string{
					utils.LocalizeSimple(localizer, "code.is.invalid"),
				},
			})
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
		return s.TeleportBackToLogin(c, InternalError_VerifyCode_002)
	}
	rageUser := getRageUserResponse.User

	_, err = userService.UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: verificationCode.Subject,
				EmailVerified: &wrapperspb.BoolValue{
					Value: true,
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return s.TeleportBackToLogin(c, InternalError_VerifyCode_004)
	}
	// one time only
	s.wellknownCookies.DeleteVerificationCodeCookie(c)

	redirectURL := "/"
	switch model.Directive {
	case models.PasswordResetDirective:
		err = s.wellknownCookies.SetPasswordResetCookie(c,
			&contracts_cookies.SetPasswordResetCookieRequest{
				PasswordReset: &contracts_cookies.PasswordReset{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetPasswordResetCookie")
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_003)
		}
		return s.RenderAutoPost(c, wellknown_echo.PasswordResetPath,
			[]models.FormParam{})
	case models.VerifyEmailDirective:
		return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
			[]models.FormParam{
				{
					Name:  "email",
					Value: model.Email,
				},
			})
	case models.MFA_VerifyEmailDirective:
		oidcIdentity := &proto_oidc_models.OIDCIdentity{
			Subject: rageUser.RootIdentity.Subject,
			Email:   rageUser.RootIdentity.Email,
			Acr: []string{
				models.ACRPassword,
				models.ACRIdpRoot,
			},
			Amr: []string{
				models.AMRPassword,
				// always true, as we are the root idp
				models.AMRIdp,
				// this is a multifactor
				models.AMRMFA,
				models.AMREmailCode,
			},
		}
		err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
			AuthCookie: &contracts_cookies.AuthCookie{
				Identity: &proto_oidc_models.Identity{
					Subject:       rageUser.RootIdentity.Subject,
					Email:         rageUser.RootIdentity.Email,
					EmailVerified: rageUser.RootIdentity.EmailVerified,
				},
				Acr: oidcIdentity.Acr,
				Amr: oidcIdentity.Amr,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("SetAuthCookie")
			// redirect to error page
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_005)
		}
		session, err := s.getSession()
		if err != nil {
			log.Error().Err(err).Msg("getSession")
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_006)
		}
		sessionRequest, err := session.Get("request")
		if err != nil {
			log.Error().Err(err).Msg("Get")
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_007)
		}
		authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

		getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().
			GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
				State: authorizationRequest.State,
			})
		if err != nil {
			log.Error().Err(err).Msg("GetAuthorizationRequestState")
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_008)
		}
		authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
		authorizationFinal.Identity = oidcIdentity
		// "urn:rage:idp:google", "urn:rage:idp:spacex", "urn:rage:idp:github-enterprise", etc.
		// "urn:rage:password", "urn:rage:2fa", "urn:rage:email", etc.
		// we are done with the state now.  Lets map it to the code so it can be looked up by the client.
		_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     authorizationFinal.Request.Code,
			AuthorizationRequestState: authorizationFinal,
		})
		if err != nil {
			log.Error().Err(err).Msg("StoreAuthorizationRequestState")
			// redirect to error page
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_009)
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
			log.Error().Err(err).Msg("StoreAuthorizationRequestState")
			return s.TeleportBackToLogin(c, InternalError_VerifyCode_010)
		}
		rootPath := echo_utils.GetMyRootPath(c)

		// redirect to the client with the code.
		redirectUri := authorizationFinal.Request.RedirectUri +
			"?code=" + authorizationFinal.Request.Code +
			"&state=" + authorizationFinal.Request.State +
			"&iss=" + rootPath
		return c.Redirect(http.StatusFound, redirectUri)
	}

	return c.Redirect(http.StatusFound, redirectURL)

}
func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()

	if err != nil {
		return nil, err
	}
	return session, nil
}

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
