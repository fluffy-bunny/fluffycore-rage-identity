package api_user_totp

import (
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_echo_wellknown "github.com/fluffy-bunny/fluffycore/echo/wellknown"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	qrcode "github.com/skip2/go-qrcode"
	gotp "github.com/xlzd/gotp"
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

const (
	InternalError_UserTOTP_001 = "rg-user-totp-001"
	InternalError_UserTOTP_002 = "rg-user-totp-002"
	InternalError_UserTOTP_003 = "rg-user-totp-003"
	InternalError_UserTOTP_004 = "rg-user-totp-004"
	InternalError_UserTOTP_005 = "rg-user-totp-005"
)

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
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
		},
		wellknown_echo.API_UserTOTP,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type TOTPStatusResponse struct {
	Enabled         bool   `json:"enabled"`
	Verified        bool   `json:"verified"`
	Secret          string `json:"secret,omitempty"`          // Only if not verified yet
	ProvisioningURI string `json:"provisioningUri,omitempty"` // Only if not verified yet
	QRCodeBase64    string `json:"qrCodeBase64,omitempty"`    // Only if not verified yet
}

// API GetUserTOTP godoc
// @Summary Get user's TOTP/authenticator status
// @Description Returns the current TOTP configuration and enrollment status
// @Tags account
// @Produce json
// @Success 200 {object} TOTPStatusResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/totp [get]
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
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_UserTOTP_001}, "  ")
	}

	user := getRageUserResponse.User
	response := TOTPStatusResponse{
		Enabled:  user.TOTP != nil && user.TOTP.Enabled,
		Verified: user.TOTP != nil && user.TOTP.Verified,
	}

	// If TOTP exists but not verified yet, include enrollment info
	if user.TOTP != nil && !user.TOTP.Verified && !fluffycore_utils.IsEmptyOrNil(user.TOTP.Secret) {
		response.Secret = user.TOTP.Secret

		// Generate provisioning URI
		otp := gotp.NewDefaultTOTP(user.TOTP.Secret)
		issuerName := "FluffyCore"
		if s.config.TOTP != nil && !fluffycore_utils.IsEmptyOrNil(s.config.TOTP.IssuerName) {
			issuerName = s.config.TOTP.IssuerName
		}

		email := user.RootIdentity.Email
		if fluffycore_utils.IsEmptyOrNil(email) {
			email = subject
		}

		response.ProvisioningURI = otp.ProvisioningUri(email, issuerName)

		// Generate QR code
		pngBytes, err := qrcode.Encode(response.ProvisioningURI, qrcode.Medium, 256)
		if err == nil {
			response.QRCodeBase64 = base64.StdEncoding.EncodeToString(pngBytes)
		}
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
