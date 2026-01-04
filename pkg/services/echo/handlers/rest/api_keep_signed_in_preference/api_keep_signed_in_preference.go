package api_keep_signed_in_preference

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	models_api_preferences "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_preferences"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
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
		wellknown_echo.API_KeepSignedInPreference,
	)
}

const (
	InternalError_KeepSignedInPreference_001 = "rg-keepsignedinpref-001"
	InternalError_KeepSignedInPreference_002 = "rg-keepsignedinpref-002"
)

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateUpdateRequest(model *models_api_preferences.UpdateKeepSignedInPreferenceRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	return nil
}

// API GetKeepSignedInPreference godoc
// @Summary Get keep signed in preference.
// @Description Check if user has set keep signed in preference
// @Tags preferences
// @Accept json
// @Produce json
// @Success 200 {object} models_api_preferences.GetKeepSignedInPreferenceResponse
// @Failure 500 {object} models_api_preferences.ErrorResponse
// @Router /api/keep-signed-in-preference [get]
func (s *service) doGet(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("rootIdentity not found")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.Redirect(http.StatusFound, "/error")
	}
	subject := rootIdentity.Subject

	response := &models_api_preferences.GetKeepSignedInPreferenceResponse{
		DoNotShowAgain: false,
		KeepSignedIn:   false,
	}
	// Check if preference cookie exists
	getKeepSigninPreferencesCookieResponse, err := s.WellknownCookies().
		GetKeepSigninPreferencesCookie(c,
			&contracts_cookies.GetKeepSigninPreferencesCookieRequest{
				Subject: subject,
			})
	if err == nil {
		if getKeepSigninPreferencesCookieResponse != nil {
			response.DoNotShowAgain = getKeepSigninPreferencesCookieResponse.KeepSigninPreferencesCookie.DoNotAskAgain
			response.KeepSignedIn = getKeepSigninPreferencesCookieResponse.KeepSigninPreferencesCookie.KeepSignedIn
		}
	}

	return c.JSON(http.StatusOK, response)
}

// API UpdateKeepSignedInPreference godoc
// @Summary Update keep signed in preference.
// @Description Set or clear keep signed in preference cookie
// @Tags preferences
// @Accept json
// @Produce json
// @Param		request body		models_api_preferences.UpdateKeepSignedInPreferenceRequest	true	"UpdateKeepSignedInPreferenceRequest"
// @Success 200 {object} models_api_preferences.UpdateKeepSignedInPreferenceResponse
// @Failure 500 {object} models_api_preferences.ErrorResponse
// @Router /api/keep-signed-in-preference [post]
func (s *service) doPost(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		log.Error().Msg("rootIdentity not found")
		return c.Redirect(http.StatusFound, "/error")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return c.Redirect(http.StatusFound, "/error")
	}

	model := &models_api_preferences.UpdateKeepSignedInPreferenceRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	if err := s.validateUpdateRequest(model); err != nil {
		log.Error().Err(err).Msg("validateUpdateRequest")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	// Try to get auth from SSO cookie first, then fall back to Auth cookie
	var subject string
	getSSOCookieResponse, err := s.WellknownCookies().GetSSOCookie(c)
	if err == nil && getSSOCookieResponse.SSOCookie != nil {
		subject = rootIdentity.Subject
	} else {
		// Fall back to Auth cookie
		getAuthCookieResponse, err := s.WellknownCookies().GetAuthCookie(c)
		if err != nil {
			log.Error().Err(err).Msg("GetAuthCookie")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to get authentication",
			})
		}

		if getAuthCookieResponse.AuthCookie == nil || getAuthCookieResponse.AuthCookie.Identity == nil {
			log.Error().Msg("No auth cookie or identity found")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Not authenticated",
			})
		}

		subject = getAuthCookieResponse.AuthCookie.Identity.Subject
	}

	// Set the preference cookie to skip the keep signed in page
	err = s.WellknownCookies().SetKeepSigninPreferencesCookie(c,
		&contracts_cookies.SetKeepSigninPreferencesCookieRequest{
			Subject: subject,
			KeepSigninPreferencesCookie: &contracts_cookies.KeepSigninPreferencesCookie{
				DoNotAskAgain: model.DoNotShowAgain,
				KeepSignedIn:  model.KeepSignedIn,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("SetKeepSigninPreferencesCookie")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to set preference",
		})
	}

	response := &models_api_preferences.UpdateKeepSignedInPreferenceResponse{
		Success: true,
	}

	return c.JSON(http.StatusOK, response)
}

func (s *service) Do(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return s.doGet(c)
	} else if c.Request().Method == http.MethodPost {
		return s.doPost(c)
	}

	return c.JSON(http.StatusMethodNotAllowed, map[string]string{
		"error": "Method not allowed",
	})
}
