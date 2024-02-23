package passwordreset

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	services_handlers_shared "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/shared"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		wellknownCookies contracts_cookies.IWellknownCookies
		passwordHasher   contracts_identity.IPasswordHasher
	}
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
			// do auto post
			//contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.PasswordResetPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type PasswordResetGetRequest struct {
	State string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
}

type PasswordResetPostRequest struct {
	State           string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	Password        string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
	ConfirmPassword string `param:"confirmPassword" query:"confirmPassword" form:"confirmPassword" json:"confirmPassword" xml:"confirmPassword"`
}

func (s *service) validatePasswordResetGetRequest(model *PasswordResetGetRequest) error {
	if fluffycore_utils.IsEmptyOrNil(model.State) {
		return status.Error(codes.InvalidArgument, "State is empty")
	}
	return nil
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &PasswordResetGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("model", model).Msg("model")
	err := s.validatePasswordResetGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validatePasswordResetGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}

	err = s.Render(c, http.StatusOK, "oidc/passwordreset/index",
		map[string]interface{}{
			"state":  model.State,
			"errors": []*services_handlers_shared.Error{},
		})
	return err
}

func (s *service) validatePasswordResetPostRequest(request *PasswordResetPostRequest) ([]*services_handlers_shared.Error, error) {
	errors := make([]*services_handlers_shared.Error, 0)
	if fluffycore_utils.IsEmptyOrNil(request.State) {
		errors = append(errors, services_handlers_shared.NewErrorF("state", "State is empty"))
	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		errors = append(errors, services_handlers_shared.NewErrorF("password", "Password is empty"))
	}
	if fluffycore_utils.IsEmptyOrNil(request.ConfirmPassword) {
		errors = append(errors, services_handlers_shared.NewErrorF("confirmPassword", "ConfirmPassword is empty"))
	}
	if request.Password != request.ConfirmPassword {
		errors = append(errors, services_handlers_shared.NewErrorF("confirmPassword", "Passwords do not match"))
	}
	if len(errors) > 0 {
		return errors, status.Error(codes.InvalidArgument, "validation failed")
	}
	return nil, nil
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	// is the request get or post?
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &PasswordResetPostRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	if fluffycore_utils.IsEmptyOrNil(model.Password) || fluffycore_utils.IsEmptyOrNil(model.ConfirmPassword) {
		return s.DoGet(c)
	}
	errors, err := s.validatePasswordResetPostRequest(model)
	if err != nil {
		return s.Render(c, http.StatusBadRequest, "oidc/passwordreset/index",
			map[string]interface{}{
				"state":  model.State,
				"errors": errors,
			})
	}

	getPasswordResetCookieResponse, err := s.wellknownCookies.GetPasswordResetCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetPasswordResetCookie")
		return c.Redirect(http.StatusFound, "/error")
	}
	if getPasswordResetCookieResponse == nil {
		return c.Redirect(http.StatusFound, "/error")
	}
	if getPasswordResetCookieResponse.PasswordReset == nil {
		return c.Redirect(http.StatusFound, "/error")
	}
	if fluffycore_utils.IsEmptyOrNil(getPasswordResetCookieResponse.PasswordReset.Subject) {
		s.wellknownCookies.DeletePasswordResetCookie(c)
		return c.Redirect(http.StatusFound, "/error")
	}
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx,
		&contracts_identity.HashPasswordRequest{
			Password: model.Password,
		})
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return c.Redirect(http.StatusFound, "/error")
	}

	getUserResponse, err := s.UserService().GetUser(ctx,
		&proto_oidc_user.GetUserRequest{
			Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
		})
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return c.Redirect(http.StatusFound, "/Error")
	}

	_, err = s.UserService().UpdateUser(ctx, &proto_oidc_user.UpdateUserRequest{
		User: &proto_oidc_models.UserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
			},
			Password: &proto_oidc_models.PasswordUpdate{
				Hash: &wrapperspb.StringValue{
					Value: hashPasswordResponse.HashedPassword,
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.Redirect(http.StatusFound, "/error")
	}

	// send the email
	_, err = s.EmailService().SendSimpleEmail(ctx,
		&contracts_email.SendSimpleEmailRequest{
			ToEmail:   getUserResponse.User.RootIdentity.Email,
			SubjectId: "password.reset.changed.subject",
			BodyId:    "password.reset.changed.message",
		})
	if err != nil {
		log.Error().Err(err).Msg("SendEmail")
		return c.Redirect(http.StatusFound, "/error")
	}

	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
		[]models.FormParam{
			{
				Name:  "state",
				Value: model.State,
			},
			{
				Name:  "email",
				Value: getUserResponse.User.RootIdentity.Email,
			},
		})

}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       code            		query     string  true  "code"
// @Success 200 {object} string
// @Router /login [get,post]
func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
