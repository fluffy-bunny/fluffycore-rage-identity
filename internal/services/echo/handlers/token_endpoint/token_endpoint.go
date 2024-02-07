package token_endpoint

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
	clientauthorization "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/echo/middleware/clientauthorization"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/wellknown/echo"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		someUtil          contracts_util.ISomeUtil
		scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache
		oidcFlowCache     contracts_eko_gocache.IOIDCFlowCache
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(
	scopedMemoryCache fluffycore_contracts_common.IScopedMemoryCache,
	oidcFlowCache contracts_eko_gocache.IOIDCFlowCache,
	someUtil contracts_util.ISomeUtil) (*service, error) {
	return &service{
		someUtil:          someUtil,
		scopedMemoryCache: scopedMemoryCache,
		oidcFlowCache:     oidcFlowCache,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.OAuth2TokenEndpointPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		clientauthorization.AuthenticateOAuth2Client(),
	}
}

// HealthCheck godoc
// @Summary OAuth2 token endpoint.
// @Description OAuth2 token endpoint.
// @Tags root
// @Accept */*
// @Produce json
// @Security BasicAuth
// @Param       response_type    query     string  true  "response_type requested"
// @Param       scope            query     string  true  "scope requested" default("openid profile email")
// @Param       state            query     string  true  "state requested"
// @Param       redirect_uri     query     string  true  "redirect_uri requested"
// @Success 200 {object} string
// @Router /token [post]
func (s *service) Do(c echo.Context) error {

	//
	return c.Render(http.StatusOK, "views/home/index", map[string]interface{}{})
}

// This should be done on your server after receiving the authorization code
func (s *service) verifyCode(ctx context.Context, code string, codeVerifier string) bool {
	_, err := s.oidcFlowCache.Get(ctx, code)
	if err != nil {
		return false
	}

	// Get the code challenge back from the authorization server (assuming it stores it)
	codeChallengeFromDB := "..."

	// Decode the code verifier and stored code challenge
	verifierBytes, _ := base64.URLEncoding.DecodeString(codeVerifier)
	challengeBytes, _ := base64.URLEncoding.DecodeString(codeChallengeFromDB)

	// Calculate the hash of the verifier using the same method as the challenge
	verifierHash := sha256.Sum256(verifierBytes)

	// Compare the stored code challenge with the calculated hash
	if len(verifierHash) != len(challengeBytes) {
		return false
	}
	for i := 0; i < len(verifierHash); i++ {
		if verifierHash[i] != challengeBytes[i] {
			return false
		}
	}
	return true
}
