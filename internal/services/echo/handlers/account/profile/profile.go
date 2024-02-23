package profile

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	"github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.ProfilePath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) DoGet(c echo.Context) error {
	return s.Render(c, http.StatusOK, "account/profile/index", map[string]interface{}{})
}
func (s *service) DoPasswordReset(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	getAuthCookieResponse, err := s.wellknownCookies.GetAuthCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthCookie")
		return c.Redirect(http.StatusFound, "/")
	}

	err = s.wellknownCookies.SetPasswordResetCookie(c,
		&contracts_cookies.SetPasswordResetCookieRequest{
			PasswordReset: &contracts_cookies.PasswordReset{
				Subject: getAuthCookieResponse.AuthCookie.Identity.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetPasswordResetCookie")
		return c.Redirect(http.StatusFound, "/error")
	}
	return s.RenderAutoPost(c, wellknown_echo.PasswordResetPath,
		[]models.FormParam{
			{
				// need to pass this as a requirment
				Name:  "state",
				Value: "profile.password-reset",
			},
			{
				Name:  "returnUrl",
				Value: wellknown_echo.ProfilePath,
			},
		})
}

type ProfileActionPost struct {
	Action string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
}

func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		model := &ProfileActionPost{}
		if err := c.Bind(model); err != nil {
			log.Error().Err(err).Msg("c.Bind")
			return c.Redirect(http.StatusFound, "/error")
		}
		switch model.Action {
		case "password-reset":
			return s.DoPasswordReset(c)
		}
		return s.DoGet(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
