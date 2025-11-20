package signup

import (
	"net/http"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
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
		config           *contracts_config.Config
		passwordHasher   contracts_identity.IPasswordHasher
		wellknownCookies contracts_cookies.IWellknownCookies
		userIdGenerator  contracts_identity.IUserIdGenerator
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_Signup_001 = "rg-signup-001"
	InternalError_Signup_002 = "rg-signup-002"
	InternalError_Signup_003 = "rg-signup-003"
	InternalError_Signup_004 = "rg-signup-004"
	InternalError_Signup_005 = "rg-signup-005"
	InternalError_Signup_006 = "rg-signup-006"
	InternalError_Signup_007 = "rg-signup-007"
	InternalError_Signup_008 = "rg-signup-008"
	InternalError_Signup_009 = "rg-signup-009"
	InternalError_Signup_010 = "rg-signup-010"
	InternalError_Signup_011 = "rg-signup-011"
	InternalError_Signup_099 = "rg-signup-099"
)

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	passwordHasher contracts_identity.IPasswordHasher,
	wellknownCookies contracts_cookies.IWellknownCookies,
	userService proto_oidc_user.IFluffyCoreRageUserServiceServer,
	userIdGenerator contracts_identity.IUserIdGenerator,
) (*service, error) {

	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		config:           config,
		passwordHasher:   passwordHasher,
		wellknownCookies: wellknownCookies,
		userIdGenerator:  userIdGenerator,
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
		wellknown_echo.SignupPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type SignupGetRequest struct {
}
type ExternalIDPAuthRequest struct {
	IDPHint string `param:"idp_hint" query:"idp_hint" form:"idp_hint" json:"idp_hint" xml:"idp_hint"`
}
type SignupPostRequest struct {
	UserName string `param:"username" query:"username" form:"username" json:"username" xml:"username"`
	Password string `param:"password" query:"password" form:"password" json:"password" xml:"password"`
	Type     string `param:"type" query:"type" form:"type" json:"type" xml:"type"`
	Action   string `param:"action" query:"action" form:"action" json:"action" xml:"action"`
}

func (s *service) DoGet(c echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &SignupGetRequest{}
	if err := c.Bind(model); err != nil {
		return s.TeleportBackToLoginWithError(c, InternalError_Signup_099, InternalError_Signup_099)
	}
	log.Debug().Interface("model", model).Msg("model")

	type row struct {
		Key   string
		Value string
	}
	idps, err := s.GetIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return s.TeleportBackToLoginWithError(c, InternalError_Signup_001, InternalError_Signup_001)
	}

	var rows []row
	//	rows = append(rows, row{Key: "code", Value: model.Code})

	return s.Render(c, http.StatusOK, "oidc/signup/index",
		map[string]interface{}{
			"errors":    rows,
			"idps":      idps,
			"directive": models.SignupDirective,
		})
}

func (s *service) validateSignupPostRequest(request *SignupPostRequest) ([]string, error) {
	localizer := s.Localizer().GetLocalizer()

	var err error
	errors := make([]string, 0)

	if fluffycore_utils.IsEmptyOrNil(request.UserName) {
		msg := utils.LocalizeSimple(localizer, "username.is.empty")
		errors = append(errors, msg)
	} else {
		_, ok := echo_utils.IsValidEmailAddress(request.UserName)
		if !ok {
			msg := utils.LocalizeWithInterperlate(localizer, "username.is.not.valid", map[string]string{"username": request.UserName})
			errors = append(errors, msg)
		}
	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		msg := utils.LocalizeSimple(localizer, "password.is.empty")
		errors = append(errors, msg)
	}

	return errors, err
}

func (s *service) DoPost(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()

	idps, err := s.GetIDPs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("getIDPs")
		return s.TeleportBackToLoginWithError(c, InternalError_Signup_002, InternalError_Signup_002)
	}
	doError := func(errors []string) error {
		return s.Render(c, http.StatusBadRequest, "oidc/signup/index",
			map[string]interface{}{
				"errors":    errors,
				"idps":      idps,
				"directive": models.SignupDirective,
			})

	}

	// is the request get or post?
	model := &SignupPostRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("Bind")
		return s.TeleportBackToLoginWithError(c, InternalError_Signup_099, InternalError_Signup_099)
	}
	log.Debug().Interface("model", model).Msg("model")
	if model.Type == "GET" {
		return s.DoGet(c)
	}
	if model.Action == "cancel" {
		return s.TeleportToPath(c, wellknown_echo.OIDCLoginPath)
	}
	errors, err := s.validateSignupPostRequest(model)
	if err != nil {
		return err
	}
	if len(errors) > 0 {
		return doError(errors)
	}
	model.UserName = strings.ToLower(model.UserName)
	// get the domain from the email
	parts := strings.Split(model.UserName, "@")
	domainPart := parts[1]
	// first lets see if this domain has been claimed.
	listIDPRequest, err := s.IdpServiceServer().ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			ClaimedDomain: &proto_types.StringArrayFilterExpression{
				Eq: domainPart,
			},
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("ListIDP")
		errors = append(errors, err.Error())
		return doError(errors)
	}
	if len(listIDPRequest.Idps) > 0 {
		// an idp has claimed this domain.
		// post to the externalIDP
		return s.RenderAutoPost(c, wellknown_echo.ExternalIDPPath,
			[]models.FormParam{

				{
					Name:  "idp_hint",
					Value: listIDPRequest.Idps[0].Slug,
				},
				{
					Name:  "directive",
					Value: models.LoginDirective,
				},
			})
	}
	// does the user exist.
	getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_Email{
				Email: strings.ToLower(model.UserName),
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			err = nil
		} else {
			log.Error().Err(err).Msg("GetRageUser")
			return s.TeleportBackToLoginWithError(c, InternalError_Signup_003, InternalError_Signup_003)
		}

	}
	if getRageUserResponse != nil {
		return doError([]string{
			utils.LocalizeWithInterperlate(localizer, "username.already.exists", map[string]string{"username": model.UserName}),
		})
	}
	subjectId := s.userIdGenerator.GenerateUserId()
	user := &proto_oidc_models.RageUser{
		RootIdentity: &proto_oidc_models.Identity{
			Subject:       subjectId,
			Email:         model.UserName,
			IdpSlug:       models.RootIdp,
			EmailVerified: false,
		},
		State: proto_oidc_models.RageUserState_USER_STATE_PENDING,
	}

	//  check password strength
	err = s.passwordHasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		return doError([]string{
			utils.LocalizeWithInterperlate(localizer, "password.is.not.acceptable",
				map[string]string{"username": model.UserName}),
		})
	}
	hashPasswordResponse, err := s.passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
		Password: model.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("GeneratePasswordHash")
		return s.TeleportBackToLoginWithError(c, InternalError_Signup_004, InternalError_Signup_004)
	}
	user.Password = &proto_oidc_models.Password{
		Hash: hashPasswordResponse.HashedPassword,
	}

	_, err = s.RageUserService().CreateRageUser(ctx, &proto_oidc_user.CreateRageUserRequest{
		User: user,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateUser")
		return s.TeleportBackToLoginWithError(c, InternalError_Signup_005, InternalError_Signup_005)
	}
	if s.config.EmailVerificationRequired {
		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:             model.UserName,
					CodeHash:          verificationCode,
					Subject:           user.RootIdentity.Subject,
					VerifyCodePurpose: contracts_cookies.VerifyCode_EmailVerification,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return s.TeleportBackToLoginWithError(c, InternalError_Signup_006, InternalError_Signup_006)
		}
		_, err = s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   model.UserName,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification.message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendSimpleEmail")
			return s.TeleportBackToLoginWithError(c, InternalError_Signup_007, InternalError_Signup_007)
		}
		formParams := []models.FormParam{

			{
				Name:  "email",
				Value: model.UserName,
			},
			{
				Name:  "directive",
				Value: models.VerifyEmailDirective,
			},
			{
				Name:  "type",
				Value: "GET",
			},
		}

		if s.config.SystemConfig.DeveloperMode {
			formParams = append(formParams, models.FormParam{
				Name:  "code",
				Value: verificationCode,
			})

		}
		return s.RenderAutoPost(c, wellknown_echo.VerifyCodePath, formParams)

	}

	return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
		[]models.FormParam{})

}

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
