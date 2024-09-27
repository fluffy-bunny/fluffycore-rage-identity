package api_verify_password_strength

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/password"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
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
		passwordHasher contracts_identity.IPasswordHasher
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

func (s *service) Ctor(
	container di.Container,
	passwordHasher contracts_identity.IPasswordHasher,
) (*service, error) {
	return &service{
		BaseHandler:    services_echo_handlers_base.NewBaseHandler(container),
		passwordHasher: passwordHasher,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_VerifyPasswordStrength,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateVerifyPasswordStrengthRequest(model *password.VerifyPasswordStrengthRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Password) {
		return status.Error(codes.InvalidArgument, "model.Password is nil")
	}

	return nil
}

// API verify-password-strength godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept json
// @Produce json
// @Param		request body		password.VerifyPasswordStrengthRequest	true	"LoginPhaseOneRequest"
// @Success 200 {object} password.VerifyPasswordStrengthResponse
// @Router /api/verify-password-strength [post]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &password.VerifyPasswordStrengthRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if err := s.validateVerifyPasswordStrengthRequest(model); err != nil {
		log.Error().Err(err).Msg("validateVerifyUsernameRequest")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	err := s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	response := &password.VerifyPasswordStrengthResponse{
		Valid: true,
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
