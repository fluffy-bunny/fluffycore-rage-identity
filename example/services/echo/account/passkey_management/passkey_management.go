package passkey_management

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
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

		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
	}
)

const (
	templateName = "account/passkey_management/index"
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	passwordHasher contracts_identity.IPasswordHasher,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies: wellknownCookies,
		passwordHasher:   passwordHasher,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.PasskeyManagementPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type PasskeyManagmentGetRequest struct {
	Action    string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
}

func (s *service) validatePasskeyManagementGetRequest(model *PasskeyManagmentGetRequest) error {
	if fluffycore_utils.IsEmptyOrNil(model.Action) {
		return status.Error(codes.InvalidArgument, "Action is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(model.ReturnUrl) {
		model.ReturnUrl = "/"
	}
	return nil
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &PasskeyManagmentGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Debug().Interface("model", model).Msg("model")
	err := s.validatePasskeyManagementGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validatePasskeyManagementGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}

	err = s.Render(c, http.StatusOK, templateName,
		map[string]interface{}{
			"action":    model.Action,
			"returnUrl": model.ReturnUrl,
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
