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
var _ contracts_tokenservice.IAuthorizationCodeClaimsAugmentor = stemService

func (s *service) Ctor() (contracts_tokenservice.IAuthorizationCodeClaimsAugmentor, error) {
	return &service{}, nil
}

func AddSingletonIClaimsAugmentor(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_tokenservice.IAuthorizationCodeClaimsAugmentor](cb, stemService.Ctor)
}

func (s *service) AugmentTokenClaims(ctx context.Context, request *contracts_tokenservice.AugmentTokenClaimsRequest) (*contracts_tokenservice.AugmentTokenClaimsResponse, error) {
	return &contracts_tokenservice.AugmentTokenClaimsResponse{
		IdTokenClaims:     request.IdTokenClaims,
		AccessTokenClaims: request.AccessTokenClaims,
	}, nil
}
