package profile

import (
	"fmt"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/internal/contracts/cookies"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/internal/wellknown/echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
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
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, err := memCache.Get("rootIdentity")
	if err != nil {
		log.Error().Err(err).Msg("memCache.Get")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity := cachedItem.(*proto_oidc_models.Identity)
	if rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.Redirect(http.StatusFound, "/error")
	}

	userService := s.UserService()
	// get the user
	getUserResponse, err := userService.GetUser(ctx,
		&proto_oidc_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("userService.GetUser")
		return c.Redirect(http.StatusFound, "/error")
	}
	user := getUserResponse.User
	if user.Profile == nil {
		user.Profile = &proto_oidc_models.Profile{}
	}
	phoneNumber := ""
	if !fluffycore_utils.IsEmptyOrNil(user.Profile.PhoneNumbers) {
		phoneNumber = user.Profile.PhoneNumbers[0].Number
	}
	return s.Render(c, http.StatusOK,
		"account/profile/index",
		map[string]interface{}{
			"displayOnly":  true,
			"formAction":   wellknown_echo.ProfilePath,
			"action":       "pi.edit",
			"email":        rootIdentity.Email,
			"given_name":   user.Profile.GivenName,
			"family_name":  user.Profile.FamilyName,
			"phone_number": phoneNumber,
		})
}
func (s *service) DoPersonalInformationEdit(c echo.Context) error {
	redirectUrl := fmt.Sprintf("%s?action=edit&returnUrl=%s",
		wellknown_echo.PersonalInformationPath,
		wellknown_echo.ProfilePath)
	return c.Redirect(http.StatusFound, redirectUrl)
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
		case "pi.edit":
			return s.DoPersonalInformationEdit(c)
		}
		return s.DoGet(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
