package passwordreset

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
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
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
}

type PasswordResetPostRequest struct {
	ReturnUrl       string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
	Password        string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
	ConfirmPassword string `param:"confirmPassword" query:"confirmPassword" form:"confirmPassword" json:"confirmPassword" xml:"confirmPassword"`
}

func (s *service) validatePasswordResetGetRequest(model *PasswordResetGetRequest) error {

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
			"returnUrl": model.ReturnUrl,
			"errors":    []string{},
		})
	return err
}

func (s *service) validatePasswordResetPostRequest(request *PasswordResetPostRequest) ([]string, error) {
	localizer := s.Localizer().GetLocalizer()

	errors := make([]string, 0)

	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		msg := utils.LocalizeSimple(localizer, "password.is.empty")
		errors = append(errors, msg)
	}
	if fluffycore_utils.IsEmptyOrNil(request.ConfirmPassword) {
		msg := utils.LocalizeSimple(localizer, "confirm_password.is.empty")
		errors = append(errors, msg)
	}
	if request.Password != request.ConfirmPassword {
		msg := utils.LocalizeSimple(localizer, "passwords.do.not.match")
		errors = append(errors, msg)
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
	doErrorReturn := func() error {
		return s.Render(c, http.StatusBadRequest, "oidc/passwordreset/index",
			map[string]interface{}{
				"returnUrl": model.ReturnUrl,
				"errors":    errors,
			})
	}
	if err != nil {
		return doErrorReturn()
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

	getUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return c.Redirect(http.StatusFound, "/Error")
	}

	_, err = s.RageUserService().UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
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
	if !fluffycore_utils.IsEmptyOrNil(model.ReturnUrl) {
		return c.Redirect(http.StatusFound, model.ReturnUrl)
	}
	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
		[]models.FormParam{

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
