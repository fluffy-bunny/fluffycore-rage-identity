package tokenservice

import (
	"context"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/tokenservice"
	mocks_oauth2 "github.com/fluffy-bunny/fluffycore/mocks/oauth2"
	"github.com/rs/xid"
)

type (
	service struct {
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_tokenservice.ITokenService = stemService
}
func (s *service) Ctor() (contracts_tokenservice.ITokenService, error) {
	return &service{}, nil
}

func AddSingletonISomeUtil(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_tokenservice.ITokenService](cb, stemService.Ctor)
}

var (
	notAllowedClaims = map[string]bool{
		"iss": true,
		"jti": true,
		"iat": true,
		"nbf": true,
		"exp": true,
	}
)

func (s *service) sanitizeClaims(claims map[string]interface{}) map[string]interface{} {
	// remove any claims that are not strings
	newClaims := make(map[string]interface{})
	for k, v := range claims {
		if _, ok := notAllowedClaims[k]; !ok {
			newClaims[k] = v
		}
	}
	return newClaims
}
func (s *service) MintToken(ctx context.Context, request *contracts_tokenservice.MintTokenRequest) (*contracts_tokenservice.MintTokenResponse, error) {
	claims := mocks_oauth2.NewClaims()
	sanitizedClaims := s.sanitizeClaims(request.Claims)
	for k, v := range sanitizedClaims {
		claims.Set(k, v)
	}

	now := time.Now()
	claims.Set("iat", now.Unix())

	// 5 minute clock skew
	nbf := now.Add(-5 * time.Minute).Unix()
	if request.NotBeforeUnix > 0 {
		nbf = int64(request.NotBeforeUnix)
	}
	claims.Set("nbf", nbf)
	exp := now.Add(time.Duration(request.DurationLifeTimeSeconds) * time.Second).Unix()

	nbfTime := time.Unix(nbf, 0)
	if nbfTime.After(now) {
		// make a new exp
		exp = nbfTime.Add(time.Duration(request.DurationLifeTimeSeconds) * time.Second).Unix()
	}
	claims.Set("exp", exp)
	claims.Set("iss", "http://localhost:9044")
	claims.Set("jti", xid.New().String())
	token, _ := mocks_oauth2.MintToken(claims)

	return &contracts_tokenservice.MintTokenResponse{
		Token:      token,
		Expiration: exp,
	}, nil
}
