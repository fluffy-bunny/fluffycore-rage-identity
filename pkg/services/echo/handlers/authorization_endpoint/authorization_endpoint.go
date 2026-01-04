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
	contracts_cache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cache"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	models_api_manifest "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/manifest"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		scopedMemoryCache                    contracts_cache.IScopedMemoryCache
		authorizationRequestStateStoreServer proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer
		clientServiceServer                  proto_oidc_client.IFluffyCoreClientServiceServer
		idpServiceServer                     proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer
		userService                          proto_oidc_user.IFluffyCoreRageUserServiceServer
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	idpServiceServer proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer,
	userService proto_oidc_user.IFluffyCoreRageUserServiceServer,
	scopedMemoryCache contracts_cache.IScopedMemoryCache,
	clientServiceServer proto_oidc_client.IFluffyCoreClientServiceServer,
	authorizationRequestStateStoreServer proto_oidc_flows.IFluffyCoreAuthorizationRequestStateStoreServer,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),

		scopedMemoryCache:                    scopedMemoryCache,
		authorizationRequestStateStoreServer: authorizationRequestStateStoreServer,
		clientServiceServer:                  clientServiceServer,
		idpServiceServer:                     idpServiceServer,
		userService:                          userService,
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
func (s *service) newSession() (contracts_sessions.ISession, error) {
	session, err := s.SessionFactory().
		GetCookieSession(s.WellknownCookieNames().GetCookieName(contracts_cookies.CookieName_OIDCSession))
	if err != nil {
		return nil, err
	}

	err = session.New()
	if err != nil {
		return nil, err
	}
	return session, nil
}

type (
	AuthorizationRequest struct {
		ClientId            string `param:"client_id" query:"client_id" form:"client_id" json:"client_id" xml:"client_id"`
		ResponseType        string `param:"response_type" query:"response_type" form:"response_type" json:"response_type" xml:"response_type"`
		Scope               string `param:"scope" query:"scope" form:"scope" json:"scope" xml:"scope"`
		State               string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
		RedirectURI         string `param:"redirect_uri" query:"redirect_uri" form:"redirect_uri" json:"redirect_uri" xml:"redirect_uri"`
		Audience            string `param:"audience" query:"audience" form:"audience" json:"audience" xml:"audience"`
		CodeChallenge       string `param:"code_challenge" query:"code_challenge" form:"code_challenge" json:"code_challenge" xml:"code_challenge"`
		CodeChallengeMethod string `param:"code_challenge_method" query:"code_challenge_method" form:"code_challenge_method" json:"code_challenge_method" xml:"code_challenge_method"`
		ACRValues           string `param:"acr_values" query:"acr_values" form:"acr_values" json:"acr_values" xml:"acr_values"`
		Nonce               string `param:"nonce" query:"nonce" form:"nonce" json:"nonce" xml:"nonce"`
		Code                string // this is the internal code that will be returned to the OIDC client
		// IDPHint is the idp_hint of the external idp that the authorization must authenticate against
		IDPHint string `json:"idp_hint,omitempty"`
		// CandidateUserID is the user_id of the candidate user that if the external IDP has no link should be linked to
		// The candidate user must exist.
		CandidateUserID string `json:"candidate_user_id,omitempty"`
	}
)

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
	echoModel := &AuthorizationRequest{}
	if err := c.Bind(echoModel); err != nil {
		return err
	}
	log.Debug().Interface("echoModel", echoModel).Msg("AuthorizationRequest")
	session, err := s.newSession()
	if err != nil {
		return err
	}
	// clean out left overs.
	s.WellknownCookies().DeleteAuthCompletedCookie(c)
	s.WellknownCookies().DeleteAuthCookie(c)

	model := &proto_oidc_models.AuthorizationRequest{}
	fluffycore_utils.ConvertStructToProto(echoModel, model)
	// TODO: validate the request
	// does the client have the permissions to do this?
	code := xid.New().String()
	model.Code = code

	// store the model in the cache.  Redis in production.
	authorizationFinal := &proto_oidc_models.AuthorizationRequestState{
		Request: model,
	}

	if fluffycore_utils.IsNotEmptyOrNil(model.AcrValues) {
		acrValues := strings.Split(model.AcrValues, " ")

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

		log.Debug().Str("idpHint", idpHint).Str("rootCandidate", candidateUserID).Msg("acrValues")
		if fluffycore_utils.IsNotEmptyOrNil(idpHint) {
			listIDPResponse, err := s.IdpServiceServer().ListIDP(ctx,
				&proto_oidc_idp.ListIDPRequest{
					Filter: &proto_oidc_idp.Filter{
						Enabled: &proto_types.BoolFilterExpression{
							Eq: true,
						},
						Slug: &proto_types.StringFilterExpression{
							Eq: idpHint,
						},
					},
				})
			if err != nil {
				log.Error().Err(err).Msg("ListIDP")
				return c.Redirect(http.StatusFound, "/error&error=invalid_idp_hint")
			}
			if listIDPResponse == nil || len(listIDPResponse.IDPs) == 0 {
				return c.Redirect(http.StatusFound, "/error&error=invalid_idp_hint")
			}
			model.IdpHint = idpHint

			model.IdpHint = idpHint
		}
		if fluffycore_utils.IsNotEmptyOrNil(candidateUserID) {
			getUserResponse, err := s.userService.GetRageUser(ctx,
				&proto_oidc_user.GetRageUserRequest{
					By: &proto_oidc_user.GetRageUserRequest_Subject{
						Subject: candidateUserID,
					},
				})
			if err != nil || getUserResponse == nil || getUserResponse.User == nil {
				candidateUserID = ""
				c.Redirect(http.StatusFound, "/error&error=invalid_root_candidate")
			}
			model.CandidateUserId = candidateUserID
		}

	}
	type authorizationStateContainer struct {
		State string `json:"state"`
	}
	// we let the client know the state so that it can use it to store semi static data against it.
	err = s.WellknownCookies().SetInsecureCookie(c,
		s.WellknownCookieNames().GetCookieName(contracts_cookies.CookieName_AuthorizationState),
		&authorizationStateContainer{
			State: model.State,
		})
	if err != nil {
		return err
	}

	_, err = s.authorizationRequestStateStoreServer.StoreAuthorizationRequestState(ctx,
		&proto_oidc_flows.StoreAuthorizationRequestStateRequest{
			AuthorizationRequestState: authorizationFinal,
			State:                     model.State,
		})
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	// set the code and state in the session
	// --~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//

	// Panic protection for session.Set operations
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Msg("panic recovered during session.Set, clearing session")
				// Clear the session cookie

				err = fmt.Errorf("session storage error: unable to store request data")
			}
		}()

		err = session.Set("request", model)
		if err != nil {
			log.Error().Err(err).Msg("failed to set session request")
			// Clear the session on error

			return
		}
		sessionId := model.Nonce
		if fluffycore_utils.IsEmptyOrNil(sessionId) {
			sessionId = xid.New().String()
		}
		session.Set("session_id", sessionId)
		session.Set("landing_page", &models_api_manifest.LandingPage{
			Page: models_api_manifest.PageUsernameEntry,
		})
		session.Save()
	}()

	if err != nil {
		// Return a friendly JSON error response
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":             "session_storage_error",
			"error_description": "Unable to process your request due to a session storage issue. This may be caused by session data being too large or configuration issues.",
			"message":           "Please try again or contact support if the problem persists.",
		})
	}

	mm, err := s.authorizationRequestStateStoreServer.GetAuthorizationRequestState(ctx, &proto_oidc_flows.GetAuthorizationRequestStateRequest{
		State: model.State,
	})
	if err != nil {
		// redirect to error page
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Debug().Interface("mm", mm).Msg("mm")
	// redirect to the server Auth login pages.
	//s
	if fluffycore_utils.IsEmptyOrNil(model.IdpHint) {
		return s.RenderAutoPost(c, wellknown_echo.OIDCLoginPath,
			[]models.FormParam{
				{
					Name:  "state",
					Value: model.State,
				},
			})

	}
	// clear out any error cookies
	s.WellknownCookies().DeleteErrorCookie(c)

	return s.RenderAutoPost(c, wellknown_echo.ExternalIDPPath,
		[]models.FormParam{
			{
				Name:  "state",
				Value: model.State,
			},
			{
				Name:  "idp_hint",
				Value: model.IdpHint,
			},
			{
				Name:  "directive",
				Value: models.LoginDirective,
			},
		})

}

func extractIdpSlug(template string) (map[string]string, error) {
	// Define the regular expression pattern
	pattern := `^urn:rage:idp:([^:]+)?$`

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
	pattern := `^urn:rage:root_candidate:([^:]+)?$`

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
