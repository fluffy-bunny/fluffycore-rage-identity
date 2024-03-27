package oidclogintotp

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
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
	qrcode "github.com/skip2/go-qrcode"
	gotp "github.com/xlzd/gotp"
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
		signinResponse   *contracts_cookies.GetSigninUserNameCookieResponse
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_OIDCLoginTOTP_001 = "rg-oidclogin-totp-001"
	InternalError_OIDCLoginTOTP_002 = "rg-oidclogin-totp-002"
	InternalError_OIDCLoginTOTP_003 = "rg-oidclogin-totp-003"
	InternalError_OIDCLoginTOTP_004 = "rg-oidclogin-totp-004"
	InternalError_OIDCLoginTOTP_005 = "rg-oidclogin-totp-005"
	InternalError_OIDCLoginTOTP_006 = "rg-oidclogin-totp-006"
	InternalError_OIDCLoginTOTP_007 = "rg-oidclogin-totp-007"
	InternalError_OIDCLoginTOTP_008 = "rg-oidclogin-totp-008"
	InternalError_OIDCLoginTOTP_009 = "rg-oidclogin-totp-009"
	InternalError_OIDCLoginTOTP_010 = "rg-oidclogin-totp-010"
	InternalError_OIDCLoginTOTP_011 = "rg-oidclogin-totp-011"

	InternalError_OIDCLoginTOTP_099 = "rg-oidclogin-totp-099"
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
		wellknown_echo.OIDCLoginTOTPPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LoginGetRequest struct {
	Error     string `param:"error" query:"error" form:"error" json:"error" xml:"error"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
}

type LoginTOTPPostRequest struct {
	UserName string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Code     string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
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
func (s *service) generatePNGQRCode(rageUser *proto_oidc_models.RageUser) string {
	totpSecret := rageUser.TOTP.Secret
	otp := gotp.NewDefaultTOTP(totpSecret)
	provisioningUri := otp.ProvisioningUri(rageUser.RootIdentity.Email, s.config.TOTPIssuerName)
	var pngB []byte
	pngB, _ = qrcode.Encode(provisioningUri, qrcode.Medium, 256)
	pngQRCode := base64.StdEncoding.EncodeToString(pngB)
	return pngQRCode
}
func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginTOTP_099)
	}
	log.Info().Interface("model", model).Msg("model")
	var rows []row
	var errors []string

	session, err := s.getSession()
	if err != nil {
		rows = append(rows, row{Key: "error", Value: err.Error()})
	}
	sessionRequest, err := session.Get("request")
	if err != nil {
		rows = append(rows, row{Key: "error", Value: err.Error()})
	}
	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	switch model.Directive {
	case models.IdentityFound:
		return s.handleIdentityFound(c, authorizationRequest.State)

	}

	renderError := func(errors []string) error {
		return s.Render(c, http.StatusBadRequest,
			"oidc/oidclogintotp/index",
			map[string]interface{}{
				"errors":    errors,
				"email":     s.signinResponse.Value.Email,
				"directive": models.LoginDirective,
			})
	}
	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: s.signinResponse.Value.Email,
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
	rageUser := getRageUserResponse.User

	pngQRCode := ""
	if !rageUser.TOTP.Verified {
		pngQRCode = s.generatePNGQRCode(rageUser)
	}

	return s.Render(c, http.StatusOK, "oidc/oidclogintotp/index",
		map[string]interface{}{
			"errors":    rows,
			"email":     s.signinResponse.Value.Email,
			"verified":  rageUser.TOTP.Verified,
			"directive": models.LoginDirective,
			"pngQRCode": pngQRCode,
		})
}

func (s *service) DoPost(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	r := c.Request()
	// is the request get or post?
	rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginTOTPPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginTOTP_099)
	}
	log.Info().Interface("model", model).Msg("model")
	if fluffycore_utils.IsEmptyOrNil(model.Code) {
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
	idps, _ := s.GetIDPs(ctx)

	authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

	model.UserName = strings.ToLower(model.UserName)

	renderError := func(rageUser *proto_oidc_models.RageUser, errors []string) error {
		pngQRCode := s.generatePNGQRCode(rageUser)
		return s.Render(c, http.StatusBadRequest,
			"oidc/oidclogintotp/index",
			map[string]interface{}{
				"errors":     errors,
				"email":      s.signinResponse.Value.Email,
				"idps":       idps,
				"directive":  models.LoginDirective,
				"hasPasskey": s.signinResponse.Value.HasPasskey,
				"pngQRCode":  pngQRCode,
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
			return c.Redirect(http.StatusFound, "/error")
		}
		err = nil
	}

	rageUser := getRageUserResponse.User
	totpSecret := rageUser.TOTP.Secret
	otp := gotp.NewDefaultTOTP(totpSecret)
	valid := otp.Verify(model.Code, time.Now().Unix())
	if !valid {
		log.Warn().Err(err).Msg("TOTP Code is invalid")
		ee := utils.LocalizeWithInterperlate(localizer, "code.is.invalid", nil)
		errors = append(errors, ee)
		return renderError(rageUser, errors)
	}
	if !rageUser.TOTP.Verified {
		rageUser.TOTP.Verified = true
		_, err = s.RageUserService().UpdateRageUser(ctx,
			&proto_oidc_user.UpdateRageUserRequest{
				User: &proto_oidc_models.RageUserUpdate{
					RootIdentity: &proto_oidc_models.IdentityUpdate{
						Subject: rageUser.RootIdentity.Subject,
					},
					TOTP: &proto_oidc_models.TOTPUpdate{
						Verified: &wrapperspb.BoolValue{Value: true},
					},
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("UpdateRageUser")
			errors = append(errors, err.Error())
			return renderError(rageUser, errors)
		}
	}
	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &proto_oidc_models.Identity{
				Subject:       rageUser.RootIdentity.Subject,
				Email:         rageUser.RootIdentity.Email,
				EmailVerified: rageUser.RootIdentity.EmailVerified,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		errors = append(errors, err.Error())
		// redirect to error page
		return renderError(rageUser, errors)
	}

	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
		State: authorizationRequest.State,
	})
	if err != nil {
		log.Warn().Err(err).Msg("GetAuthorizationRequestState")
		errors = append(errors, err.Error())
		return renderError(rageUser, errors)
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
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
		},
	}

	// "urn:rage:idp:google", "urn:rage:idp:spacex", "urn:rage:idp:github-enterprise", etc.
	// "urn:rage:password", "urn:rage:2fa", "urn:rage:email", etc.
	// we are done with the state now.  Lets map it to the code so it can be looked up by the client.
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
		State:                     authorizationFinal.Request.Code,
		AuthorizationRequestState: authorizationFinal,
	})
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationRequestState")
		// redirect to error page
		errors = append(errors, err.Error())
		return renderError(rageUser, errors)
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
		return renderError(rageUser, errors)
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
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	signinResponse, err := s.wellknownCookies.GetSigninUserNameCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetSigninUserNameCookie")
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginTOTP_004)
	}
	s.signinResponse = signinResponse

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
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginTOTP_002)
	}
	_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
		State:                     authorizationFinal.Request.Code,
		AuthorizationRequestState: authorizationFinal,
	})
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationRequestState")
		// redirect to error page
		return s.TeleportBackToLogin(c, InternalError_OIDCLoginTOTP_003)
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
