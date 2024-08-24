package AuthorizationCodeClaimsAugmentor

/*
Default Claims Augmentor.
Does Nothing
*/
import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/tokenservice"
)

type (
	service struct {
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_tokenservice.IAuthorizationCodeClaimsAugmentor = stemService
}
func (s *service) Ctor() (contracts_tokenservice.IAuthorizationCodeClaimsAugmentor, error) {
	return &service{}, nil
}

func AddSingletonIClaimsAugmentor(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_tokenservice.IAuthorizationCodeClaimsAugmentor](cb, stemService.Ctor)
}

func (s *service) AugmentIdentityTokenClaims(ctx context.Context, request *contracts_tokenservice.AugmentIdentityTokenClaimsRequest) (*contracts_tokenservice.AugmentIdentityTokenClaimsResponse, error) {
	request.Claims.Set("id_test", "hello")
	return &contracts_tokenservice.AugmentIdentityTokenClaimsResponse{
		Claims: request.Claims,
	}, nil
}

func (s *service) AugmentAccessTokenClaims(ctx context.Context, request *contracts_tokenservice.AugmentAccessTokenClaimsRequest) (*contracts_tokenservice.AugmentAccessTokenClaimsResponse, error) {
	request.Claims.Set("access_test", "hello")
	return &contracts_tokenservice.AugmentAccessTokenClaimsResponse{
		Claims: request.Claims,
	}, nil
}
