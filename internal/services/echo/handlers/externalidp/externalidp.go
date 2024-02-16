package externalidp

import (
	"net/http"
	"strings"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	contracts_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/oauth2factory"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/utils"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
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
		oauth2Factory contracts_oauth2factory.IOAuth2Factory
		config        *contracts_config.Config
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	config *contracts_config.Config,
	container di.Container,
	oauth2Factory contracts_oauth2factory.IOAuth2Factory) (*service, error) {

	return &service{
		BaseHandler:   services_echo_handlers_base.NewBaseHandler(container),
		oauth2Factory: oauth2Factory,
		config:        config,
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
	Directive   string `param:"directive" query:"directive" form:"directive" json:"directive" xml:"directive"`
	State       string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
	IDPSlug     string `param:"idp_slug" query:"idp_slug" form:"idp_slug" json:"idp_slug" xml:"idp_slug"`
	RedirectURL string `param:"redirect_url" query:"redirect_url" form:"redirect_url" json:"redirect_url" xml:"redirect_url"`
}

func (s *service) validateLoginGetRequest(model *ExternalIDPAuthRequest) error {

	if fluffycore_utils.IsEmptyOrNil(model.State) {
		return status.Error(codes.InvalidArgument, "State is required")
	}
	if fluffycore_utils.IsEmptyOrNil(model.IDPSlug) {
		return status.Error(codes.InvalidArgument, "IDPSlug is required")
	}
	if fluffycore_utils.IsEmptyOrNil(model.Directive) {
		return status.Error(codes.InvalidArgument, "Directive is required")
	}
	if !(model.Directive == "login" || model.Directive == "signup") {
		return status.Error(codes.InvalidArgument, "Directive must be 'login' or 'signup'")
	}
	return nil
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
	if err := s.validateLoginGetRequest(model); err != nil {
		log.Error().Err(err).Msg("validateLoginGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("model", model).Msg("model")

	getIDPBySlugResponse, err := s.IdpServiceServer().GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: model.IDPSlug,
		})
	if err != nil {
		log.Error().Err(err).Msg("GetIDPBySlug")
		return c.Redirect(http.StatusFound, "/error")
	}
	idp := getIDPBySlugResponse.Idp
	externalState := xid.New().String()
	if idp.Protocol != nil {
		log.Info().Interface("getIDPBySlugResponse", getIDPBySlugResponse).Msg("getIDPBySlugResponse")
		switch v := idp.Protocol.Value.(type) {
		case *proto_oidc_models.Protocol_Github:
			{
				codeChallenge, verifier := generateCodeChallenge()
				err = s.ExternalOauth2FlowStore().StoreExternalOauth2Final(ctx,
					externalState,
					&models.ExternalOauth2Final{
						Request: &models.ExternalOauth2Request{
							IDPSlug:               model.IDPSlug,
							ClientID:              v.Github.ClientId,
							State:                 model.State,
							CodeChallenge:         codeChallenge,
							CodeChallengeMethod:   "S256",
							CodeChallengeVerifier: verifier,
							RedirectURL:           model.RedirectURL,
							Directive:             model.Directive,
							ParentState:           model.State,
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
				u := oauth2Config.AuthCodeURL(externalState,
					oauth2.SetAuthURLParam("code_challenge", codeChallenge),
					oauth2.SetAuthURLParam("code_challenge_method", "S256"))
				return c.Redirect(http.StatusFound, u)

			}
		case *proto_oidc_models.Protocol_Oauth2:
			{
				codeChallenge, verifier := generateCodeChallenge()

				err = s.ExternalOauth2FlowStore().StoreExternalOauth2Final(ctx,
					externalState,
					&models.ExternalOauth2Final{
						Request: &models.ExternalOauth2Request{
							IDPSlug:               model.IDPSlug,
							ClientID:              v.Oauth2.ClientId,
							State:                 model.State,
							CodeChallenge:         codeChallenge,
							CodeChallengeMethod:   "S256",
							CodeChallengeVerifier: verifier,
							RedirectURL:           model.RedirectURL,
							Directive:             model.Directive,
							ParentState:           model.State,
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

				u := config.AuthCodeURL(externalState,
					oauth2.SetAuthURLParam("code_challenge", codeChallenge),
					oauth2.SetAuthURLParam("code_challenge_method", "S256"))
				return c.Redirect(http.StatusFound, u)

			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				nonce := xid.New().String()

				//codeChallenge, verifier := generateCodeChallenge()
				err = s.ExternalOauth2FlowStore().StoreExternalOauth2Final(ctx,
					externalState,
					&models.ExternalOauth2Final{
						Request: &models.ExternalOauth2Request{
							IDPSlug:  model.IDPSlug,
							ClientID: v.Oidc.ClientId,
							State:    model.State,
							//			CodeChallenge:         codeChallenge,
							//			CodeChallengeMethod:   "S256",
							//			CodeChallengeVerifier: verifier,
							Nonce:       nonce,
							RedirectURL: model.RedirectURL,
							Directive:   model.Directive,
							ParentState: model.State,
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
				u := oauth2Config.AuthCodeURL(externalState, authCodeOptions...)
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
