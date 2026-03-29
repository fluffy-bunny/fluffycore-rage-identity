package password

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	pkg_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	api_profile "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/api_profile"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config         *contracts_config.Config
		userService    proto_external_user.IFluffyCoreUserServiceServer
		passwordHasher contracts_identity.IPasswordHasher
		emailService   contracts_email.IEmailService
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	userService proto_external_user.IFluffyCoreUserServiceServer,
	passwordHasher contracts_identity.IPasswordHasher,
	emailService contracts_email.IEmailService,
) (*service, error) {
	return &service{
		BaseHandler:    services_echo_handlers_base.NewBaseHandler(container, config),
		config:         config,
		userService:    userService,
		passwordHasher: passwordHasher,
		emailService:   emailService,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		wellknown_echo.HTMXManagementPasswordPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	// Non-HTMX GET requests (e.g. browser refresh) need the full shell page
	if c.Request().Method == http.MethodGet && !components.IsHTMXRequest(c) {
		return c.Redirect(http.StatusFound, wellknown_echo.HTMXManagementPath+"?redirect="+c.Request().URL.Path)
	}
	r := c.Request()
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	return c.NoContent(http.StatusNotFound)
}

func (s *service) getProfileData(c *echo.Context) (*api_profile.Profile, error) {
	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	profile := &api_profile.Profile{
		Email:   rootIdentity.Email,
		Subject: rootIdentity.Subject,
	}
	cp := s.ClaimsPrincipal()
	acrClaims := cp.GetClaimsByType("acr")
	for _, claim := range acrClaims {
		if claim.Value == pkg_models.ACRClaimedDomain {
			profile.IsClaimedDomain = true
			break
		}
	}
	return profile, nil
}

func (s *service) DoGet(c *echo.Context) error {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	profile, err := s.getProfileData(c)
	if err != nil {
		return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
			RenderContext: rc,
			Stage:         components.PasswordStageInitial,
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}

	return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
		RenderContext:   rc,
		Stage:           components.PasswordStageInitial,
		Email:           profile.Email,
		IsClaimedDomain: profile.IsClaimedDomain,
	}))
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	action := c.FormValue("action")
	profile, err := s.getProfileData(c)
	if err != nil {
		return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
			RenderContext: rc,
			Stage:         components.PasswordStageInitial,
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}

	switch action {
	case "send-code":
		// Generate hashed verification code (matching the pattern in pkg/services/echo/handlers/htmx/password)
		codeResult, err := echo_utils.GenerateHashedVerificationCode(ctx, s.passwordHasher, 6)
		if err != nil {
			log.Error().Err(err).Msg("GenerateHashedVerificationCode")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageInitial,
				Email:         profile.Email,
				Error:         rc.L("mgmt_something_went_wrong"),
			}))
		}

		// Log code in developer mode for testing
		if s.config.SystemConfig.DeveloperMode {
			log.Info().Str("code", codeResult.PlainCode).Str("email", profile.Email).Msg("DEV MODE: Password reset verification code")
		}

		// Store hashed code in cookie (persists across requests, unlike ScopedMemoryCache)
		plainCode := ""
		if s.config.SystemConfig.DeveloperMode {
			plainCode = codeResult.PlainCode
		}
		err = s.WellknownCookies().SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Email:             profile.Email,
					CodeHash:          codeResult.HashedCode,
					PlainCode:         plainCode,
					Subject:           profile.Subject,
					VerifyCodePurpose: contracts_cookies.VerifyCode_PasswordReset,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageInitial,
				Email:         profile.Email,
				Error:         rc.L("mgmt_something_went_wrong"),
			}))
		}

		// Send email with plain code
		_, err = s.emailService.SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   profile.Email,
				SubjectId: "forgotpassword.email.subject",
				BodyId:    "password.reset.message",
				Data:      map[string]string{"code": codeResult.PlainCode},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendEmail")
			// Still proceed to verify stage even if email fails in dev mode
		}

		devCode := ""
		if s.config.SystemConfig.DeveloperMode {
			devCode = codeResult.PlainCode
		}

		return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
			RenderContext: rc,
			Stage:         components.PasswordStageVerify,
			Email:         profile.Email,
			DevCode:       devCode,
		}))

	case "verify-code":
		submittedCode := c.FormValue("code")

		// Read hashed code from cookie
		vcResp, err := s.WellknownCookies().GetVerificationCodeCookie(c)
		if err != nil || vcResp == nil || vcResp.VerificationCode == nil {
			log.Error().Err(err).Msg("GetVerificationCodeCookie")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageVerify,
				Email:         profile.Email,
				Error:         rc.L("mgmt_invalid_code"),
			}))
		}

		// Verify code using bcrypt comparison
		err = echo_utils.VerifyVerificationCode(ctx, s.passwordHasher, submittedCode, vcResp.VerificationCode.CodeHash)
		if err != nil {
			log.Warn().Str("submitted", submittedCode).Msg("Verification code mismatch")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageVerify,
				Email:         profile.Email,
				Error:         rc.L("mgmt_invalid_code"),
			}))
		}

		// Code verified — delete verification cookie and set password reset cookie
		s.WellknownCookies().DeleteVerificationCodeCookie(c)
		err = s.WellknownCookies().SetPasswordResetCookie(c, &contracts_cookies.SetPasswordResetCookieRequest{
			PasswordReset: &contracts_cookies.PasswordReset{
				Subject: profile.Subject,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("SetPasswordResetCookie")
		}

		return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
			RenderContext: rc,
			Stage:         components.PasswordStageReset,
			Email:         profile.Email,
		}))

	case "reset-password":
		// Verify password reset cookie exists (proves code was verified)
		prResp, err := s.WellknownCookies().GetPasswordResetCookie(c)
		if err != nil || prResp == nil || prResp.PasswordReset == nil {
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageInitial,
				Email:         profile.Email,
				Error:         rc.L("mgmt_unexpected_error"),
			}))
		}

		newPassword := c.FormValue("newPassword")
		confirmPassword := c.FormValue("confirmPassword")

		if newPassword != confirmPassword {
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageReset,
				Email:         profile.Email,
				Error:         rc.L("mgmt_passwords_do_not_match"),
			}))
		}

		// Hash the password
		hashResp, err := s.passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
			Password: newPassword,
		})
		if err != nil {
			log.Error().Err(err).Msg("HashPassword")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageReset,
				Email:         profile.Email,
				Error:         rc.L("mgmt_password_too_weak"),
			}))
		}

		memCacheItem, ok := s.ScopedMemoryCache().Get("rootIdentity")
		if !ok {
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageInitial,
				Error:         rc.L("mgmt_unexpected_error"),
			}))
		}
		rootIdentity := memCacheItem.(*proto_oidc_models.Identity)

		getUserResponse, err := s.userService.GetUser(ctx,
			&proto_external_user.GetUserRequest{
				Subject: rootIdentity.Subject,
			})
		if err != nil {
			log.Error().Err(err).Msg("GetUser")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageReset,
				Email:         profile.Email,
				Error:         rc.L("mgmt_something_went_wrong"),
			}))
		}
		user := getUserResponse.User

		_, err = s.userService.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
			User: &proto_external_models.ExampleUserUpdate{
				Id: user.Id,
				RageUser: &proto_oidc_models.RageUserUpdate{
					Password: &proto_oidc_models.PasswordUpdate{
						Hash: &wrapperspb.StringValue{Value: hashResp.HashedPassword},
					},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("UpdateUser")
			return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
				RenderContext: rc,
				Stage:         components.PasswordStageReset,
				Email:         profile.Email,
				Error:         rc.L("mgmt_something_went_wrong"),
			}))
		}

		// Send password changed notification email
		_, _ = s.emailService.SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   profile.Email,
				SubjectId: "password.reset.changed.subject",
				BodyId:    "password.reset.changed.message",
			})

		// Clean up cookies
		s.WellknownCookies().DeletePasswordResetCookie(c)
		s.WellknownCookies().DeleteVerificationCodeCookie(c)

		return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
			RenderContext: rc,
			Stage:         components.PasswordStageSuccess,
			Email:         profile.Email,
		}))
	}

	return components.RenderNode(c, http.StatusOK, components.PasswordPage(&components.PasswordPageData{
		RenderContext: rc,
		Stage:         components.PasswordStageInitial,
		Email:         profile.Email,
	}))
}
