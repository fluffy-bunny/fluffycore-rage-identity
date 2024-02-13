package tokenservice

import (
	"context"
	"testing"

	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/tokenservice"
	require "github.com/stretchr/testify/require"
)

func TestMintToken(t *testing.T) {
	tokenService, err := stemService.Ctor()
	require.NoError(t, err)
	require.NotNil(t, tokenService)

	ctx := context.Background()

	response, err := tokenService.MintToken(ctx, &contracts_tokenservice.MintTokenRequest{
		Claims: map[string]interface{}{
			"sub":   "1234567890",
			"name":  "John Doe",
			"email": "john.doe@gmail.com",
			"jti":   "ShouldBeRemoved",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotEmpty(t, response.Token)
}
