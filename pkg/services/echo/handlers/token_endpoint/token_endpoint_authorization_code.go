package token_endpoint

import (
	"fmt"
	"net/http"
	"regexp"

	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/tokenservice"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	proto_events_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/events/types"
	proto_oidc_flows "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/flows"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_services_claims "github.com/fluffy-bunny/fluffycore/services/claims"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	oauth2 "github.com/go-oauth2/oauth2/v4"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type TokenEndpointAuthorizationCodeRequest struct {
	GrantType    string `param:"grant_type" query:"grant_type" form:"grant_type" json:"grant_type" xml:"grant_type"`
	Code         string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	CodeVerifier string `param:"code_verifier" query:"code_verifier" form:"code_verifier" json:"code_verifier" xml:"code_verifier"`
	RedirectURI  string `param:"redirect_uri" query:"redirect_uri" form:"redirect_uri" json:"redirect_uri" xml:"redirect_uri"`
}

type TokenEndpointAuthorizationCodeResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

func (s *service) validateTokenEndpointAuthorizationCodeRequest(req *TokenEndpointAuthorizationCodeRequest) error {
	if req.GrantType != string(oauth2.AuthorizationCode) {
		return status.Error(codes.InvalidArgument, "invalid_request - grant_type must be authorization_code")
	}
	if fluffycore_utils.IsEmptyOrNil(req.Code) {
		return status.Error(codes.InvalidArgument, "invalid_request - code is empty")
	}
	if req.RedirectURI == "" {
		return status.Error(codes.InvalidArgument, "invalid_request - redirect_uri is empty")
	}
	return nil
}
func extractIdpSlug(template string) (string, error) {
	// Define the regular expression pattern
	pattern := `^urn:rage:idp:([^:]+)?$`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	// Match the template against the regular expression
	match := re.FindStringSubmatch(template)
	if match == nil {
		return "", fmt.Errorf("invalid template format")
	}
	return match[1], nil
}
func sanitizeArray(input []string) []string {
	mm := make(map[string]bool)
	for _, v := range input {
		if fluffycore_utils.IsNotEmptyOrNil(v) {
			mm[v] = true
		}
	}

	var output []string
	for k := range mm {
		output = append(output, k)
	}
	return output
}
func (s *service) handleAuthorizationCode(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().
		Str("grant_type", string(oauth2.AuthorizationCode)).Logger()

	clientI, ok := s.scopedMemoryCache.Get("client")
	if !ok {
		log.Error().Msg("s.scopedMemoryCache.Get client")
		return c.String(http.StatusBadRequest, "s.scopedMemoryCache.Get client")
	}
	client, ok := clientI.(*proto_oidc_models.Client)
	if !ok {
		log.Error().Msg("clientI.(*proto_oidc_models.Client)")
		return c.String(http.StatusBadRequest, "clientI.(*proto_oidc_models.Client)")
	}

	allowedGrantType := false
	// is this client allowed to use the authorization_code grant type?
	for _, gt := range client.AllowedGrantTypes {
		allowedGrantType = (gt == string(oauth2.AuthorizationCode))
		if allowedGrantType {
			break
		}
	}
	if !allowedGrantType {
		log.Error().Msg("allowedGrantType")
		return c.String(http.StatusUnauthorized, "allowedGrantType")
	}

	req := &TokenEndpointAuthorizationCodeRequest{}
	if err := c.Bind(req); err != nil {
		log.Error().Err(err).Msg("Bind")
		return err
	}
	err := s.validateTokenEndpointAuthorizationCodeRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("validateTokenEndpointAuthorizationCodeRequest")
		return c.String(http.StatusBadRequest, err.Error())
	}
	log.Debug().Interface("req", req).Msg("req")

	getAuthorizationRequestStateResponse, err := s.authorizationRequestStateStore.GetAuthorizationRequestState(ctx,
		&proto_oidc_flows.GetAuthorizationRequestStateRequest{
			State: req.Code,
		})
	if err != nil {
		log.Warn().Err(err).Msg("GetAuthorizationRequestState")
		return c.String(http.StatusBadRequest, err.Error())
	}
	authorizationFinal := getAuthorizationRequestStateResponse.AuthorizationRequestState
	idClaims := fluffycore_services_claims.NewClaims()
	//--REQUIRED--
	idClaims.Set("aud", client.ClientId)
	idClaims.Set("nonce", authorizationFinal.Request.Nonce)
	idClaims.Set("sub", authorizationFinal.Identity.Subject)
	//--REQUIRED FOR US --
	idClaims.Set("client_id", client.ClientId)
	idClaims.Set("email", authorizationFinal.Identity.Email)
	idClaims.Set("email_verified", authorizationFinal.Identity.EmailVerified)
	acrClaims := []string{
		// always true
		models.ACRIdpRoot,
	}
	if len(authorizationFinal.Identity.Acr) > 0 {
		acrClaims = append(acrClaims, authorizationFinal.Identity.Acr...)
	}
	idpClaims := []string{}
	for _, acrValue := range acrClaims {
		idp, err := extractIdpSlug(acrValue)
		if err != nil {
			continue
		}
		idpClaims = append(idpClaims, idp)
	}
	amrClaims := authorizationFinal.Identity.Amr
	amrClaims = sanitizeArray(amrClaims)
	acrClaims = sanitizeArray(acrClaims)
	idpClaims = sanitizeArray(idpClaims)

	idClaims.Set("acr", acrClaims)
	idClaims.Set("amr", amrClaims)
	idClaims.Set("idp", idpClaims)
	//--OPTIONAL--

	// this one is opaque and is only really good for calling user_info endpoint
	accessTokenClaims := fluffycore_services_claims.NewClaims()
	accessTokenClaims.Set("client_id", client.ClientId)
	accessTokenClaims.Set("aud", client.ClientId)
	accessTokenClaims.Set("sub", authorizationFinal.Identity.Subject)

	augmentTokenClaimsResponse, err := s.claimsaugmentor.AugmentTokenClaims(ctx,
		&contracts_tokenservice.AugmentTokenClaimsRequest{
			IdTokenClaims:     idClaims,
			AccessTokenClaims: accessTokenClaims,
		})
	if err != nil {
		log.Error().Err(err).Msg("AugmentIdentityTokenClaims")
		return c.String(http.StatusBadRequest, err.Error())
	}

	idToken, err := s.tokenService.MintToken(ctx, &contracts_tokenservice.MintTokenRequest{
		Claims:                  augmentTokenClaimsResponse.IdTokenClaims,
		DurationLifeTimeSeconds: 3600,
		NotBeforeUnix:           0,
	})
	if err != nil {
		log.Warn().Err(err).Msg("MintToken - idToken")
		return c.String(http.StatusBadRequest, err.Error())
	}
	accessToken, err := s.tokenService.MintToken(ctx, &contracts_tokenservice.MintTokenRequest{
		Claims:                  augmentTokenClaimsResponse.AccessTokenClaims,
		DurationLifeTimeSeconds: 3600,
		NotBeforeUnix:           0,
	})
	if err != nil {
		log.Warn().Err(err).Msg("MintToken - accessToken")
		return c.String(http.StatusBadRequest, err.Error())
	}
	response := &TokenEndpointAuthorizationCodeResponse{
		AccessToken: accessToken.Token,
		TokenType:   "bearer",
		ExpiresIn:   3600,
		// TODO: Not something we want to support for an authentication only service.
		//RefreshToken: "refresh_token",
		IDToken: idToken.Token,
	}
	// return as json
	defer func() {
		log.Debug().Interface("response", response).Msg("response")
		_, err := s.authorizationRequestStateStore.DeleteAuthorizationRequestState(ctx, &proto_oidc_flows.DeleteAuthorizationRequestStateRequest{
			State: req.Code,
		})
		if err != nil {
			log.Error().Err(err).Msg("DeleteAuthorizationRequestState")
		}
	}()
	s.eventSink.OnEvent(ctx, &proto_events_types.Event{
		Event: &proto_events_types.Event_LoginEvent{
			LoginEvent: &proto_events_types.LoginEvent{
				Subject:        authorizationFinal.Identity.Subject,
				Email:          authorizationFinal.Identity.Email,
				ClientId:       client.ClientId,
				Acr:            acrClaims,
				Amr:            amrClaims,
				Idp:            idpClaims,
				LoginEventType: proto_events_types.LoginEventType_LOGIN_EVENT_TYPE_SUCCESS,
			},
		},
	})
	return c.JSONPretty(http.StatusOK, response, "  ")
}
