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
	contracts_codeexchange "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/codeexchange"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/identity"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/client"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
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
	genericOIDCCodeExchange contracts_codeexchange.IGenericOIDCCodeExchange) (*service, error) {
	return &service{
		BaseHandler:             services_echo_handlers_base.NewBaseHandler(container),
		clientServiceServer:     clientServiceServer,
		githubCodeExchange:      githubCodeExchange,
		genericOIDCCodeExchange: genericOIDCCodeExchange,
		config:                  config,
		userIdGenerator:         userIdGenerator,
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

	externalOAuth2State, err := s.ExternalOauth2FlowStore().GetExternalOauth2Final(ctx, model.State)
	if err != nil {
		log.Error().Err(err).Msg("GetExternalOauth2Final")
		return c.Redirect(http.StatusTemporaryRedirect, "/login?error=invalid_state")
	}
	parentState := externalOAuth2State.Request.ParentState
	// if we get here we are going to NUKE the cache for this transaction
	defer func() {
		s.ExternalOauth2FlowStore().DeleteExternalOauth2Final(ctx, model.State)
	}()
	oidcFinalState, err := s.OIDCFlowStore().GetAuthorizationFinal(ctx, parentState)
	if err != nil {
		log.Error().Err(err).Msg("GetAuthorizationFinal")
		redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
			wellknown_echo.OIDCLoginPath,
			parentState, models.InternalError)
		return c.Redirect(http.StatusFound, redirectURL)
	}

	getIDPBySlugResponse, err := s.IdpServiceServer().GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: externalOAuth2State.Request.IDPHint,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return c.Redirect(http.StatusFound, "/error")
	}
	var exchangeCodeResponse *contracts_codeexchange.ExchangeCodeResponse
	idp := getIDPBySlugResponse.Idp

	isMetadataBoolSet := func(key string, idp *proto_oidc_models.IDP) bool {
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
	if idp.Protocol != nil {
		log.Info().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			{
				exchangeCodeResponse, err = s.githubCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPHint:      externalOAuth2State.Request.IDPHint,
					ClientID:     externalOAuth2State.Request.ClientID,
					Nonce:        externalOAuth2State.Request.Nonce,
					Code:         model.Code,
					CodeVerifier: externalOAuth2State.Request.CodeChallenge,
				})
				if err != nil {
					log.Error().Err(err).Msg("ExchangeCode")
					return c.Redirect(http.StatusTemporaryRedirect, "/login?error=exchange_code")
				}
			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				exchangeCodeResponse, err = s.genericOIDCCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPHint:  externalOAuth2State.Request.IDPHint,
					ClientID: externalOAuth2State.Request.ClientID,
					Nonce:    externalOAuth2State.Request.Nonce,
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

		externalIdentity := &models.Identity{
			Subject: rawToken.Subject(),
			Email:   email.(string),
			ACR: []string{
				fmt.Sprintf("urn:mastodon:idp:%s", externalOAuth2State.Request.IDPHint),
			},
			AMR: []string{
				models.AMRIdp,
			},
			EmailVerified: emailVerified,
		}
		getUserByEmail := func(email string) (*proto_oidc_models.User, error) {
			// is this user already linked.
			listUserResponse, err := s.UserService().ListUser(ctx, &proto_oidc_user.ListUserRequest{
				Filter: &proto_oidc_user.Filter{
					RootIdentity: &proto_oidc_user.IdentityFilter{
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
		loginLinkedUser := func(user *proto_oidc_models.User) error {

			oidcFinalState.Identity = &models.Identity{
				Subject: user.RootIdentity.Subject,
				Email:   user.RootIdentity.Email,
				ACR: []string{
					fmt.Sprintf("urn:mastodon:idp:%s", externalOAuth2State.Request.IDPHint),
				},
				AMR: []string{
					models.AMRIdp,
				},
			}
			err = s.OIDCFlowStore().StoreAuthorizationFinal(ctx, parentState, oidcFinalState)
			if err != nil {
				log.Error().Err(err).Msg("StoreAuthorizationFinal")
				redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
					wellknown_echo.OIDCLoginPath,
					parentState, models.InternalError)
				return c.Redirect(http.StatusFound, redirectURL)
			}
			// redirect back
			redirectURL := fmt.Sprintf("%s?state=%s&directive=%s", wellknown_echo.OIDCLoginPath, parentState, models.IdentityFound)
			return c.Redirect(http.StatusFound, redirectURL)
		}

		// is this user already linked.
		listUserResponse, err := s.UserService().ListUser(ctx, &proto_oidc_user.ListUserRequest{
			Filter: &proto_oidc_user.Filter{
				LinkedIdentity: &proto_oidc_user.IdentityFilter{
					Subject: &proto_types.IDFilterExpression{
						Eq: rawToken.Subject(),
					},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("ListUser")
			redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
				wellknown_echo.OIDCLoginPath,
				parentState, models.InternalError)
			return c.Redirect(http.StatusFound, redirectURL)
		}

		if len(listUserResponse.Users) > 0 {
			user := listUserResponse.Users[0]
			return loginLinkedUser(user)
		}
		linkUserAndLogin := func(candidateUserID string, externalIdentity *models.Identity) error {
			getUserResponse, err := s.UserService().GetUser(ctx, &proto_oidc_user.GetUserRequest{
				Subject: candidateUserID,
			})
			if err != nil {
				log.Error().Err(err).Msg("GetUser")
				redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
					wellknown_echo.OIDCLoginPath,
					parentState, models.InternalError)
				return c.Redirect(http.StatusFound, redirectURL)
			}
			user := getUserResponse.User
			if user == nil {
				log.Error().Msg("user not found")
				redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
					wellknown_echo.OIDCLoginPath,
					parentState, models.InternalError)
				return c.Redirect(http.StatusFound, redirectURL)
			}
			_, err = s.UserService().LinkUsers(ctx, &proto_oidc_user.LinkUsersRequest{
				RootSubject: candidateUserID,
				ExternalIdentity: &proto_oidc_models.Identity{
					Subject:       externalIdentity.Subject,
					Email:         externalIdentity.Email,
					IdpSlug:       externalOAuth2State.Request.IDPHint,
					EmailVerified: externalIdentity.EmailVerified,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("LinkUsers")
				redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
					wellknown_echo.OIDCLoginPath,
					parentState, models.InternalError)

				return c.Redirect(http.StatusFound, redirectURL)
			}
			return loginLinkedUser(user)
		}
		// not found, redirect to OIDC LoginPage telling the user to do the signup dance
		if externalOAuth2State.Request.Directive == models.LoginDirective {

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
				if !fluffycore_utils.IsEmptyOrNil(oidcFinalState.Request.CandidateUserID) {
					// CandidateUserID hint
					//--------------------------------------------------------------------------------------------
					return linkUserAndLogin(user.RootIdentity.Subject, externalIdentity)
				}

				// is AUTO-ACCOUNT creation enabled for this IDP?
				if isMetadataBoolSet(models.Wellknown_IDP_Metadata_AutoCreate, idp) {
					emailVerified := false
					emailVerificationRequired := isMetadataBoolSet(models.Wellknown_IDP_Metadata_EmailVerificationRequired, idp)
					if emailVerificationRequired {
						if externalIdentity.EmailVerified {
							emailVerified = true
						}
					} else {
						emailVerified = true
					}
					createUserResponse, err := s.UserService().CreateUser(ctx,
						&proto_oidc_user.CreateUserRequest{
							User: &proto_oidc_models.User{
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
											IdpSlug:       externalOAuth2State.Request.IDPHint,
											EmailVerified: externalIdentity.EmailVerified,
										},
									},
								},
							},
						})
					if err != nil {
						log.Error().Err(err).Msg("CreateUser")
						redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
							wellknown_echo.OIDCLoginPath,
							parentState, models.InternalError)
						return c.Redirect(http.StatusFound, redirectURL)
					}
					return loginLinkedUser(createUserResponse.User)
				}
				// we bounce the user back to go through a sigunup flow
				redirectURL := fmt.Sprintf("%s?state=%s&error=%s",
					wellknown_echo.OIDCLoginPath,
					parentState, models.ExternalIDPNotLinked)
				return c.Redirect(http.StatusFound, redirectURL)
			}

		}
		if externalOAuth2State.Request.Directive == models.SignupDirective {
			oidcFinalState.Directive = models.SignupDirective
			oidcFinalState.ExternalIdentity = externalIdentity
			err = s.OIDCFlowStore().StoreAuthorizationFinal(ctx, parentState, oidcFinalState)
			if err != nil {
				log.Error().Err(err).Msg("StoreAuthorizationFinal")
				return c.Redirect(http.StatusTemporaryRedirect, "/login?error=store_authorization_final")
			}
			// we don't store the external identity on a missed match.  User has to go through the trouble of a signup
			redirectURL := fmt.Sprintf("%s?state=%s", wellknown_echo.OIDCLoginPath, parentState)
			return c.Redirect(http.StatusFound, redirectURL)
		}

	}
	return c.Redirect(http.StatusTemporaryRedirect, "/error?state="+parentState)
}
