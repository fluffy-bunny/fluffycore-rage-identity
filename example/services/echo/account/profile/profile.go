package profile

import (
	"fmt"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		wellknownCookies contracts_cookies.IWellknownCookies
		userService      proto_external_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	userService proto_external_user.IFluffyCoreUserServiceServer,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		wellknownCookies: wellknownCookies,
		userService:      userService,
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
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("rootIdentity not found in cache")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.Redirect(http.StatusFound, "/error")
	}

	// get the user
	getUserResponse, err := s.userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("userService.GetUser")
		return c.Redirect(http.StatusFound, "/error")
	}
	user := getUserResponse.User
	if user.Profile == nil {
		user.Profile = &proto_external_models.Profile{}
	}
	phoneNumber := ""
	if fluffycore_utils.IsNotEmptyOrNil(user.Profile.PhoneNumbers) {
		phoneNumber = user.Profile.PhoneNumbers[0].Number
	}
	return s.Render(c, http.StatusOK,
		"account/profile/index",
		map[string]interface{}{
			"displayOnly":  true,
			"formAction":   wellknown_echo.ProfilePath,
			"action":       "pi-edit",
			"email":        rootIdentity.Email,
			"given_name":   user.Profile.GivenName,
			"family_name":  user.Profile.FamilyName,
			"phone_number": phoneNumber,
			"user":         user,
		})
}

func (s *service) DoPasskeyManagment(c echo.Context) error {
	redirectUrl := fmt.Sprintf("%s?action=edit&returnUrl=%s",
		wellknown_echo.PasskeyManagementPath,
		wellknown_echo.ProfilePath)
	return c.Redirect(http.StatusFound, redirectUrl)
}
func (s *service) DoTOTPManagment(c echo.Context) error {
	redirectUrl := fmt.Sprintf("%s?action=edit&returnUrl=%s",
		wellknown_echo.TOTPPath,
		wellknown_echo.ProfilePath)
	return c.Redirect(http.StatusFound, redirectUrl)
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
		case "totp-management":
			return s.DoTOTPManagment(c)
		case "pi-edit":
			return s.DoPersonalInformationEdit(c)
		case "passkeys":
			return s.DoPasskeyManagment(c)

		}
		return s.DoGet(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
