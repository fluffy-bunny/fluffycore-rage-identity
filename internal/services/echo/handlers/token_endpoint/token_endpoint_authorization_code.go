package token_endpoint

import (
	"net/http"

	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	oauth2 "github.com/go-oauth2/oauth2/v4"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type TokenEndpointAuthorizationCodeRequest struct {
	GrantType    string `param:"grant_type" query:"grant_type" form:"grant_type" json:"grant_type" xml:"grant_type"`
	Code         string `param:"code" query:"code" form:"code" json:"code" xml:"code"`
	CodeVerifier string `param:"code_verifier" query:"code_verifier" form:"code_verifier" json:"code_verifier" xml:"code_verifier"`
	RedirectURI  string `param:"redirect_uri" query:"redirect_uri" form:"redirect_uri" json:"redirect_uri" xml:"redirect_uri"`
}

type TokenEndpointAuthorizationCodeResponse struct {
	AccessToken  string `json:"access_token",omitempty`
	TokenType    string `json:"token_type",omitempty`
	ExpiresIn    int    `json:"expires_in",omitempty`
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

	response := &TokenEndpointAuthorizationCodeResponse{
		AccessToken:  "access_token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		RefreshToken: "refresh_token",
		IDToken:      "id_token",
	}
	// return as json
	return c.JSONPretty(http.StatusOK, response, "  ")
}
