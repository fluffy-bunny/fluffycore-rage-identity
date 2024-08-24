package tokenservice

import (
	"context"

	fluffycore_contracts_claims "github.com/fluffy-bunny/fluffycore/contracts/claims"
)

type (
	MintTokenRequest struct {
		// Claims is a map of claims to be included in the token.
		// The standard claims of the token provider are added automatically.  i.e. issuer, etc
		Claims fluffycore_contracts_claims.IClaims
		// DurationLifeTimeSeconds is the duration of the token in seconds.
		// the final expiration time is calculated as NotBeforeUnix + DurationLifeTimeSeconds
		DurationLifeTimeSeconds int
		// NotBeforeUnix is the unix time in seconds that the token is not valid before.
		NotBeforeUnix int64
	}
	MintTokenResponse struct {
		Token      string
		Expiration int64
	}
	ITokenService interface {
		MintToken(ctx context.Context, request *MintTokenRequest) (*MintTokenResponse, error)
	}
	AugmentIdentityTokenClaimsRequest struct {
		Claims fluffycore_contracts_claims.IClaims
	}
	AugmentIdentityTokenClaimsResponse struct {
		Claims fluffycore_contracts_claims.IClaims
	}
	AugmentAccessTokenClaimsRequest struct {
		Claims fluffycore_contracts_claims.IClaims
	}
	AugmentAccessTokenClaimsResponse struct {
		Claims fluffycore_contracts_claims.IClaims
	}
	IAuthorizationCodeClaimsAugmentor interface {
		AugmentIdentityTokenClaims(ctx context.Context, request *AugmentIdentityTokenClaimsRequest) (*AugmentIdentityTokenClaimsResponse, error)
		AugmentAccessTokenClaims(ctx context.Context, request *AugmentAccessTokenClaimsRequest) (*AugmentAccessTokenClaimsResponse, error)
	}
)
