package api_verify_username

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/verify_username"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
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
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService

}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_OIDCLogin_001      = "rg-oidclogin-001"
	InternalError_OIDCLogin_002      = "rg-oidclogin-002"
	InternalError_OIDCLogin_003      = "rg-oidclogin-003"
	InternalError_OIDCLogin_004      = "rg-oidclogin-004"
	InternalError_OIDCLogin_005      = "rg-oidclogin-005"
	InternalError_OIDCLogin_006      = "rg-oidclogin-006"
	InternalError_OIDCLogin_007      = "rg-oidclogin-007"
	InternalError_OIDCLogin_008      = "rg-oidclogin-008"
	InternalError_OIDCLogin_009      = "rg-oidclogin-009"
	InternalError_OIDCLogin_010      = "rg-oidclogin-010"
	InternalError_OIDCLogin_011      = "rg-oidclogin-011"
	InternalError_VerifyUsername_099 = "rg-oidclogin-099"
)

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.API_VerifyUsername,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) validateVerifyUsernameRequest(model *verify_username.VerifyUsernameRequest) error {
	if fluffycore_utils.IsNil(model) {
		return status.Error(codes.InvalidArgument, "model is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(model.UserName) {
		return status.Error(codes.InvalidArgument, "model.Username is nil")
	}

	return nil
}

// API Manifest godoc
// @Summary get the login manifest.
// @Description This is the configuration of the server..
// @Tags root
// @Accept json
// @Produce json
// @Success 200 {object} verify_username.VerifyUsernameResponse
// @Router /api/verify-username [post]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &verify_username.VerifyUsernameRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	if err := s.validateVerifyUsernameRequest(model); err != nil {
		log.Error().Err(err).Msg("validateVerifyUsernameRequest")
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: model.UserName,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: err.Error()}, "  ")
	}
	user := getRageUserResponse.User
	passkeyAvailable := false
	if fluffycore_utils.IsNotNil(user.WebAuthN) &&
		fluffycore_utils.IsNotEmptyOrNil(user.WebAuthN.Credentials) {
		passkeyAvailable = true

	}
	response := &verify_username.VerifyUsernameResponse{
		UserName:         model.UserName,
		PasskeyAvailable: passkeyAvailable,
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
