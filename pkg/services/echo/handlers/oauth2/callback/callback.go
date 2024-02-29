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
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
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
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OAuth2CallbackPath,
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
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &CallbackRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log = log.With().Interface("model", model).Logger()

	getExternalOauth2FinalResponse, err := s.ExternalOauth2FlowStore().GetExternalOauth2Final(ctx, &proto_oidc_flows.GetExternalOauth2FinalRequest{
		State: model.State,
	})
	if err != nil {
		log.Error().Err(err).Msg("GetExternalOauth2Final")
		return c.Redirect(http.StatusTemporaryRedirect, "/login?error=invalid_state")
	}
	externalOauth2Final := getExternalOauth2FinalResponse.ExternalOauth2Final
	parentState := externalOauth2Final.Request.ParentState
	// if we get here we are going to NUKE the cache for this transaction
	defer func() {
		s.ExternalOauth2FlowStore().DeleteExternalOauth2Final(ctx, &proto_oidc_flows.DeleteExternalOauth2FinalRequest{
			State: model.State,
		})
	}()
	doInternalErrorPost := func() error {
		formParams := []models.FormParam{
			{
				Name:  "state",
				Value: parentState,
			},
			{
				Name:  "error",
				Value: models.InternalError,
			},
		}
		return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath, formParams)

	}

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
	doEmailVerification := func(user *proto_oidc_models.RageUser) error {
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
			return c.Redirect(http.StatusFound, "/error")
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
			return c.Redirect(http.StatusFound, "/error")
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
	getAuthorizationRequestStateResponse, err := s.OIDCFlowStore().GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
		State: parentState,
	})
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationRequestState")
		return doInternalErrorPost()
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState

	getIDPBySlugResponse, err := s.IdpServiceServer().GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: externalOauth2Final.Request.IdpHint,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return c.Redirect(http.StatusFound, "/error")
	}
	var exchangeCodeResponse *contracts_codeexchange.ExchangeCodeResponse
	idp := getIDPBySlugResponse.Idp

	if idp.Protocol != nil {
		log.Info().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			{
				exchangeCodeResponse, err = s.githubCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPHint:      externalOauth2Final.Request.IdpHint,
					ClientID:     externalOauth2Final.Request.ClientId,
					Nonce:        externalOauth2Final.Request.Nonce,
					Code:         model.Code,
					CodeVerifier: externalOauth2Final.Request.CodeChallenge,
				})
				if err != nil {
					log.Error().Err(err).Msg("ExchangeCode")
					return c.Redirect(http.StatusTemporaryRedirect, "/login?error=exchange_code")
				}
			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				exchangeCodeResponse, err = s.genericOIDCCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPHint:  externalOauth2Final.Request.IdpHint,
					ClientID: externalOauth2Final.Request.ClientId,
					Nonce:    externalOauth2Final.Request.Nonce,
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
				fmt.Sprintf("urn:mastodon:idp:%s", externalOauth2Final.Request.IdpHint),
			},
			Amr: []string{
				models.AMRIdp,
			},
			EmailVerified: emailVerified,
		}

		getUserByEmail := func(email string) (*proto_oidc_models.RageUser, error) {
			// is this user already linked.
			listUserResponse, err := s.RageUserService().ListRageUser(ctx, &proto_oidc_user.ListRageUserRequest{
				Filter: &proto_oidc_models.RageUserFilter{
					RootIdentity: &proto_oidc_models.IdentityFilter{
						Email: &proto_types.StringFilterExpression{
							Eq: strings.ToLower(email),
						},
					},
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("ListUser")
				return nil, err
			}
			if len(listUserResponse.Users) > 0 {
				return listUserResponse.Users[0], nil
			}
			return nil, status.Error(codes.NotFound, "user not found")

		}
		loginLinkedUser := func(user *proto_oidc_models.RageUser) error {
			if !user.RootIdentity.EmailVerified {
				return doEmailVerification(user)
			}
			authorizationFinal.Identity = &proto_oidc_models.OIDCIdentity{
				Subject: user.RootIdentity.Subject,
				Email:   user.RootIdentity.Email,
				Acr: []string{
					fmt.Sprintf("urn:mastodon:idp:%s", externalOauth2Final.Request.IdpHint),
				},
				Amr: []string{
					models.AMRIdp,
				},
			}
			_, err = s.OIDCFlowStore().StoreAuthorizationRequestState(ctx, &proto_oidc_flows.StoreAuthorizationRequestStateRequest{
				State:                     parentState,
				AuthorizationRequestState: authorizationFinal,
			})
			if err != nil {
				log.Error().Err(err).Msg("StoreAuthorizationRequestState")
				return doInternalErrorPost()
			}
			// redirect back
			return doLoginBounceBack()
		}

		// is this user already linked.
		listUserResponse, err := s.RageUserService().ListRageUser(ctx, &proto_oidc_user.ListRageUserRequest{
			Filter: &proto_oidc_models.RageUserFilter{
				LinkedIdentity: &proto_oidc_models.IdentityFilter{
					Subject: &proto_types.IDFilterExpression{
						Eq: rawToken.Subject(),
					},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("ListUser")
			return doInternalErrorPost()
		}

		if len(listUserResponse.Users) > 0 {
			user := listUserResponse.Users[0]
			return loginLinkedUser(user)
		}

		linkUser := func(candidateUserID string, externalIdentity *proto_oidc_models.OIDCIdentity) (*proto_oidc_models.RageUser, error) {
			getUserResponse, err := s.RageUserService().GetRageUser(ctx, &proto_oidc_user.GetRageUserRequest{
				Subject: candidateUserID,
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
			_, err = s.RageUserService().LinkRageUsers(ctx, &proto_oidc_user.LinkRageUsersRequest{
				RootSubject: candidateUserID,
				ExternalIdentity: &proto_oidc_models.Identity{
					Subject:       externalIdentity.Subject,
					Email:         externalIdentity.Email,
					IdpSlug:       externalOauth2Final.Request.IdpHint,
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
				return doInternalErrorPost()
			}
			return loginLinkedUser(user)
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
									IdpSlug:       externalOauth2Final.Request.IdpHint,
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
		// not found, redirect to OIDC LoginPage telling the user to do the signup dance
		if externalOauth2Final.Request.Directive == models.LoginDirective {

			// a perfect email match beats out a candidate user
			//--------------------------------------------------------------------------------------------

			// auto link if we get an email hit
			user, err := getUserByEmail(externalIdentity.Email)
			if err == nil && user != nil {
				// Perfect email match
				//--------------------------------------------------------------------------------------------
				return linkUserAndLogin(user.RootIdentity.Subject, externalIdentity)
			} else {
				// do we have a candidate user to link to?
				if !fluffycore_utils.IsEmptyOrNil(authorizationFinal.Request.CandidateUserId) {
					// CandidateUserID hint
					//--------------------------------------------------------------------------------------------
					return linkUserAndLogin(authorizationFinal.Request.CandidateUserId, externalIdentity)
				}

				// is AUTO-ACCOUNT creation enabled for this IDP?
				if idp.AutoCreate {
					user, err := doAutoCreateUser()
					if err != nil {
						log.Error().Err(err).Msg("doAutoCreateUser")
						return doInternalErrorPost()
					}
					return loginLinkedUser(user)

				}
				// we bounce the user back to go through a sigunup flow
				return doInternalErrorPost()
			}

		}

		if externalOauth2Final.Request.Directive == models.SignupDirective {

			user, err := doAutoCreateUser()
			if err != nil {
				log.Error().Err(err).Msg("doAutoCreateUser")
				return doInternalErrorPost()
			}
			emailVerificationRequired := idp.EmailVerificationRequired
			if !emailVerificationRequired {
				return loginLinkedUser(user)
			}
			return doEmailVerification(user)

		}

	}
	return c.Redirect(http.StatusTemporaryRedirect, "/error?state="+parentState)
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
