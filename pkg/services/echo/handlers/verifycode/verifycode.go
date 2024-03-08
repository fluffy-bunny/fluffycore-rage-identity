package forgotpassword

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
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
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		wellknownCookies: wellknownCookies,
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
		wellknown_echo.VerifyCodePath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type VerifyCodeGetRequest struct {
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Code      string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
}

type VerifyCodePostRequest struct {
	Email     string `param:"email" query:"email" form:"email" json:"email" xml:"email"`
	Code      string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	Directive string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
	Type      string `param:"type" query:"type" form:"type" json:"type" xml:"type"`
}

func (s *service) validateVerifyCodeGetRequest(model *VerifyCodeGetRequest) error {

	if fluffycore_utils.IsEmptyOrNil(model.Email) {
		return status.Error(codes.InvalidArgument, "Email is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Directive) {
		return status.Error(codes.InvalidArgument, "Directive is empty")
	}
	return nil
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &VerifyCodeGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("model", model).Msg("model")
	err := s.validateVerifyCodeGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validateVerifyCodeGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}

	err = s.Render(c, http.StatusOK, "oidc/verifycode/index",
		map[string]interface{}{
			"email":     model.Email,
			"code":      model.Code,
			"directive": model.Directive,
			"errors":    make([]string, 0),
		})
	return err
}

func (s *service) validateVerifyCodePostRequest(request *VerifyCodePostRequest) ([]string, error) {
	localizer := s.Localizer().GetLocalizer()
	var err error
	errors := make([]string, 0)

	if fluffycore_utils.IsEmptyOrNil(request.Email) {
		msg := utils.LocalizeSimple(localizer, "email.is.empty")
		errors = append(errors, msg)
	}
	if fluffycore_utils.IsEmptyOrNil(request.Code) {
		msg := utils.LocalizeSimple(localizer, "code.is.empty")
		errors = append(errors, msg)
	}
	_, ok := echo_utils.IsValidEmailAddress(request.Email)
	if !ok {
		msg := utils.LocalizeWithInterperlate(localizer, "email.is.not.valid", map[string]string{"email": request.Email})
		errors = append(errors, msg)
	}
	if fluffycore_utils.IsEmptyOrNil(request.Directive) {
		msg := utils.LocalizeSimple(localizer, "directive.is.empty")
		errors = append(errors, msg)
	}
	return errors, err
}

func (s *service) DoPost(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	r := c.Request()
	// is the request get or post?
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &VerifyCodePostRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	errors, err := s.validateVerifyCodePostRequest(model)
	if err != nil {
		return s.Render(c, http.StatusBadRequest, "oidc/verifycode/index",
			map[string]interface{}{
				"email":     model.Email,
				"code":      model.Code,
				"directive": model.Directive,
				"errors":    errors,
			})
	}
	if model.Type == "GET" {
		return s.DoGet(c)
	}
	getVerificationCodeCookieResponse, err := s.wellknownCookies.GetVerificationCodeCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("GetVerificationCodeCookie")
		return s.RenderAutoPost(c, wellknown_echo.ForgotPasswordPath,
			[]models.FormParam{

				{
					Name:  "email",
					Value: model.Email,
				},
				{
					Name:  "type",
					Value: "GET",
				},
			})
	}
	verificationCode := getVerificationCodeCookieResponse.VerificationCode
	code := verificationCode.Code

	if code != model.Code {
		return s.Render(c, http.StatusBadRequest, "oidc/verifycode/index",
			map[string]interface{}{
				"email":     model.Email,
				"code":      model.Code,
				"directive": model.Directive,
				"errors": []string{
					utils.LocalizeSimple(localizer, "code.is.invalid"),
				},
			})
	}
	userService := s.RageUserService()

	_, err = userService.UpdateRageUser(ctx, &proto_oidc_user.UpdateRageUserRequest{
		User: &proto_oidc_models.RageUserUpdate{
			RootIdentity: &proto_oidc_models.IdentityUpdate{
				Subject: verificationCode.Subject,
				EmailVerified: &wrapperspb.BoolValue{
					Value: true,
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("UpdateUser")
		return c.Redirect(http.StatusFound, "/error")
	}
	// one time only
	s.wellknownCookies.DeleteVerificationCodeCookie(c)

	redirectURL := "/"
	switch model.Directive {
	case models.PasswordResetDirective:
		err = s.wellknownCookies.SetPasswordResetCookie(c,
			&contracts_cookies.SetPasswordResetCookieRequest{
				PasswordReset: &contracts_cookies.PasswordReset{
					Subject: verificationCode.Subject,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetPasswordResetCookie")
			return c.Redirect(http.StatusFound, "/error")
		}
		return s.RenderAutoPost(c, wellknown_echo.PasswordResetPath,
			[]models.FormParam{})
	case models.VerifyEmailDirective:
		return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
			[]models.FormParam{
				{
					Name:  "email",
					Value: model.Email,
				},
			})
	}

	return c.Redirect(http.StatusFound, redirectURL)

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
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
