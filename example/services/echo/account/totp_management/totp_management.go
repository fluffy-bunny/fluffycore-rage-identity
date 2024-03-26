package totp_management

import (
	"net/http"

	"encoding/base64"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
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

		wellknownCookies            contracts_cookies.IWellknownCookies
		passwordHasher              contracts_identity.IPasswordHasher
		fluffyCoreUserServiceServer proto_external_user.IFluffyCoreUserServiceServer
	}
)

const (
	templateName = "account/totp_management/index"
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
	fluffyCoreUserServiceServer proto_external_user.IFluffyCoreUserServiceServer,
) (*service, error) {
	return &service{
		BaseHandler:                 services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies:            wellknownCookies,
		passwordHasher:              passwordHasher,
		fluffyCoreUserServiceServer: fluffyCoreUserServiceServer,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.TOTPPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type TOTPManagmentGetRequest struct {
	Action    string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
}

func (s *service) validateTOTPManagementGetRequest(model *TOTPManagmentGetRequest) error {
	if fluffycore_utils.IsEmptyOrNil(model.Action) {
		return status.Error(codes.InvalidArgument, "Action is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(model.ReturnUrl) {
		model.ReturnUrl = "/"
	}
	return nil
}
func (s *service) getUser(c echo.Context) (*proto_external_models.ExampleUser, error) {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	memCache := s.ScopedMemoryCache()
	cachedItem, err := memCache.Get("rootIdentity")
	if err != nil {
		log.Error().Err(err).Msg("memCache.Get")
		return nil, err
	}
	rootIdentity := cachedItem.(*proto_oidc_models.Identity)
	if rootIdentity == nil {
		log.Error().Msg("rootIdentity is nil")
		return nil, status.Error(codes.NotFound, "rootIdentity is nil")
	}
	userService := s.fluffyCoreUserServiceServer
	// get the user
	getUserResponse, err := userService.GetUser(ctx,
		&proto_external_user.GetUserRequest{
			Subject: rootIdentity.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("userService.GetUser")
		return nil, err
	}
	if getUserResponse.User.Profile == nil {
		getUserResponse.User.Profile = &proto_external_models.Profile{}
	}
	return getUserResponse.User, nil

}
func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	model := &TOTPManagmentGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("model", model).Msg("model")
	err := s.validateTOTPManagementGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validatePasskeyManagementGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}

	user, err := s.getUser(c)
	if err != nil {
		return c.Redirect(http.StatusFound, "/error")
	}

	rageUser := user.RageUser
	totpSecret := rageUser.TOTP.Secret
	otp := gotp.NewDefaultTOTP(totpSecret)

	provisioningUri := otp.ProvisioningUri(rageUser.RootIdentity.Email, "rage.identity")
	var pngB []byte
	pngB, _ = qrcode.Encode(provisioningUri, qrcode.Medium, 256)
	base64Str := base64.StdEncoding.EncodeToString(pngB)

	err = s.Render(c, http.StatusOK, templateName,
		map[string]interface{}{
			"action":    model.Action,
			"returnUrl": model.ReturnUrl,
			"pngQRCode": base64Str,
		})
	return err
}

func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
