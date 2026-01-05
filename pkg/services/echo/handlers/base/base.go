package base

import (
	"context"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/localizer"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	core_echo_templates "github.com/fluffy-bunny/fluffycore/echo/templates"
	core_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	BaseHandler struct {
		Container                      di.Container
		Localizer                      func() contracts_localizer.ILocalizer
		ClaimsPrincipal                func() fluffycore_contracts_common.IClaimsPrincipal
		EchoContextAccessor            func() fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor
		IdpServiceServer               func() proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer
		RageUserService                func() proto_oidc_user.IFluffyCoreRageUserServiceServer
		AuthorizationRequestStateStore func() proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer
		ScopedMemoryCache              func() contracts_cache.IScopedMemoryCache
		EmailService                   func() contracts_email.IEmailService
		SessionFactory                 func() contracts_sessions.ISessionFactory
		OIDCSession                    func() contracts_oidc_session.IOIDCSession
		WellknownCookies               func() contracts_cookies.IWellknownCookies
		WellknownCookieNames           func() contracts_cookies.IWellknownCookieNames

		localizer                      contracts_localizer.ILocalizer
		claimsPrincipal                fluffycore_contracts_common.IClaimsPrincipal
		echoContextAccessor            fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor
		idpServiceServer               proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer
		rageUserService                proto_oidc_user.IFluffyCoreRageUserServiceServer
		authorizationRequestStateStore proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer
		scopedMemoryCache              contracts_cache.IScopedMemoryCache
		emailService                   contracts_email.IEmailService
		sessionFactory                 contracts_sessions.ISessionFactory
		oidcSession                    contracts_oidc_session.IOIDCSession
		wellknownCookies               contracts_cookies.IWellknownCookies
		wellknownCookieNames           contracts_cookies.IWellknownCookieNames

		config *contracts_config.Config
	}
)

func NewBaseHandler(container di.Container, config *contracts_config.Config) *BaseHandler {

	obj := &BaseHandler{Container: container, config: config}
	obj.Localizer = obj.getLocalizer
	obj.ClaimsPrincipal = obj.getClaimsPrincipal
	obj.EchoContextAccessor = obj.getEchoContextAccessor
	obj.IdpServiceServer = obj.getIdpServiceServer
	obj.RageUserService = obj.getUserService
	obj.AuthorizationRequestStateStore = obj.getOIDCFlowStore
	obj.ScopedMemoryCache = obj.getScopedMemoryCache
	obj.EmailService = obj.getEmailService
	obj.SessionFactory = obj.getSessionFactory
	obj.OIDCSession = obj.getOIDCSession
	obj.WellknownCookies = obj.getWellknownCookies
	obj.WellknownCookieNames = obj.getWellknownCookieNames

	return obj

}

func (b *BaseHandler) GetManifest(c echo.Context) (*models_api_manifest.Manifest, error) {
	ctx := c.Request().Context()

	idps, err := b.GetIDPs(ctx)
	if err != nil {
		return nil, err
	}
	manifest := &models_api_manifest.Manifest{
		DevelopmentMode: b.config.SystemConfig.DeveloperMode,
	}
	manifest.DisableLocalAccountCreation = b.config.DisableLocalAccountCreation
	manifest.DisableSocialAccounts = b.config.DisableSocialAccounts
	manifest.SocialIdps = make([]models_api_manifest.IDP, 0)
	for _, idp := range idps {
		if idp.Enabled && !idp.Hidden {
			manifest.SocialIdps = append(manifest.SocialIdps, models_api_manifest.IDP{
				Slug: idp.Slug,
			})
		}
	}
	manifest.PasskeyEnabled = false
	if b.config.WebAuthNConfig != nil {
		manifest.PasskeyEnabled = b.config.WebAuthNConfig.Enabled
	}
	// we may have a session in flight and got redirect back here.
	// we may have an external OIDC callback that requires a verify code to continue.
	session, err := b.getSession()
	if err == nil {
		sessionIdI, err := session.Get("session_id")
		if err == nil && sessionIdI != nil {
			sessionId, ok := sessionIdI.(string)
			if ok {
				manifest.SessionId = sessionId
			}
		}
		sessionRequest, err := session.Get("request")
		if err == nil {
			authorizationRequest := sessionRequest.(*proto_oidc_models.AuthorizationRequest)

			if authorizationRequest != nil {
				landingPageI, err := session.Get("landingPage")
				if err == nil && landingPageI != nil {
					landingPage, ok := landingPageI.(*models_api_manifest.LandingPage)
					if ok && landingPage != nil {
						manifest.LandingPage = landingPage
					}
					// get rid of it.
					//session.Set("landingPage", nil)
					//session.Save()
				}
			}
		}

	}
	return manifest, nil
}

func (b *BaseHandler) getWellknownCookies() contracts_cookies.IWellknownCookies {
	if b.wellknownCookies == nil {
		b.wellknownCookies = di.Get[contracts_cookies.IWellknownCookies](b.Container)
	}
	return b.wellknownCookies
}
func (b *BaseHandler) getWellknownCookieNames() contracts_cookies.IWellknownCookieNames {
	if b.wellknownCookieNames == nil {
		b.wellknownCookieNames = di.Get[contracts_cookies.IWellknownCookieNames](b.Container)
	}
	return b.wellknownCookieNames
}
func (b *BaseHandler) getSession() (contracts_sessions.ISession, error) {
	session, err := b.getOIDCSession().GetSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (b *BaseHandler) getOIDCSession() contracts_oidc_session.IOIDCSession {
	if b.oidcSession == nil {
		b.oidcSession = di.Get[contracts_oidc_session.IOIDCSession](b.Container)
	}
	return b.oidcSession
}
func (b *BaseHandler) getSessionFactory() contracts_sessions.ISessionFactory {
	if b.sessionFactory == nil {
		b.sessionFactory = di.Get[contracts_sessions.ISessionFactory](b.Container)
	}
	return b.sessionFactory
}
func (b *BaseHandler) getEmailService() contracts_email.IEmailService {
	if b.emailService == nil {
		b.emailService = di.Get[contracts_email.IEmailService](b.Container)
	}
	return b.emailService
}
func (b *BaseHandler) getScopedMemoryCache() contracts_cache.IScopedMemoryCache {
	if b.scopedMemoryCache == nil {
		b.scopedMemoryCache = di.Get[contracts_cache.IScopedMemoryCache](b.Container)
	}
	return b.scopedMemoryCache
}

func (b *BaseHandler) getLocalizer() contracts_localizer.ILocalizer {
	if b.localizer == nil {
		b.localizer = di.Get[contracts_localizer.ILocalizer](b.Container)
	}
	return b.localizer
}
func (b *BaseHandler) getClaimsPrincipal() fluffycore_contracts_common.IClaimsPrincipal {
	if b.claimsPrincipal == nil {
		b.claimsPrincipal = di.Get[fluffycore_contracts_common.IClaimsPrincipal](b.Container)
	}
	return b.claimsPrincipal
}
func (b *BaseHandler) getEchoContextAccessor() fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor {
	if b.echoContextAccessor == nil {
		b.echoContextAccessor = di.Get[fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor](b.Container)
	}
	return b.echoContextAccessor
}
func (b *BaseHandler) getIdpServiceServer() proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer {
	if b.idpServiceServer == nil {
		b.idpServiceServer = di.Get[proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer](b.Container)
	}
	return b.idpServiceServer
}
func (b *BaseHandler) getUserService() proto_oidc_user.IFluffyCoreRageUserServiceServer {
	if b.rageUserService == nil {
		b.rageUserService = di.Get[proto_oidc_user.IFluffyCoreRageUserServiceServer](b.Container)
	}
	return b.rageUserService
}
func (b *BaseHandler) getOIDCFlowStore() proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer {
	if b.authorizationRequestStateStore == nil {
		b.authorizationRequestStateStore = di.Get[proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer](b.Container)
	}
	return b.authorizationRequestStateStore
}

func (b *BaseHandler) RenderAutoPost(c echo.Context, action string, formData []models.FormParam) error {
	data := map[string]interface{}{
		"form_params": formData,
		"action":      action,
	}
	return b.Render(c, http.StatusFound, "oidc/autopost/index", data)
}

func (b *BaseHandler) Render(c echo.Context, code int, name string, data map[string]interface{}) error {
	localizer := b.Localizer().GetLocalizer()
	data["LocalizeMessage"] = func(key string) string {
		message, _ := localizer.LocalizeMessage(&i18n.Message{ID: key})
		return message
	}
	data["isAuthenticated"] = func() bool {
		if b.ClaimsPrincipal == nil {
			return false
		}
		isAuthenticated := b.ClaimsPrincipal().HasClaimType(core_wellknown.ClaimTypeAuthenticated)
		return isAuthenticated
	}
	data["getUsername"] = func() string {
		claims := b.ClaimsPrincipal().GetClaimsByType("email")
		if len(claims) > 0 {
			return claims[0].Value
		}
		return "Account"
	}
	data["paths"] = wellknown_echo.NewPaths()
	data["username"] = "Account"
	if b.ClaimsPrincipal != nil {
		data["claims"] = b.ClaimsPrincipal().GetClaims()
		claims := b.ClaimsPrincipal().GetClaimsByType("email")
		if len(claims) > 0 {
			data["username"] = claims[0].Value
		}
	}
	type auth struct {
		CSRF string `param:"csrf" query:"csrf" header:"csrf" form:"csrf" json:"csrf" xml:"csrf"`
	}
	csrfValue := c.Get("csrf")
	csrfStr := ""
	if csrfValue != nil {
		if str, ok := csrfValue.(string); ok {
			csrfStr = str
		}
	}
	authArtifacts := &auth{
		CSRF: csrfStr,
	}
	data["security"] = authArtifacts
	data["csrf"] = authArtifacts.CSRF

	return core_echo_templates.Render(c, code, name, data)

}

func (b *BaseHandler) GetIDPs(ctx context.Context) ([]*proto_oidc_models.IDP, error) {
	listIDPResponse, err := b.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			Enabled: &proto_types.BoolFilterExpression{
				Eq: true,
			},
			Hidden: &proto_types.BoolFilterExpression{
				Eq: false,
			},
		},
	})
	if err != nil {

		return nil, err
	}
	return listIDPResponse.IDPs, nil
}

type ProcessFinalAuthenticationStateRequest struct {
	AuthorizationRequest *proto_oidc_models.AuthorizationRequest
	Identity             *proto_oidc_models.OIDCIdentity
	RootPath             string
}

type ProcessFinalAuthenticationStateResponse struct {
	RedirectURI string
}

// ProcessFinalAuthenticationState handles the final state after successful authentication
func (b *BaseHandler) ProcessFinalAuthenticationState(
	ctx context.Context,
	c echo.Context,
	request *ProcessFinalAuthenticationStateRequest,
) (*ProcessFinalAuthenticationStateResponse, error) {
	// Set the auth cookie
	err := b.WellknownCookies().SetAuthCookie(c,
		&contracts_cookies.SetAuthCookieRequest{
			AuthCookie: &contracts_cookies.AuthCookie{
				Identity: &proto_oidc_models.Identity{
					Subject:       request.Identity.Subject,
					Email:         request.Identity.Email,
					EmailVerified: request.Identity.EmailVerified,
				},
			},
		})
	if err != nil {
		return nil, err
	}

	// Get the authorization request state
	getAuthorizationRequestStateResponse, err := b.AuthorizationRequestStateStore().
		GetAuthorizationRequestState(ctx,
			&proto_oidc_flows.GetAuthorizationRequestStateRequest{
				State: request.AuthorizationRequest.State,
			})
	if err != nil {
		return nil, err
	}

	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	authorizationFinal.Identity = request.Identity

	// Store the authorization request state with the code
	_, err = b.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx,
		&proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     authorizationFinal.Request.Code,
			AuthorizationRequestState: authorizationFinal,
		})
	if err != nil {
		return nil, err
	}

	// Delete the old state
	b.AuthorizationRequestStateStore().DeleteAuthorizationRequestState(ctx,
		&proto_oidc_flows.DeleteAuthorizationRequestStateRequest{
			State: request.AuthorizationRequest.State,
		})

	// Store with the state key
	_, err = b.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx,
		&proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			State:                     request.AuthorizationRequest.State,
			AuthorizationRequestState: authorizationFinal,
		})
	if err != nil {
		return nil, err
	}

	// Build the redirect URI
	redirectUri := authorizationFinal.Request.RedirectUri +
		"?code=" + authorizationFinal.Request.Code +
		"&state=" + authorizationFinal.Request.State +
		"&iss=" + request.RootPath

	return &ProcessFinalAuthenticationStateResponse{
		RedirectURI: redirectUri,
	}, nil
}

func (b *BaseHandler) TeleportBackToLoginWithError(c echo.Context, code, msg string) error {
	formParams := []models.FormParam{
		{
			Name:  "error_code",
			Value: code,
		},
		{
			Name:  "error",
			Value: msg,
		},
	}
	return b.RenderAutoPost(c, wellknown_echo.OIDCLoginPath, formParams)

}
func (b *BaseHandler) TeleportToPath(c echo.Context, path string) error {
	formParams := []models.FormParam{}
	return b.RenderAutoPost(c, path, formParams)

}
