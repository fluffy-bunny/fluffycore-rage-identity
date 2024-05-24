package callback

/*
reference: https://github.com/go-oauth2/oauth2/blob/master/example/client/client.go
*/
import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_codeexchange "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/codeexchange"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	jwxt "github.com/lestrrat-go/jwx/v2/jwt"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		clientServiceServer     proto_oidc_client.IFluffyCoreClientServiceServer
		githubCodeExchange      contracts_codeexchange.IGithubCodeExchange
		genericOIDCCodeExchange contracts_codeexchange.IGenericOIDCCodeExchange
		config                  *contracts_config.Config
		userIdGenerator         contracts_identity.IUserIdGenerator
		wellknownCookies        contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_Callback_001 = "rg-callback-001"
	InternalError_Callback_002 = "rg-callback-002"
	InternalError_Callback_003 = "rg-callback-003"
	InternalError_Callback_004 = "rg-callback-004"
	InternalError_Callback_005 = "rg-callback-005"
	InternalError_Callback_006 = "rg-callback-006"
	InternalError_Callback_007 = "rg-callback-007"
	InternalError_Callback_008 = "rg-callback-008"
	InternalError_Callback_009 = "rg-callback-009"
	InternalError_Callback_010 = "rg-callback-010"
	InternalError_Callback_011 = "rg-callback-011"
	InternalError_Callback_099 = "rg-callback-099" // 99 is a bind problem
)

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
	clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer,
	githubCodeExchange contracts_codeexchange.IGithubCodeExchange,
	userIdGenerator contracts_identity.IUserIdGenerator,
	wellknownCookies contracts_cookies.IWellknownCookies,
	genericOIDCCodeExchange contracts_codeexchange.IGenericOIDCCodeExchange) (*service, error) {
	return &service{
		BaseHandler:             services_echo_handlers_base.NewBaseHandler(container),
		clientServiceServer:     clientServiceServer,
		githubCodeExchange:      githubCodeExchange,
		genericOIDCCodeExchange: genericOIDCCodeExchange,
		config:                  config,
		userIdGenerator:         userIdGenerator,
		wellknownCookies:        wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder, callbackPath string) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		callbackPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		//clientauthorization.AuthenticateOAuth2Client(),
	}
}

type CallbackRequest struct {
	Code  string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	State string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       code    				query     string  true  "code requested"
// @Param       state            		query     string  true  "state requested"
// @Success 200 {object} string
// @Router /oauth2/callback [get]
func (s *service) Do(c echo.Context) error {
	localizer := s.Localizer().GetLocalizer()

	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &CallbackRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return s.TeleportBackToLogin(c, InternalError_Callback_099)
	}
	log = log.With().Interface("model", model).Logger()
	var idp *proto_oidc_models.IDP
	getExternalOauth2CookieResponse, err := s.wellknownCookies.GetExternalOauth2Cookie(c,
		&contracts_cookies.GetExternalOauth2CookieRequest{
			State: model.State,
		})

	if err != nil {
		log.Error().Err(err).Msg("GetExternalOauth2Final")
		return c.Redirect(http.StatusTemporaryRedirect, "/login?error=invalid_state")
	}
	externalOauth2State := getExternalOauth2CookieResponse.ExternalOAuth2State
	parentState := externalOauth2State.Request.ParentState
	// if we get here we are going to NUKE the cache for this transaction
	s.wellknownCookies.DeleteExternalOauth2Cookie(c,
		&contracts_cookies.DeleteExternalOauth2CookieRequest{
			State: model.State,
		})

	doLoginBounceBack := func() error {
		formParams := []models.FormParam{
			{
				Name:  "state",
				Value: parentState,
			},
			{
				Name:  "directive",
				Value: models.IdentityFound,
			},
		}
		return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath, formParams)

	}
	doEmailVerification := func(user *proto_oidc_models.RageUser, directive string) error {
		verificationCode := echo_utils.GenerateRandomAlphaNumericString(6)
		err = s.wellknownCookies.SetVerificationCodeCookie(c,
			&contracts_cookies.SetVerificationCodeCookieRequest{
				VerificationCode: &contracts_cookies.VerificationCode{
					Subject: user.RootIdentity.Subject,
					Email:   user.RootIdentity.Email,
					Code:    verificationCode,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetVerificationCodeCookie")
			return s.TeleportBackToLogin(c, InternalError_Callback_007)
		}
		_, err = s.EmailService().SendSimpleEmail(ctx,
			&contracts_email.SendSimpleEmailRequest{
				ToEmail:   user.RootIdentity.Email,
				SubjectId: "email.verification.subject",
				BodyId:    "email.verification..message",
				Data: map[string]string{
					"code": verificationCode,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SendSimpleEmail")
			return s.TeleportBackToLogin(c, InternalError_Callback_009)
		}
		formParams := []models.FormParam{
			{
				Name:  "state",
				Value: parentState,
			},
			{
				Name:  "email",
				Value: user.RootIdentity.Email,
			},
			{
				Name:  "directive",
				Value: directive,
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
	getAuthorizationRequestStateResponse, err := s.AuthorizationRequestStateStore().GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
		State: parentState,
	})
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return s.TeleportBackToLogin(c, InternalError_Callback_001)
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState

	getIDPBySlugResponse, err := s.IdpServiceServer().GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: externalOauth2State.Request.IdpHint,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return s.TeleportBackToLogin(c, InternalError_Callback_008)
	}
	var exchangeCodeResponse *contracts_codeexchange.ExchangeCodeResponse
	idp = getIDPBySlugResponse.Idp

	if idp.Protocol != nil {
		log.Info().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			{
				exchangeCodeResponse, err = s.githubCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPHint:      externalOauth2State.Request.IdpHint,
					ClientID:     externalOauth2State.Request.ClientId,
					Nonce:        externalOauth2State.Request.Nonce,
					Code:         model.Code,
					CodeVerifier: externalOauth2State.Request.CodeChallenge,
				})
				if err != nil {
					log.Error().Err(err).Msg("ExchangeCode")
					return c.Redirect(http.StatusTemporaryRedirect, "/login?error=exchange_code")
				}
			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				exchangeCodeResponse, err = s.genericOIDCCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPHint:  externalOauth2State.Request.IdpHint,
					ClientID: externalOauth2State.Request.ClientId,
					Nonce:    externalOauth2State.Request.Nonce,
					Code:     model.Code,
					//			CodeVerifier: finalCache.Request.CodeChallenge,
				})
				if err != nil {
					log.Error().Err(err).Msg("ExchangeCode")
					return c.Redirect(http.StatusTemporaryRedirect, "/login?error=exchange_code")
				}

			}
		}
	}
	if exchangeCodeResponse != nil && !fluffycore_utils.IsEmptyOrNil(exchangeCodeResponse.IdToken) {
		// now we do the link dance
		parseOptions := []jwxt.ParseOption{
			jwxt.WithVerify(false),
			jwxt.WithAcceptableSkew(time.Minute * 5),
		}
		rawToken, err := jwxt.ParseString(exchangeCodeResponse.IdToken, parseOptions...)
		if err != nil {
			log.Error().Err(err).Msg("ParseString")
			return c.Redirect(http.StatusTemporaryRedirect, "/login?error=parse_id_token")
		}
		email, ok := rawToken.Get("email")
		if !ok {
			log.Error().Msg("email not found")
			return c.Redirect(http.StatusTemporaryRedirect, "/login?error=email_not_found")
		}
		emailVerified := false
		emailVerifiedC, ok := rawToken.Get("email_verified")
		if ok {
			bval, ok := emailVerifiedC.(bool)
			if ok {
				emailVerified = bval
			}
		}

		externalIdentity := &proto_oidc_models.OIDCIdentity{
			Subject: rawToken.Subject(),
			Email:   email.(string),
			Acr: []string{
				fmt.Sprintf("urn:rage:idp:%s", externalOauth2State.Request.IdpHint),
			},
			Amr: []string{
				models.AMRIdp,
			},
			EmailVerified: emailVerified,
		}

		getUserByEmail := func(email string) (*proto_oidc_models.RageUser, error) {
			// is this user already linked.
			getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
				&proto_oidc_user.GetRageUserRequest{
					By: &proto_oidc_user.GetRageUserRequest_Email{
						Email: strings.ToLower(email),
					},
				})

			if err != nil {
				log.Error().Err(err).Msg("ListUser")
				return nil, err
			}
			if getRageUserResponse != nil {
				return getRageUserResponse.User, nil
			}
			return nil, status.Error(codes.NotFound, "user not found")

		}
		loginLinkedUser := func(user *proto_oidc_models.RageUser, directive string) error {
			if idp.MultiFactorRequired || !user.RootIdentity.EmailVerified {
				return doEmailVerification(user, directive)
			}
			authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
				Subject: user.RootIdentity.Subject,
				Email:   user.RootIdentity.Email,
				Acr: []string{
					fmt.Sprintf("urn:rage:idp:%s", externalOauth2State.Request.IdpHint),
				},
				Amr: []string{
					models.AMRIdp,
				},
			}
			_, err = s.AuthorizationRequestStateStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
				State:                     parentState,
				AuthorizationRequestState: authorizationFinal,
			})
			if err != nil {
				log.Error().Err(err).Msg("StoreAuthorizationRequestState")
				return s.TeleportBackToLogin(c, InternalError_Callback_002)
			}
			// redirect back
			return doLoginBounceBack()
		}

		// is this user already linked.
		getRageUserResponse, err := s.RageUserService().GetRageUser(ctx,
			&proto_oidc_user.GetRageUserRequest{
				By: &proto_oidc_user.GetRageUserRequest_ExternalIdentity{
					ExternalIdentity: &proto_oidc_models.Identity{
						Subject: externalIdentity.Subject,
						IdpSlug: externalOauth2State.Request.IdpHint,
					},
				},
			})

		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				err = nil
			} else {
				log.Error().Err(err).Msg("GetUser")
				return s.TeleportBackToLogin(c, InternalError_Callback_003)
			}
		}

		if getRageUserResponse != nil {
			user := getRageUserResponse.User
			return loginLinkedUser(user, models.MFA_VerifyEmailDirective)
		}

		linkUser := func(candidateUserID string, externalIdentity *proto_oidc_models.OIDCIdentity) (*proto_oidc_models.RageUser, error) {
			getUserResponse, err := s.RageUserService().GetRageUser(ctx,
				&proto_oidc_user.GetRageUserRequest{
					By: &proto_oidc_user.GetRageUserRequest_Subject{
						Subject: candidateUserID,
					},
				})
			if err != nil {
				log.Error().Err(err).Msg("GetUser")
				return nil, err
			}
			user := getUserResponse.User
			if user == nil {
				log.Error().Msg("user not found")
				return nil, err
			}
			_, err = s.RageUserService().LinkRageUser(ctx, &proto_oidc_user.LinkRageUserRequest{
				RootSubject: candidateUserID,
				ExternalIdentity: &proto_oidc_models.Identity{
					Subject:       externalIdentity.Subject,
					Email:         externalIdentity.Email,
					IdpSlug:       externalOauth2State.Request.IdpHint,
					EmailVerified: externalIdentity.EmailVerified,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("LinkUsers")
				return nil, err
			}
			return user, nil
		}
		linkUserAndLogin := func(candidateUserID string, externalIdentity *proto_oidc_models.OIDCIdentity) error {
			user, err := linkUser(candidateUserID, externalIdentity)
			if err != nil {
				log.Error().Err(err).Msg("LinkUsers")
				return s.TeleportBackToLogin(c, InternalError_Callback_004)
			}
			return loginLinkedUser(user, models.MFA_VerifyEmailDirective)
		}
		doAutoCreateUser := func() (*proto_oidc_models.RageUser, error) {
			emailVerified := false
			emailVerificationRequired := idp.EmailVerificationRequired
			if emailVerificationRequired {
				emailVerified = false
				/*
					the external IDP may say its been verified, but we don't trust it, we want our own verification
					if externalIdentity.EmailVerified {
						emailVerified = true
					}
				*/
			} else {
				emailVerified = true
			}
			createUserResponse, err := s.RageUserService().CreateRageUser(ctx,
				&proto_oidc_user.CreateRageUserRequest{
					User: &proto_oidc_models.RageUser{
						RootIdentity: &proto_oidc_models.Identity{
							Subject:       s.userIdGenerator.GenerateUserId(),
							IdpSlug:       models.WellknownIdpRoot,
							Email:         externalIdentity.Email,
							EmailVerified: emailVerified,
						},
						LinkedIdentities: &proto_oidc_models.LinkedIdentities{
							Identities: []*proto_oidc_models.Identity{
								{
									Subject:       externalIdentity.Subject,
									Email:         externalIdentity.Email,
									IdpSlug:       externalOauth2State.Request.IdpHint,
									EmailVerified: externalIdentity.EmailVerified,
								},
							},
						},
					},
				})
			if err != nil {
				log.Error().Err(err).Msg("CreateUser")
				return nil, err
			}
			return createUserResponse.User, nil
		}
		// in both cases we do an auto link if we get a perfect hit.
		// a perfect email match beats out a candidate user
		//--------------------------------------------------------------------------------------------
		// auto link if we get an email hit
		user, err := getUserByEmail(externalIdentity.Email)
		if err == nil && user != nil {
			// Perfect email match
			//--------------------------------------------------------------------------------------------
			return linkUserAndLogin(user.RootIdentity.Subject, externalIdentity)
		}

		switch externalOauth2State.Request.Directive {
		case models.LoginDirective:
			// do we have a candidate user to link to?
			if fluffycore_utils.IsNotEmptyOrNil(authorizationFinal.Request.CandidateUserId) {
				// CandidateUserID hint
				//--------------------------------------------------------------------------------------------
				return linkUserAndLogin(authorizationFinal.Request.CandidateUserId, externalIdentity)
			}

			// is AUTO-ACCOUNT creation enabled for this IDP?
			if idp.AutoCreate {
				user, err := doAutoCreateUser()
				if err != nil {
					log.Error().Err(err).Msg("doAutoCreateUser")
					return s.TeleportBackToLogin(c, InternalError_Callback_005)
				}
				return loginLinkedUser(user, models.MFA_VerifyEmailDirective)

			}
			// we bounce the user back to go through a sigunup flow
			msg := utils.LocalizeWithInterperlate(localizer, "username.not.found", map[string]string{"username": externalIdentity.Email})
			return s.TeleportBackToLogin(c, msg)
		case models.SignupDirective:
			user, err := doAutoCreateUser()
			if err != nil {
				log.Error().Err(err).Msg("doAutoCreateUser")
				return s.TeleportBackToLogin(c, InternalError_Callback_006)
			}
			emailVerificationRequired := idp.EmailVerificationRequired
			if !emailVerificationRequired {
				return loginLinkedUser(user, models.VerifyEmailDirective)
			}
			return doEmailVerification(user, models.VerifyEmailDirective)
		}
	}
	return s.TeleportBackToLogin(c, InternalError_Callback_011)
}

// IsMetadataBoolSet checks if the key is set in the metadata and if the value is a boolean
func IsMetadataBoolSet(key string, idp *proto_oidc_models.IDP) bool {
	if fluffycore_utils.IsEmptyOrNil(idp.Metadata) {
		return false
	}
	v, ok := idp.Metadata[key]
	if ok {
		// convert string to boolean
		bVal, err := strconv.ParseBool(v)
		if err == nil {
			return bVal
		}
	}
	return false
}
