package token_endpoint

import (
	"net/http"

	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/tokenservice"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
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
func (s *service) handleAuthorizationCode(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().
		Str("grant_type", string(oauth2.AuthorizationCode)).Logger()

	clientI, err := s.scopedMemoryCache.Get("client")
	if err != nil {
		log.Error().Err(err).Msg("s.scopedMemoryCache.Get client")
		return c.String(http.StatusBadRequest, err.Error())
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
	err = s.validateTokenEndpointAuthorizationCodeRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("validateTokenEndpointAuthorizationCodeRequest")
		return c.String(http.StatusBadRequest, err.Error())
	}
	log.Debug().Interface("req", req).Msg("req")

	authFinal, err := s.oidcFlowStore.GetAuthorizationFinal(ctx, req.Code)
	if err != nil {
		log.Warn().Err(err).Msg("GetAuthorizationFinal")
		return c.String(http.StatusBadRequest, err.Error())
	}
	idClaims := fluffycore_services_claims.NewClaims()
	//--REQUIRED--
	idClaims.Set("aud", client.ClientId)
	idClaims.Set("nonce", authFinal.Request.Nonce)
	idClaims.Set("sub", authFinal.Identity.Subject)
	//--REQUIRED FOR US --
	idClaims.Set("client_id", client.ClientId)
	idClaims.Set("email", authFinal.Identity.Email)
	//--OPTIONAL--

	idToken, err := s.tokenService.MintToken(ctx, &contracts_tokenservice.MintTokenRequest{
		Claims:                  idClaims,
		DurationLifeTimeSeconds: 3600,
		NotBeforeUnix:           0,
	})
	if err != nil {
		log.Warn().Err(err).Msg("MintToken - idToken")
		return c.String(http.StatusBadRequest, err.Error())
	}
	// this one is opaque and is only really good for calling user_info endpoint
	accessTokenClaims := fluffycore_services_claims.NewClaims()
	accessTokenClaims.Set("client_id", client.ClientId)
	accessTokenClaims.Set("aud", client.ClientId)
	accessTokenClaims.Set("sub", authFinal.Identity.Subject)

	accessToken, err := s.tokenService.MintToken(ctx, &contracts_tokenservice.MintTokenRequest{
		Claims:                  accessTokenClaims,
		DurationLifeTimeSeconds: 3600,
		NotBeforeUnix:           0,
	})
	if err != nil {
		log.Warn().Err(err).Msg("MintToken - accessToken")
		return c.String(http.StatusBadRequest, err.Error())
	}
	response := &TokenEndpointAuthorizationCodeResponse{
		AccessToken:  accessToken.Token,
		TokenType:    "bearer",
		ExpiresIn:    3600,
		RefreshToken: "refresh_token",
		IDToken:      idToken.Token,
	}
	// return as json
	defer func() {
		log.Debug().Interface("response", response).Msg("response")
		err := s.oidcFlowStore.DeleteAuthorizationFinal(ctx, req.Code)
		if err != nil {
			log.Error().Err(err).Msg("DeleteAuthorizationFinal")
		}
	}()
	return c.JSONPretty(http.StatusOK, response, "  ")
}
