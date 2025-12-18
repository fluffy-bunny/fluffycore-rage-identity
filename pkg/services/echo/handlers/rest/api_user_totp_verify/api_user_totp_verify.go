package api_user_totp_verify

import (
	"net/http"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	gotp "github.com/xlzd/gotp"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
	}
)

var stemService = (*service)(nil)
var _ contracts_handler.IHandler = stemService

const (
	InternalError_TOTPVerify_001 = "rg-totp-verify-001"
	InternalError_TOTPVerify_002 = "rg-totp-verify-002"
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
		wellknown_echo.API_UserTOTPVerify,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type VerifyTOTPRequest struct {
	Code string `json:"code" binding:"required"`
}

// API VerifyTOTP godoc
// @Summary Verify TOTP code and complete enrollment
// @Description Verifies the TOTP code from authenticator app and enables TOTP
// @Tags account
// @Accept json
// @Produce json
// @Param request body VerifyTOTPRequest true "Verification request"
// @Success 200 {object} wellknown_echo.RestSuccessResponse
// @Failure 400 {object} wellknown_echo.RestErrorResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Failure 404 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/totp/verify [post]
func (s *service) Do(c echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()

	// Get authenticated user's subject from claims
	claimsPrincipal := s.ClaimsPrincipal()
	subjectClaims := claimsPrincipal.GetClaimsByType(fluffycore_echo_wellknown.ClaimTypeSubject)
	if fluffycore_utils.IsEmptyOrNil(subjectClaims) {
		return c.JSONPretty(http.StatusUnauthorized, wellknown_echo.RestErrorResponse{Error: "Unauthorized"}, "  ")
	}
	claim := subjectClaims[0]
	if fluffycore_utils.IsEmptyOrNil(claim.Value) {
		return c.JSONPretty(http.StatusUnauthorized, wellknown_echo.RestErrorResponse{Error: "Unauthorized"}, "  ")
	}
	subject := claim.Value

	// Bind request
	req := &VerifyTOTPRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Invalid request"}, "  ")
	}

	if fluffycore_utils.IsEmptyOrNil(req.Code) {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Code is required"}, "  ")
	}

	// Get the user from the store
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: subject,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.JSONPretty(http.StatusNotFound, wellknown_echo.RestErrorResponse{Error: "User not found"}, "  ")
		}
		log.Error().Err(err).Msg("GetRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_TOTPVerify_001}, "  ")
	}

	user := getRageUserResponse.User

	// Check if TOTP is set up
	if user.TOTP == nil || fluffycore_utils.IsEmptyOrNil(user.TOTP.Secret) {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "TOTP not enrolled. Call /api/totp/enroll first."}, "  ")
	}

	// Verify the code
	otp := gotp.NewDefaultTOTP(user.TOTP.Secret)
	valid := otp.Verify(req.Code, time.Now().Unix())
	if !valid {
		return c.JSONPretty(http.StatusBadRequest, wellknown_echo.RestErrorResponse{Error: "Invalid TOTP code"}, "  ")
	}

	// Update user - mark as verified and enabled
	_, err = s.RageUserService().UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: subject,
			},
			TOTP: &proto_oidc_models.TOTPUpdate{
				Enabled:  &wrapperspb.BoolValue{Value: true},
				Verified: &wrapperspb.BoolValue{Value: true},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_TOTPVerify_002}, "  ")
	}

	return c.JSONPretty(http.StatusOK, wellknown_echo.RestSuccessResponse{Message: "TOTP verified and enabled successfully"}, "  ")
}
