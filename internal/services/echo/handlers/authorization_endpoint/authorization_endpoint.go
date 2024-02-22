package authorization_endpoint

/*
reference: https://developers.onelogin.com/openid-connect/api/authorization-code
*/
import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/eko_gocache"
	models "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/models"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/wellknown/echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/client"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/idp"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler

		scopedMemoryCache   fluffycore_contracts_common.IScopedMemoryCache
		oidcFlowStore       contracts_eko_gocache.IOIDCFlowStore
		clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer
		idpServiceServer    proto_oidc_idp.IFluffyCoreIDPServiceServer
		userService         proto_oidc_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	container di.Container,
	idpServiceServer proto_oidc_idp.IFluffyCoreIDPServiceServer,
	userService proto_oidc_user.IFluffyCoreUserServiceServer,

	scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache,
	clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer,
	oidcFlowStore contracts_eko_gocache.IOIDCFlowStore) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),

		scopedMemoryCache:   scopedMemoryCache,
		oidcFlowStore:       oidcFlowStore,
		clientServiceServer: clientServiceServer,
		idpServiceServer:    idpServiceServer,
		userService:         userService,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OIDCAuthorizationEndpointPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		//clientauthorization.AuthenticateOAuth2Client(),
	}
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       client_id    			query     string  true  "client_id requested"
// @Param       response_type    		query     string  true  "response_type requested"
// @Param       scope            		query     string  true  "scope requested" default("openid profile email")
// @Param       state            		query     string  true  "state requested"
// @Param       redirect_uri     		query     string  true  "redirect_uri requested"
// @Param       audience     	 		query     string  false  "audience requested"
// @Param       code_challenge   		query     string  false  "PKCE challenge code"
// @Param       code_challenge_method 	query     string  false  "PKCE challenge method" default("S256")
// @Param       acr_values 				query     string  false  "acr_values requested"
// @Success 200 {object} string
// @Router /oidc/v1/auth [get]
func (s *service) Do(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &models.AuthorizationRequest{}
	if err := c.Bind(model); err != nil {
		return err
	}
	log.Debug().Interface("model", model).Msg("AuthorizationRequest")
	// TODO: validate the request
	// does the client have the permissions to do this?
	code := xid.New().String()
	model.Code = code
	// store the model in the cache.  Redis in production.
	authorizationFinal := &models.AuthorizationFinal{
		Request: model,
	}

	if !fluffycore_utils.IsEmptyOrNil(model.ACRValues) {
		acrValues := strings.Split(model.ACRValues, " ")

		idpHint := ""
		candidateUserID := ""
		for _, acrValue := range acrValues {
			d, err := extractIdpSlug(acrValue)
			if err == nil {
				v, ok := d["idp_hint"]
				if ok {
					idpHint = v
				}
			}
			d, err = extractRootCandidate(acrValue)
			if err == nil {
				v, ok := d["user_id"]
				if ok {
					candidateUserID = v
				}
			}
		}

		log.Info().Str("idpHint", idpHint).Str("rootCandidate", candidateUserID).Msg("acrValues")
		if !fluffycore_utils.IsEmptyOrNil(idpHint) {
			getIDPBySlugResponse, err := s.idpServiceServer.GetIDPBySlug(ctx, &proto_oidc_idp.GetIDPBySlugRequest{
				Slug: idpHint,
			})
			if err != nil ||
				getIDPBySlugResponse == nil ||
				getIDPBySlugResponse.Idp == nil {
				idpHint = ""
				c.Redirect(http.StatusFound, "/error&error=invalid_idp_hint")
			}
			model.IDPHint = idpHint
		}
		if !fluffycore_utils.IsEmptyOrNil(candidateUserID) {
			getUserResponse, err := s.userService.GetUser(ctx, &proto_oidc_user.GetUserRequest{
				Subject: candidateUserID,
			})
			if err != nil || getUserResponse == nil || getUserResponse.User == nil {
				candidateUserID = ""
				c.Redirect(http.StatusFound, "/error&error=invalid_root_candidate")
			}
			model.CandidateUserID = candidateUserID
		}

	}

	err := s.oidcFlowStore.StoreAuthorizationFinal(ctx, model.State, authorizationFinal)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}

	mm, err := s.oidcFlowStore.GetAuthorizationFinal(ctx, model.State)
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Info().Interface("mm", mm).Msg("mm")
	// redirect to the server Auth login pages.
	//s
	if fluffycore_utils.IsEmptyOrNil(model.IDPHint) {
		return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
			[]models.FormParam{
				{
					Name:  "state",
					Value: model.State,
				},
			})

	}
	return s.RenderAutoPost(c, wellknown_echo.ExternalIDPPath,
		[]models.FormParam{
			{
				Name:  "state",
				Value: model.State,
			},
			{
				Name:  "idp_hint",
				Value: model.IDPHint,
			},
			{
				Name:  "directive",
				Value: models.LoginDirective,
			},
		})

}

func extractIdpSlug(template string) (map[string]string, error) {
	// Define the regular expression pattern
	pattern := `^urn:mastodon:idp:([^:]+)?$`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Match the template against the regular expression
	match := re.FindStringSubmatch(template)
	if match == nil {
		return nil, fmt.Errorf("invalid template format")
	}

	// Extract and store the values
	info := make(map[string]string)
	info["idp_hint"] = match[1]

	return info, nil
}
func extractRootCandidate(template string) (map[string]string, error) {
	// Define the regular expression pattern
	pattern := `^urn:mastodon:root_candidate:([^:]+)?$`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Match the template against the regular expression
	match := re.FindStringSubmatch(template)
	if match == nil {
		return nil, fmt.Errorf("invalid template format")
	}

	// Extract and store the values
	info := make(map[string]string)
	info["user_id"] = match[1]

	return info, nil
}
