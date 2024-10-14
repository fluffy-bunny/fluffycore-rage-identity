package personal_information

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
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

		wellknownCookies            contracts_cookies.IWellknownCookies
		passwordHasher              contracts_identity.IPasswordHasher
		fluffyCoreUserServiceServer proto_external_user.IFluffyCoreUserServiceServer
	}
)

const (
	templateName = "account/personal_information/index"
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
			contracts_handler.POST,
		},
		wellknown_echo.PersonalInformationPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type PersonalInformationGetRequest struct {
	Action    string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
	ReturnUrl string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
}

type PersonalInformationPostRequest struct {
	Action      string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
	ReturnUrl   string `param:"returnUrl" query:"returnUrl" form:"returnUrl" json:"returnUrl" xml:"returnUrl"`
	GivenName   string `param:"given_name" query:"given_name" form:"given_name" json:"given_name" xml:"given_name"`
	FamilyName  string `param:"family_name" query:"family_name" form:"family_name" json:"family_name" xml:"family_name"`
	PhoneNumber string `param:"phone_number" query:"phone_number" form:"phone_number" json:"phone_number" xml:"phone_number"`
}

func (s *service) validatePersonalInformationGetRequest(model *PersonalInformationGetRequest) error {
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
	model := &PersonalInformationGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Debug().Interface("model", model).Msg("model")
	err := s.validatePersonalInformationGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validatePersonalInformationGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}
	user, err := s.getUser(c)
	if err != nil {
		return c.Redirect(http.StatusFound, "/error")
	}
	phoneNumber := ""
	if !fluffycore_utils.IsEmptyOrNil(user.Profile.PhoneNumbers) {
		phoneNumber = user.Profile.PhoneNumbers[0].Number
	}

	err = s.Render(c, http.StatusOK, templateName,
		map[string]interface{}{
			"action":       model.Action,
			"returnUrl":    model.ReturnUrl,
			"formAction":   wellknown_echo.PersonalInformationPath,
			"errors":       []string{},
			"email":        user.RageUser.RootIdentity.Email,
			"given_name":   user.Profile.GivenName,
			"family_name":  user.Profile.FamilyName,
			"phone_number": phoneNumber,
		})
	return err
}

func (s *service) validatePersonalInformationPostRequest(request *PersonalInformationPostRequest) ([]string, error) {
	localizer := s.Localizer().GetLocalizer()

	errors := make([]string, 0)
	if fluffycore_utils.IsEmptyOrNil(request.Action) {
		ee := utils.LocalizeWithInterperlate(localizer, "action.is.empty", nil)
		errors = append(errors, ee)
	}
	if fluffycore_utils.IsEmptyOrNil(request.ReturnUrl) {
		ee := utils.LocalizeWithInterperlate(localizer, "returnurl.is.empty", nil)
		errors = append(errors, ee)
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
	model := &PersonalInformationPostRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Debug().Interface("model", model).Msg("model")

	user, err := s.getUser(c)
	if err != nil {
		return c.Redirect(http.StatusFound, "/error")
	}
	phoneNumber := ""
	if !fluffycore_utils.IsEmptyOrNil(user.Profile.PhoneNumbers) {
		phoneNumber = user.Profile.PhoneNumbers[0].Number
	}

	errors, err := s.validatePersonalInformationPostRequest(model)
	doErrorReturn := func() error {
		err = s.Render(c, http.StatusOK, templateName,
			map[string]interface{}{
				"action":       model.Action,
				"returnUrl":    model.ReturnUrl,
				"formAction":   wellknown_echo.PersonalInformationPath,
				"errors":       errors,
				"email":        user.RageUser.RootIdentity.Email,
				"given_name":   user.Profile.GivenName,
				"family_name":  user.Profile.FamilyName,
				"phone_number": phoneNumber,
			})
		return err
	}
	if err != nil {
		return doErrorReturn()
	}

	user.Profile.GivenName = model.GivenName
	user.Profile.FamilyName = model.FamilyName
	user.Profile.PhoneNumbers = []*proto_types.PhoneNumberDTO{
		{
			Id:     "0",
			Number: model.PhoneNumber,
		},
	}
	_, err = s.fluffyCoreUserServiceServer.UpdateUser(ctx,
		&proto_external_user.UpdateUserRequest{
			User: &proto_external_models.ExampleUserUpdate{
				Id: user.Id,

				Profile: &proto_external_models.ProfileUpdate{
					GivenName:  &wrapperspb.StringValue{Value: model.GivenName},
					FamilyName: &wrapperspb.StringValue{Value: model.FamilyName},
					PhoneNumbers: []*proto_types.PhoneNumberDTOUpdate{
						{
							Id:     "0",
							Number: &wrapperspb.StringValue{Value: model.PhoneNumber},
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("ListUser")
		return c.Redirect(http.StatusFound, "/Error")
	}

	// send the email
	_, err = s.EmailService().SendSimpleEmail(ctx,
		&contracts_email.SendSimpleEmailRequest{
			ToEmail:   user.RageUser.RootIdentity.Email,
			SubjectId: "password.reset.changed.subject",
			BodyId:    "password.reset.changed.message",
		})
	if err != nil {
		log.Error().Err(err).Msg("SendEmail")
		return c.Redirect(http.StatusFound, "/error")
	}
	return c.Redirect(http.StatusFound, model.ReturnUrl)

}

func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodPost:
		return s.DoPost(c)
	case http.MethodGet:
		return s.DoGet(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
