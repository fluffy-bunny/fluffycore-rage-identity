package api

import (
	"net/http"
	"strings"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oauth2factory"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/external_idp"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		oauth2Factory    contracts_oauth2factory.IOAuth2Factory
		config           *contracts_config.Config
		oidcSession      contracts_oidc_session.IOIDCSession
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

const (
	// make sure only one is shown.  This is an internal error code to point the developer to the code that is failing
	InternalError_ExternalIDP_001 = "rg-external-idp-001"
	InternalError_ExternalIDP_002 = "rg-external-idp-002"
	InternalError_ExternalIDP_003 = "rg-external-idp-003"
	InternalError_ExternalIDP_004 = "rg-external-idp-004"
	InternalError_ExternalIDP_005 = "rg-external-idp-005"
	InternalError_ExternalIDP_006 = "rg-external-idp-006"
	InternalError_ExternalIDP_007 = "rg-external-idp-007"
	InternalError_ExternalIDP_008 = "rg-external-idp-008"
	InternalError_ExternalIDP_009 = "rg-external-idp-009"
	InternalError_ExternalIDP_010 = "rg-external-idp-010"
	InternalError_ExternalIDP_011 = "rg-external-idp-011"
	InternalError_ExternalIDP_099 = "rg-external-idp-099" // 99 is a bind problem
)

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	oauth2Factory contracts_oauth2factory.IOAuth2Factory,
	oidcSession contracts_oidc_session.IOIDCSession,
	wellknownCookies contracts_cookies.IWellknownCookies,
) (*service, error) {

	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container),
		oauth2Factory:    oauth2Factory,
		config:           config,
		oidcSession:      oidcSession,
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
		wellknown_echo.API_Start_ExternalLogin,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) getSession() (contracts_sessions.ISession, error) {
	session, err := s.oidcSession.GetSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *service) validateLoginGetRequest(model *external_idp.StartExternalIDPLoginRequest) error {

	if fluffycore_utils.IsEmptyOrNil(model.Slug) {
		return status.Error(codes.InvalidArgument, "Slug is required")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Directive) {
		return status.Error(codes.InvalidArgument, "Directive is required")
	}
	if !(model.Directive == "login" || model.Directive == "signup") {
		return status.Error(codes.InvalidArgument, "Directive must be 'login' or 'signup'")
	}
	return nil
}

type ErrorResponse struct {
	Error        string `json:"error"`
	InternalCode string `json:"internalCode"`
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	rootPath := echo_utils.GetMyRootPath(c)

	// is the request get or post?
	//rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &external_idp.StartExternalIDPLoginRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		var response = &ErrorResponse{
			Error:        err.Error(),
			InternalCode: InternalError_ExternalIDP_099,
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	if err := s.validateLoginGetRequest(model); err != nil {
		log.Error().Err(err).Msg("validateLoginGetRequest")
		var response = &ErrorResponse{
			Error:        err.Error(),
			InternalCode: InternalError_ExternalIDP_002,
		}
		return c.JSON(http.StatusBadRequest, response)

	}
	log.Debug().Interface("model", model).Msg("model")
	session, err := s.getSession()
	if err != nil {
		var response = &ErrorResponse{
			Error:        err.Error(),
			InternalCode: InternalError_ExternalIDP_003,
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	dd, err := session.Get("request")
	if err != nil {
		var response = &ErrorResponse{
			Error:        err.Error(),
			InternalCode: InternalError_ExternalIDP_004,
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	dd2 := dd.(*proto_oidc_models.AuthorizationRequest)
	getIDPBySlugResponse, err := s.IdpServiceServer().GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: model.Slug,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		var response = &ErrorResponse{
			Error:        err.Error(),
			InternalCode: InternalError_ExternalIDP_005,
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	idp := getIDPBySlugResponse.Idp
	externalState := xid.New().String()
	if idp.Protocol != nil {
		log.Debug().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch v := idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			{
				codeChallenge, verifier := generateCodeChallenge()
				externalOAuth2State := &proto_oidc_models.ExternalOauth2State{
					Request: &proto_oidc_models.ExternalOauth2Request{
						IdpHint:               model.Slug,
						ClientId:              v.Github.ClientId,
						State:                 dd2.State,
						CodeChallenge:         codeChallenge,
						CodeChallengeMethod:   "S256",
						CodeChallengeVerifier: verifier,
						Directive:             model.Directive,
						ParentState:           dd2.State,
					},
				}
				err = s.wellknownCookies.SetExternalOauth2Cookie(c, &contracts_cookies.SetExternalOauth2CookieRequest{
					State:               externalState,
					ExternalOAuth2State: externalOAuth2State,
				})

				if err != nil {
					log.Error().Err(err).Msg("SetExternalOauth2Cookie")
					// redirect to error page
					var response = &ErrorResponse{
						Error:        err.Error(),
						InternalCode: InternalError_ExternalIDP_006,
					}
					return c.JSON(http.StatusBadRequest, response)
				}
				getConfigResponse, err := s.oauth2Factory.GetConfig(ctx,
					&contracts_oauth2factory.GetConfigRequest{
						IDPHint: model.Slug,
					})
				if err != nil {
					log.Error().Err(err).Msg("Failed to get oauth2Config")
					var response = &ErrorResponse{
						Error:        err.Error(),
						InternalCode: InternalError_ExternalIDP_007,
					}
					return c.JSON(http.StatusBadRequest, response)
				}
				oauth2Config := getConfigResponse.Config
				u := oauth2Config.AuthCodeURL(externalState,
					oauth2.SetAuthURLParam("code_challenge", codeChallenge),
					oauth2.SetAuthURLParam("code_challenge_method", "S256"))
				return c.JSON(http.StatusOK, &external_idp.StartExternalIDPLoginResponse{
					RedirectURI: u,
				})

			}
		case *proto_oidc_models.Protocol_Oauth2:
			{
				codeChallenge, verifier := generateCodeChallenge()
				externalOAuth2State := &proto_oidc_models.ExternalOauth2State{
					Request: &proto_oidc_models.ExternalOauth2Request{
						IdpHint:               model.Slug,
						ClientId:              v.Oauth2.ClientId,
						State:                 dd2.State,
						CodeChallenge:         codeChallenge,
						CodeChallengeMethod:   "S256",
						CodeChallengeVerifier: verifier,
						Directive:             model.Directive,
						ParentState:           dd2.State,
					},
				}
				err = s.wellknownCookies.SetExternalOauth2Cookie(c, &contracts_cookies.SetExternalOauth2CookieRequest{
					State:               externalState,
					ExternalOAuth2State: externalOAuth2State,
				})

				if err != nil {
					log.Error().Err(err).Msg("SetExternalOauth2Cookie")
					// redirect to error page
					var response = &ErrorResponse{
						Error:        err.Error(),
						InternalCode: InternalError_ExternalIDP_008,
					}
					return c.JSON(http.StatusBadRequest, response)
				}
				scopes := strings.Split(v.Oauth2.Scope, " ")
				config := oauth2.Config{
					ClientID:     v.Oauth2.ClientId,
					ClientSecret: v.Oauth2.ClientSecret,
					Scopes:       scopes,
					RedirectURL:  rootPath + s.config.OIDCConfig.OAuth2CallbackPath,
					Endpoint: oauth2.Endpoint{
						AuthURL:  v.Oauth2.AuthorizationEndpoint,
						TokenURL: v.Oauth2.TokenEndpoint,
					},
				}

				u := config.AuthCodeURL(externalState,
					oauth2.SetAuthURLParam("code_challenge", codeChallenge),
					oauth2.SetAuthURLParam("code_challenge_method", "S256"))
				return c.JSON(http.StatusOK, &external_idp.StartExternalIDPLoginResponse{
					RedirectURI: u,
				})

			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				nonce := xid.New().String()
				externalOAuth2State := &proto_oidc_models.ExternalOauth2State{
					Request: &proto_oidc_models.ExternalOauth2Request{
						IdpHint:     model.Slug,
						ClientId:    v.Oidc.ClientId,
						State:       dd2.State,
						Directive:   model.Directive,
						ParentState: dd2.State,
						Nonce:       nonce,
					},
				}
				//codeChallenge, verifier := generateCodeChallenge()
				err = s.wellknownCookies.SetExternalOauth2Cookie(c, &contracts_cookies.SetExternalOauth2CookieRequest{
					State:               externalState,
					ExternalOAuth2State: externalOAuth2State,
				})

				if err != nil {
					log.Error().Err(err).Msg("SetExternalOauth2Cookie")
					// redirect to error page
					var response = &ErrorResponse{
						Error:        err.Error(),
						InternalCode: InternalError_ExternalIDP_009,
					}
					return c.JSON(http.StatusBadRequest, response)
				}
				getConfigResponse, err := s.oauth2Factory.GetConfig(ctx,
					&contracts_oauth2factory.GetConfigRequest{
						IDPHint: model.Slug,
					})
				if err != nil {
					log.Error().Err(err).Msg("Failed to get oauth2Config")
					var response = &ErrorResponse{
						Error:        err.Error(),
						InternalCode: InternalError_ExternalIDP_010,
					}
					return c.JSON(http.StatusBadRequest, response)
				}
				oauth2Config := getConfigResponse.Config
				authCodeOptions := []oauth2.AuthCodeOption{
					oauth2.SetAuthURLParam("nonce", nonce),
				}
				u := oauth2Config.AuthCodeURL(externalState, authCodeOptions...)
				return c.JSON(http.StatusOK, &external_idp.StartExternalIDPLoginResponse{
					RedirectURI: u,
				})
			}
		}
	}
	var response = &ErrorResponse{
		Error:        status.Error(codes.InvalidArgument, "Unsupported Protocol").Error(),
		InternalCode: InternalError_ExternalIDP_011,
	}
	return c.JSON(http.StatusBadRequest, response)

}

func generateCodeChallenge() (string, string) {
	// Generate a random 32-byte verifier string
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		panic(err) // Handle error appropriately in production code
	}
	verifier := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(verifierBytes)

	// Calculate the SHA256 hash of the verifier
	hash := sha256.Sum256([]byte(verifier))

	// Base64-encode the hash using URL-safe encoding without padding
	challenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	// Replace any '+' or '/' characters with '-' and '_', respectively
	challenge = strings.ReplaceAll(challenge, "+", "-")
	challenge = strings.ReplaceAll(challenge, "/", "_")

	return challenge, verifier
}

// ExternalIDP godoc
// @Summary starts an external login ceremony with an external IDP
// @Description starts an external login ceremony with an external IDP.
// @Tags 		root
// @Accept 		*/*
// @Produce 	json
// @Param		external_idp body		external_idp.StartExternalIDPLoginRequest	true	"StartExternalIDPLoginRequest"
// @Success 	200 				{object} 	external_idp.StartExternalIDPLoginResponse
// @Failure		400					{object}	ErrorResponse
// @Router /api/start-external-login [post]
func (s *service) Do(c echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		return s.DoPost(c)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
