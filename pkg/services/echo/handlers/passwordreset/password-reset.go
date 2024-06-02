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

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_PasswordReset_001 = "rg-password-reset-001"
	InternalError_PasswordReset_002 = "rg-password-reset-002"
	InternalError_PasswordReset_003 = "rg-password-reset-003"
	InternalError_PasswordReset_004 = "rg-password-reset-004"
	InternalError_PasswordReset_005 = "rg-password-reset-005"
	InternalError_PasswordReset_006 = "rg-password-reset-006"
	InternalError_PasswordReset_007 = "rg-password-reset-007"
	InternalError_PasswordReset_008 = "rg-password-reset-008"
	InternalError_PasswordReset_009 = "rg-password-reset-009"
	InternalError_PasswordReset_010 = "rg-password-reset-010"
	InternalError_PasswordReset_011 = "rg-password-reset-011"

	InternalError_PasswordReset_099 = "rg-password-reset-099"
)

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
	Action          string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
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
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_099)
	}
	log.Info().Interface("model", model).Msg("model")
	err := s.validatePasswordResetGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validatePasswordResetGetRequest")
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_002)
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
		log.Error().Err(err).Msg("c.Bind")
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_099)
	}
	log.Info().Interface("model", model).Msg("model")

	if model.Action == "cancel" {
		return s.TeleportToPath(c, wellknown_echo.OIDCLoginPath)
	}
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
	localizer := s.Localizer().GetLocalizer()

	getPasswordResetCookieResponse, err := s.wellknownCookies.GetPasswordResetCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetPasswordResetCookie")
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_003)
	}
	if getPasswordResetCookieResponse == nil {
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_004)
	}
	if getPasswordResetCookieResponse.PasswordReset == nil {
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_005)
	}
	if fluffycore_utils.IsEmptyOrNil(getPasswordResetCookieResponse.PasswordReset.Subject) {
		s.wellknownCookies.DeletePasswordResetCookie(c)
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_006)
	}
	getUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Subject{
				Subject: getPasswordResetCookieResponse.PasswordReset.Subject,
			},
		})
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_008)
	}

	// do password acceptablity check
	err = s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("IsAcceptablePassword")
		msg := utils.LocalizeSimple(localizer, "password.is.not.acceptable")
		errors = append(errors, msg)
		return doErrorReturn()
	}
	// hash the password
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx,
		&contracts_identity.HashPasswordRequest{
			Password: model.Password,
		})
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_007)
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
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_009)
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
		return s.TeleportBackToLogin(c, InternalError_PasswordReset_010)
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
