package tokenservice

import "context"

type (
	MintTokenRequest struct {
		// Claims is a map of claims to be included in the token.
		// The standard claims of the token provider are added automatically.  i.e. issuer, etc
		Claims map[string]interface{}
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
)
