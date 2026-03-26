package home

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/htmx/components"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v5"
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
		webAuthNConfig   *contracts_webauthn.WebAuthNConfig
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	oidcSession contracts_oidc_session.IOIDCSession,
	webAuthNConfig *contracts_webauthn.WebAuthNConfig,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		wellknownCookies: wellknownCookies,
		passwordHasher:   passwordHasher,
		oidcSession:      oidcSession,
		webAuthNConfig:   webAuthNConfig,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXHomePath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
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

func (s *service) renderHome(c *echo.Context, code int, errors []string, email string) error {
	ctx := c.Request().Context()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)
	manifest, _ := s.GetManifest(c)
	idps, _ := s.GetIDPs(ctx)

	disableSignup := false
	if manifest != nil {
		disableSignup = manifest.DisableLocalAccountCreation
	}

	return components.RenderNode(c, code, components.HomePartial(components.HomeData{
		RenderContext:   rc,
		Errors:          errors,
		Email:           email,
		SocialIdps:      idps,
		DisableSignup:   disableSignup,
		EnabledWebAuthN: s.webAuthNConfig.Enabled,
	}))
}

func (s *service) DoGet(c *echo.Context) error {
	return s.renderHome(c, http.StatusOK, nil, "")
}

type HomePostRequest struct {
	Email   string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	IDPHint string `param:"idp_hint" query:"idp_hint" form:"idp_hint" json:"idp_hint" xml:"idp_hint"`
}

func (s *service) DoPost(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &HomePostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.renderHome(c, http.StatusBadRequest, []string{"Invalid request"}, "")
	}

	// If an IDP hint was provided, redirect to external IDP
	if fluffycore_utils.IsNotEmptyOrNil(model.IDPHint) {
		c.Response().Header().Set("HX-Redirect", wellknown_echo.ExternalIDPPath+"?idp_hint="+model.IDPHint+"&directive=login")
		return c.NoContent(http.StatusOK)
	}

	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		msg := utils.LocalizeSimple(localizer, "username.is.empty")
		return s.renderHome(c, http.StatusBadRequest, []string{msg}, "")
	}

	model.Email = strings.ToLower(model.Email)
	email, ok := echo_utils.IsValidEmailAddress(model.Email)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "username.not.valid", map[string]string{"username": model.Email})
		return s.renderHome(c, http.StatusBadRequest, []string{msg}, model.Email)
	}

	// Check if domain is claimed by an external IDP
	parts := strings.Split(email, "@")
	domainPart := parts[1]
	listIDPRequest, err := s.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			ClaimedDomains: &proto_types.StringArrayFilterExpression{
				Eq: domainPart,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("ListIDP")
		return s.renderHome(c, http.StatusInternalServerError, []string{err.Error()}, model.Email)
	}
	if len(listIDPRequest.IDPs) > 0 {
		// Domain claimed by external IDP - redirect there
		c.Response().Header().Set("HX-Redirect", wellknown_echo.ExternalIDPPath+"?idp_hint="+listIDPRequest.IDPs[0].Slug+"&directive=login")
		return c.NoContent(http.StatusOK)
	}

	// Check if user exists
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: model.Email,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			// User not found - show signup option
			msg := utils.LocalizeWithInterperlate(localizer, "user.not.found", map[string]string{"username": model.Email})
			return s.renderHome(c, http.StatusOK, []string{msg}, model.Email)
		}
		log.Error().Err(err).Msg("GetRageUser")
		return s.renderHome(c, http.StatusInternalServerError, []string{err.Error()}, model.Email)
	}

	if getRageUserResponse == nil {
		msg := utils.LocalizeWithInterperlate(localizer, "user.not.found", map[string]string{"username": model.Email})
		return s.renderHome(c, http.StatusOK, []string{msg}, model.Email)
	}

	user := getRageUserResponse.User
	hasPasskey := false
	if s.webAuthNConfig.Enabled {
		if user.WebAuthN != nil && fluffycore_utils.IsNotEmptyOrNil(user.WebAuthN.Credentials) {
			hasPasskey = true
		}
	}

	// Set the signin username cookie
	err = s.wellknownCookies.SetSigninUserNameCookie(c,
		&contracts_cookies.SetSigninUserNameCookieRequest{
			Value: &contracts_cookies.SigninUserNameCookie{
				Email:      model.Email,
				HasPasskey: hasPasskey,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetSigninUserNameCookie")
		return s.renderHome(c, http.StatusInternalServerError, []string{err.Error()}, model.Email)
	}

	session, err := s.oidcSession.GetSession()
	if err == nil {
		session.Set("landing_page", &models_api_manifest.LandingPage{
			Page: models_api_manifest.PagePasswordEntry,
		})
		session.Save()
	}

	// Render password partial
	rc := components.NewRenderContext(c, localizer)
	return components.RenderNode(c, http.StatusOK, components.PasswordPartial(components.PasswordData{
		RenderContext: rc,
		Email:         model.Email,
		Errors:        []string{},
	}))
}
