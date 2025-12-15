package api_user_totp_enroll

import (
	"encoding/base64"
	"net/http"

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
	qrcode "github.com/skip2/go-qrcode"
	gotp "github.com/xlzd/gotp"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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
	InternalError_TOTPEnroll_001 = "rg-totp-enroll-001"
	InternalError_TOTPEnroll_002 = "rg-totp-enroll-002"
	InternalError_TOTPEnroll_003 = "rg-totp-enroll-003"
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
			contracts_handler.POST,
		},
		wellknown_echo.API_UserTOTPEnroll,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type TOTPEnrollResponse struct {
	Secret          string `json:"secret"`
	ProvisioningURI string `json:"provisioningUri"`
	QRCodeBase64    string `json:"qrCodeBase64"`
}

// API EnrollTOTP godoc
// @Summary Start TOTP enrollment
// @Description Generates a new TOTP secret and returns QR code for authenticator app
// @Tags account
// @Produce json
// @Success 200 {object} TOTPEnrollResponse
// @Failure 401 {object} wellknown_echo.RestErrorResponse
// @Failure 409 {object} wellknown_echo.RestErrorResponse "TOTP already verified"
// @Failure 500 {object} wellknown_echo.RestErrorResponse
// @Router /api/totp/enroll [post]
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
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_TOTPEnroll_001}, "  ")
	}

	user := getRageUserResponse.User

	// Check if already verified - don't allow re-enrollment
	if user.TOTP != nil && user.TOTP.Verified {
		return c.JSONPretty(http.StatusConflict, wellknown_echo.RestErrorResponse{Error: "TOTP already verified. Disable first to re-enroll."}, "  ")
	}

	// Generate new TOTP secret
	secret := gotp.RandomSecret(16)

	// Update user with new secret (not verified, not enabled yet)
	_, err = s.RageUserService().UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: subject,
			},
			TOTP: &proto_oidc_models.TOTPUpdate{
				Secret:   &wrapperspb.StringValue{Value: secret},
				Enabled:  &wrapperspb.BoolValue{Value: false},
				Verified: &wrapperspb.BoolValue{Value: false},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateRageUser")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_TOTPEnroll_002}, "  ")
	}

	// Generate provisioning URI
	otp := gotp.NewDefaultTOTP(secret)
	issuerName := "FluffyCore"
	if s.config.TOTP != nil && !fluffycore_utils.IsEmptyOrNil(s.config.TOTP.IssuerName) {
		issuerName = s.config.TOTP.IssuerName
	}

	email := user.RootIdentity.Email
	if fluffycore_utils.IsEmptyOrNil(email) {
		email = subject
	}

	provisioningURI := otp.ProvisioningUri(email, issuerName)

	// Generate QR code
	pngBytes, err := qrcode.Encode(provisioningURI, qrcode.Medium, 256)
	if err != nil {
		log.Error().Err(err).Msg("qrcode.Encode")
		return c.JSONPretty(http.StatusInternalServerError, wellknown_echo.RestErrorResponse{Error: InternalError_TOTPEnroll_003}, "  ")
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	response := TOTPEnrollResponse{
		Secret:          secret,
		ProvisioningURI: provisioningURI,
		QRCodeBase64:    qrCodeBase64,
	}

	return c.JSONPretty(http.StatusOK, response, "  ")
}
