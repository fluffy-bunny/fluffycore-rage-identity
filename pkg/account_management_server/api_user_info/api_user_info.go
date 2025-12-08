package api_user_info

import (
	"net/http"
	"strconv"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	models_api_login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config          *contracts_config.Config
		claimsPrincipal fluffycore_contracts_common.IClaimsPrincipal
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	claimsPrincipal fluffycore_contracts_common.IClaimsPrincipal,
) (*service, error) {
	return &service{
		BaseHandler:     services_echo_handlers_base.NewBaseHandler(container, config),
		config:          config,
		claimsPrincipal: claimsPrincipal,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.API_IsAuthorized,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// API UserInfo godoc
// @Summary Get current user information.
// @Description Returns the current authenticated user's information from the auth session.
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} login_models.UserInfoResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Router /api/is-authorized [get]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Get subject from claims (middleware ensures user is authenticated)
	subjectClaims := s.claimsPrincipal.GetClaimsByType(fluffycore_echo_wellknown.ClaimTypeSubject)
	if len(subjectClaims) == 0 {
		log.Error().Msg("No subject claim found")
		return c.JSONPretty(http.StatusUnauthorized, wellknown_echo.RestErrorResponse{Error: "Not authenticated"}, "  ")
	}

	// Get email from claims
	emailClaims := s.claimsPrincipal.GetClaimsByType("email")
	email := ""
	if len(emailClaims) > 0 {
		email = emailClaims[0].Value
	}

	// Get email_verified from claims
	emailVerifiedClaims := s.claimsPrincipal.GetClaimsByType("email_verified")
	emailVerified := false
	if len(emailVerifiedClaims) > 0 {
		emailVerified, _ = strconv.ParseBool(emailVerifiedClaims[0].Value)
	}

	// Return user information
	response := &models_api_login_models.UserInfoResponse{
		Subject:       subjectClaims[0].Value,
		Email:         email,
		EmailVerified: emailVerified,
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
