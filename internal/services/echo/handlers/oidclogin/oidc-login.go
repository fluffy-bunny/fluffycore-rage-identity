package oidclogin

import (
	"fmt"
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	services_handlers_shared "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/shared"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config           *contracts_config.Config
		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
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
	State     string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Error     string `param:"error" query:"error" form:"error" json:"error" xml:"error"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
}
type ExternalIDPAuthRequest struct {
	IDPHint string `param:"idp_hint" query:"idp_hint" form:"idp_hint" json:"idp_hint" xml:"idp_hint"`
}
type LoginPostRequest struct {
	State    string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	UserName string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
}

type row struct {
	Key   string
	Value string
}

func (s *service) handleIdentityFound(c echo.Context, state string) error {
	r := c.Request()
	// is the request get or post?
	rootPath := s.config.OIDCConfig.BaseUrl
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	mm, err := s.OIDCFlowStore().GetAuthorizationFinal(ctx, state)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationFinal")
		// redirect to error page
		redirectUrl := fmt.Sprintf("%s?state=%s&error=%s", wellknown_echo.OIDCLoginPath, state, models.InternalError)
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	if mm.Identity == nil {
		redirectUrl := fmt.Sprintf("%s?state=%s&error=%s", wellknown_echo.OIDCLoginPath, state, models.InternalError)
		return c.Redirect(http.StatusFound, redirectUrl)
	}

	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: mm.Identity,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}

	err = s.OIDCFlowStore().StoreAuthorizationFinal(ctx, mm.Request.Code, mm)
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationFinal")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	s.OIDCFlowStore().DeleteAuthorizationFinal(ctx, state)

	// redirect to the client with the code.
	redirectUri := mm.Request.RedirectURI +
		"?code=" + mm.Request.Code +
		"&state=" + mm.Request.State +
		"&iss=" + rootPath
	return c.Redirect(http.StatusFound, redirectUri)

}
func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginGetRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	var rows []row
	rows = append(rows, row{Key: "state", Value: model.State})

	switch model.Directive {
	case models.IdentityFound:
		return s.handleIdentityFound(c, model.State)

	}
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		rows = append(rows, row{Key: "error", Value: err.Error()})
	}

	return s.Render(c, http.StatusOK, "oidc/oidclogin/index",
		map[string]interface{}{
			"defs":      rows,
			"idps":      idps,
			"state":     model.State,
			"email":     model.Email,
			"directive": models.LoginDirective,
		})
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LoginPostRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")
	if fluffycore_utils.IsEmptyOrNil(model.UserName) &&
		fluffycore_utils.IsEmptyOrNil(model.Password) {
		return s.DoGet(c)
	}
	model.UserName = strings.ToLower(model.UserName)

	var errors []*services_handlers_shared.Error
	errors = append(errors, services_handlers_shared.NewErrorF("state", model.State))
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		errors = append(errors, services_handlers_shared.NewErrorF("error", err.Error()))
	}

	// does the user exist.
	listUserResponse, err := s.UserService().ListUser(ctx, &proto_oidc_user.ListUserRequest{
		Filter: &proto_oidc_user.Filter{
			RootIdentity: &proto_oidc_user.IdentityFilter{
				Email: &proto_types.StringFilterExpression{
					Eq: model.UserName,
				},
			},
		},
	})
	if len(listUserResponse.Users) == 0 {
		errors = append(errors, services_handlers_shared.NewErrorF("username", "username:%s not found", model.UserName))
		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     model.State,
				"idps":      idps,
				"defs":      errors,
				"directive": models.LoginDirective,
			})

	}
	if err != nil {
		log.Warn().Err(err).Msg("ListUser")
		errors = append(errors, services_handlers_shared.NewErrorF("error", err.Error()))

		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     model.State,
				"idps":      idps,
				"defs":      errors,
				"directive": models.LoginDirective,
			})
	}
	user := listUserResponse.Users[0]
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
			return c.Redirect(http.StatusFound, "/error")
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
				Name:  "state",
				Value: model.State,
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
				Value: verificationCode,
			})
		}
		return s.RenderAutoPost(c, wellknown_echo.VerifyCodePath, formParams)
	}

	if user.Password == nil {
		errors = append(errors, services_handlers_shared.NewErrorF("username", "username:%s does not have a password", model.UserName))

		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     model.State,
				"idps":      idps,
				"defs":      errors,
				"directive": models.LoginDirective,
			})
	}

	err = s.passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
		Password:       model.Password,
		HashedPassword: user.Password.Hash,
	})
	if err != nil {
		log.Warn().Err(err).Msg("ComparePasswordHash")
		errors = append(errors, services_handlers_shared.NewErrorF("password", "password is invalid"))
		return s.Render(c, http.StatusBadRequest, "oidc/oidclogin/index",
			map[string]interface{}{
				"state":     model.State,
				"idps":      idps,
				"defs":      errors,
				"directive": models.LoginDirective,
			})
	}
	err = s.wellknownCookies.SetAuthCookie(c, &contracts_cookies.SetAuthCookieRequest{
		AuthCookie: &contracts_cookies.AuthCookie{
			Identity: &models.Identity{
				Subject:       user.RootIdentity.Subject,
				Email:         user.RootIdentity.Email,
				EmailVerified: user.RootIdentity.EmailVerified,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("SetAuthCookie")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}

	mm, err := s.OIDCFlowStore().GetAuthorizationFinal(ctx, model.State)
	if err != nil {
		log.Warn().Err(err).Msg("GetAuthorizationFinal")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	mm.Identity = &models.Identity{
		Subject: user.RootIdentity.Subject,
		Email:   user.RootIdentity.Email,
		ACR: []string{
			models.ACRPassword,
			models.ACRIdpRoot,
		},
		AMR: []string{
			models.AMRPassword,
			// always true, as we are the root idp
			models.AMRIdp,
		},
	}

	// "urn:mastodon:idp:google", "urn:mastodon:idp:spacex", "urn:mastodon:idp:github-enterprise", etc.
	// "urn:mastodon:password", "urn:mastodon:2fa", "urn:mastodon:email", etc.
	// we are done with the state now.  Lets map it to the code so it can be looked up by the client.
	err = s.OIDCFlowStore().StoreAuthorizationFinal(ctx, mm.Request.Code, mm)
	if err != nil {
		log.Warn().Err(err).Msg("StoreAuthorizationFinal")
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	s.OIDCFlowStore().DeleteAuthorizationFinal(ctx, model.State)

	err = s.OIDCFlowStore().StoreAuthorizationFinal(ctx, model.State, mm)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	// redirect to the client with the code.
	redirectUri := mm.Request.RedirectURI +
		"?code=" + mm.Request.Code +
		"&state=" + mm.Request.State +
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
