package externalidp

import (
	"net/http"
	"strings"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/oauth2factory"
	contracts_util "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/util"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	fluffycore_echo_contracts_contextaccessor "github.com/fluffy-bunny/fluffycore/echo/contracts/contextaccessor"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
	oauth2 "golang.org/x/oauth2"
)

type (
	service struct {
		services_echo_handlers_base.BaseHandler
		container               di.Container
		externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore
		idpServiceServer        proto_oidc_idp.IFluffyCoreIDPServiceServer
		someUtil                contracts_util.ISomeUtil
		oauth2Factory           contracts_oauth2factory.IOAuth2Factory
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(someUtil contracts_util.ISomeUtil,
	container di.Container,
	externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore,
	claimsPrincipal fluffycore_contracts_common.IClaimsPrincipal,
	idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer,
	oauth2Factory contracts_oauth2factory.IOAuth2Factory,
	echoContextAccessor fluffycore_echo_contracts_contextaccessor.IEchoContextAccessor) (*service, error) {

	return &service{
		BaseHandler: services_echo_handlers_base.BaseHandler{
			ClaimsPrincipal: claimsPrincipal, EchoContextAccessor: echoContextAccessor},
		container:               container,
		someUtil:                someUtil,
		externalOauth2FlowStore: externalOauth2FlowStore,
		idpServiceServer:        idpServiceServer,
		oauth2Factory:           oauth2Factory,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.POST,
		},
		wellknown_echo.ExternalIDPPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type ExternalIDPAuthRequest struct {
	IDPSlug string `param:"idp_slug" query:"idp_slug" form:"idp_slug" json:"idp_slug" xml:"idp_slug"`
}

func (s *service) DoPost(c echo.Context) error {
	r := c.Request()
	rootPath := echo_utils.GetMyRootPath(c)

	// is the request get or post?
	//rootPath := echo_utils.GetMyRootPath(c)
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &ExternalIDPAuthRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Info().Interface("model", model).Msg("model")

	getIDPBySlugResponse, err := s.idpServiceServer.GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: model.IDPSlug,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return c.Redirect(http.StatusFound, "/error")
	}
	idp := getIDPBySlugResponse.Idp
	if idp.Protocol != nil {
		log.Info().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch v := idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			{
				state := xid.New().String()
				codeChallenge, verifier := generateCodeChallenge()
				err = s.externalOauth2FlowStore.StoreExternalOauth2Final(ctx, state, &models.ExternalOauth2Final{
					Request: &models.ExternalOauth2Request{
						IDPSlug:               model.IDPSlug,
						ClientID:              v.Github.ClientId,
						State:                 state,
						CodeChallenge:         codeChallenge,
						CodeChallengeMethod:   "S256",
						CodeChallengeVerifier: verifier,
					},
				})
				if err != nil {
					log.Error().Err(err).Msg("StoreExternalOauth2Final")
					// redirect to error page
					return c.Redirect(http.StatusFound, "/error")
				}
				getConfigResponse, err := s.oauth2Factory.GetConfig(ctx,
					&contracts_oauth2factory.GetConfigRequest{
						IDPSlug: model.IDPSlug,
					})
				if err != nil {
					log.Error().Err(err).Msg("Failed to get oauth2Config")
					return c.Redirect(http.StatusFound, "/error")
				}
				oauth2Config := getConfigResponse.Config
				u := oauth2Config.AuthCodeURL(state,
					oauth2.SetAuthURLParam("code_challenge", codeChallenge),
					oauth2.SetAuthURLParam("code_challenge_method", "S256"))
				return c.Redirect(http.StatusFound, u)

			}
		case *proto_oidc_models.Protocol_Oauth2:
			{
				state := xid.New().String()
				codeChallenge, verifier := generateCodeChallenge()

				err = s.externalOauth2FlowStore.StoreExternalOauth2Final(ctx, state, &models.ExternalOauth2Final{
					Request: &models.ExternalOauth2Request{
						IDPSlug:               model.IDPSlug,
						ClientID:              v.Oauth2.ClientId,
						State:                 state,
						CodeChallenge:         codeChallenge,
						CodeChallengeMethod:   "S256",
						CodeChallengeVerifier: verifier,
					},
				})
				if err != nil {
					log.Error().Err(err).Msg("StoreExternalOauth2Final")
					// redirect to error page
					return c.Redirect(http.StatusFound, "/error")
				}
				scopes := strings.Split(v.Oauth2.Scope, " ")
				config := oauth2.Config{
					ClientID:     v.Oauth2.ClientId,
					ClientSecret: v.Oauth2.ClientSecret,
					Scopes:       scopes,
					RedirectURL:  rootPath + wellknown_echo.OAuth2CallbackPath,
					Endpoint: oauth2.Endpoint{
						AuthURL:  v.Oauth2.AuthorizationEndpoint,
						TokenURL: v.Oauth2.TokenEndpoint,
					},
				}

				u := config.AuthCodeURL(state,
					oauth2.SetAuthURLParam("code_challenge", codeChallenge),
					oauth2.SetAuthURLParam("code_challenge_method", "S256"))

				return c.Redirect(http.StatusFound, u)
			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				state := xid.New().String()
				nonce := xid.New().String()

				//codeChallenge, verifier := generateCodeChallenge()
				err = s.externalOauth2FlowStore.StoreExternalOauth2Final(ctx, state, &models.ExternalOauth2Final{
					Request: &models.ExternalOauth2Request{
						IDPSlug:  model.IDPSlug,
						ClientID: v.Oidc.ClientId,
						State:    state,
						//			CodeChallenge:         codeChallenge,
						//			CodeChallengeMethod:   "S256",
						//			CodeChallengeVerifier: verifier,
						Nonce: nonce,
					},
				})
				if err != nil {
					log.Error().Err(err).Msg("StoreExternalOauth2Final")
					// redirect to error page
					return c.Redirect(http.StatusFound, "/error")
				}
				getConfigResponse, err := s.oauth2Factory.GetConfig(ctx,
					&contracts_oauth2factory.GetConfigRequest{
						IDPSlug: model.IDPSlug,
					})
				if err != nil {
					log.Error().Err(err).Msg("Failed to get oauth2Config")
					return c.Redirect(http.StatusFound, "/error")
				}
				oauth2Config := getConfigResponse.Config
				authCodeOptions := []oauth2.AuthCodeOption{
					oauth2.SetAuthURLParam("nonce", nonce),
				}
				u := oauth2Config.AuthCodeURL(state, authCodeOptions...)
				return c.Redirect(http.StatusFound, u)
			}
		}
	}

	return c.Redirect(http.StatusFound, "/error")

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
