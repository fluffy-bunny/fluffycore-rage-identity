package callback

/*
reference: https://github.com/go-oauth2/oauth2/blob/master/example/client/client.go
*/
import (
	"net/http"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_codeexchange "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/codeexchange"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	contracts_util "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/util"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/client"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	jwxt "github.com/lestrrat-go/jwx/v2/jwt"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		someUtil                contracts_util.ISomeUtil
		scopedMemoryCache       fluffycore_contracts_common.IScopedMemoryCache
		externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore
		clientServiceServer     proto_oidc_client.IFluffyCoreClientServiceServer
		idpServiceServer        proto_oidc_idp.IFluffyCoreIDPServiceServer
		githubCodeExchange      contracts_codeexchange.IGithubCodeExchange
		genericOIDCCodeExchange contracts_codeexchange.IGenericOIDCCodeExchange
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache,
	clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer,
	externalOauth2FlowStore contracts_eko_gocache.IExternalOauth2FlowStore,
	idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer,
	githubCodeExchange contracts_codeexchange.IGithubCodeExchange,
	genericOIDCCodeExchange contracts_codeexchange.IGenericOIDCCodeExchange,
	someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		someUtil:                someUtil,
		scopedMemoryCache:       scopedMemoryCache,
		externalOauth2FlowStore: externalOauth2FlowStore,
		clientServiceServer:     clientServiceServer,
		idpServiceServer:        idpServiceServer,
		githubCodeExchange:      githubCodeExchange,
		genericOIDCCodeExchange: genericOIDCCodeExchange,
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
	log.Info().Interface("model", model).Msg("model")

	finalCache, err := s.externalOauth2FlowStore.GetExternalOauth2Final(ctx, model.State)
	if err != nil {
		log.Error().Err(err).Msg("GetExternalOauth2Final")
		return c.Redirect(http.StatusTemporaryRedirect, "/login?error=invalid_state")
	}
	// if we get here we are going to NUKE the cache for this transaction
	defer func() {
		s.externalOauth2FlowStore.DeleteExternalOauth2Final(ctx, model.State)
	}()
	getIDPBySlugResponse, err := s.idpServiceServer.GetIDPBySlug(ctx,
		&proto_oidc_idp.GetIDPBySlugRequest{
			Slug: finalCache.Request.IDPSlug,
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
					IDPSlug:      finalCache.Request.IDPSlug,
					ClientID:     finalCache.Request.ClientID,
					Nonce:        finalCache.Request.Nonce,
					Code:         model.Code,
					CodeVerifier: finalCache.Request.CodeChallenge,
				})
				if err != nil {
					log.Error().Err(err).Msg("ExchangeCode")
					return c.Redirect(http.StatusTemporaryRedirect, "/login?error=exchange_code")
				}
			}
		case *proto_oidc_models.Protocol_Oidc:
			{
				exchangeCodeResponse, err = s.genericOIDCCodeExchange.ExchangeCode(ctx, &contracts_codeexchange.ExchangeCodeRequest{
					IDPSlug:  finalCache.Request.IDPSlug,
					ClientID: finalCache.Request.ClientID,
					Nonce:    finalCache.Request.Nonce,
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

		// just save the id_token.  Its verifiable by the backend
		c.SetCookie(&http.Cookie{
			Name:     "_external_user",
			Value:    exchangeCodeResponse.IdToken,
			Path:     "/",
			Expires:  rawToken.Expiration(),
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

	}
	return c.Redirect(http.StatusTemporaryRedirect, "/login?code="+model.Code)
}
