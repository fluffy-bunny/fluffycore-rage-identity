package tokenservice

import (
	"context"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/internal/contracts/config"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/internal/contracts/tokenservice"
	fluffycore_contracts_claims "github.com/fluffy-bunny/fluffycore/contracts/claims"
	fluffycore_contracts_jwtminter "github.com/fluffy-bunny/fluffycore/contracts/jwtminter"
	fluffycore_services_claims "github.com/fluffy-bunny/fluffycore/services/claims"
	xid "github.com/rs/xid"
)

type (
	service struct {
		jwtMinter  fluffycore_contracts_jwtminter.IJWTMinter
		oidcConfig *contracts_config.OIDCConfig
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_tokenservice.ITokenService = stemService
}
func (s *service) Ctor(
	oidcConfig *contracts_config.OIDCConfig,
	jwtMinter fluffycore_contracts_jwtminter.IJWTMinter,
) (contracts_tokenservice.ITokenService, error) {
	return &service{
		jwtMinter:  jwtMinter,
		oidcConfig: oidcConfig,
	}, nil
}

func AddSingletonITokenService(cb di.ContainerBuilder) {
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

func dedupStringArray(arr []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for v := range arr {
		if !encountered[arr[v]] {
			encountered[arr[v]] = true
			result = append(result, arr[v])
		}
	}
	return result
}
func (s *service) sanitizeClaims(claims fluffycore_contracts_claims.IClaims) fluffycore_contracts_claims.IClaims {
	// remove any claims that are not strings
	newClaims := fluffycore_services_claims.NewClaims()
	for _, v := range claims.Claims() {
		claimType := v.Type()
		if _, ok := notAllowedClaims[claimType]; !ok {
			// is string array
			value := v.Value()
			va, ok := value.([]string)
			if ok {
				value = dedupStringArray(va)
			}
			newClaims.Set(claimType, value)
		}
	}
	return newClaims
}
func (s *service) MintToken(ctx context.Context, request *contracts_tokenservice.MintTokenRequest) (*contracts_tokenservice.MintTokenResponse, error) {

	claims := s.sanitizeClaims(request.Claims)

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
	claims.Set("iss", s.oidcConfig.BaseUrl)
	claims.Set("jti", xid.New().String())
	token, err := s.jwtMinter.MintToken(ctx, claims)
	if err != nil {
		return nil, err
	}

	return &contracts_tokenservice.MintTokenResponse{
		Token:      token,
		Expiration: exp,
	}, nil
}
