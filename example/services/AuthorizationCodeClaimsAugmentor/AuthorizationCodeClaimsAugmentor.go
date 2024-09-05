package AuthorizationCodeClaimsAugmentor

/*
Default Claims Augmentor.
Does Nothing
*/
import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/tokenservice"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		fluffyCoreUserServiceServer proto_external_user.IFluffyCoreUserServiceServer
	}
)

var stemService = (*service)(nil)

var _ contracts_tokenservice.IAuthorizationCodeClaimsAugmentor = stemService

func (s *service) Ctor(fluffyCoreUserServiceServer proto_external_user.IFluffyCoreUserServiceServer) (contracts_tokenservice.IAuthorizationCodeClaimsAugmentor, error) {
	return &service{
		fluffyCoreUserServiceServer: fluffyCoreUserServiceServer,
	}, nil
}

func AddSingletonIClaimsAugmentor(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_tokenservice.IAuthorizationCodeClaimsAugmentor](cb, stemService.Ctor)
}

func (s *service) AugmentTokenClaims(ctx context.Context, request *contracts_tokenservice.AugmentTokenClaimsRequest) (*contracts_tokenservice.AugmentTokenClaimsResponse, error) {
	subjectI := request.IdTokenClaims.Get("sub")
	if subjectI == nil {
		return nil, status.Error(codes.InvalidArgument, "sub claim is missing")
	}

	getUserResponse, err := s.fluffyCoreUserServiceServer.GetUser(ctx, &proto_external_user.GetUserRequest{
		Subject: subjectI.(string),
	})
	if err != nil {
		return nil, err
	}
	request.IdTokenClaims.Set("integrity_id", "0")
	request.AccessTokenClaims.Set("integrity_id", "0")

	if fluffycore_utils.IsNotEmptyOrNil(getUserResponse.User.Metadata) {
		for _, metadata := range getUserResponse.User.Metadata {
			if metadata.Key == "integrity_id" {
				request.IdTokenClaims.Set("integrity_id", metadata.Value)
				request.AccessTokenClaims.Set("integrity_id", metadata.Value)
				break
			}
		}
	}

	return &contracts_tokenservice.AugmentTokenClaimsResponse{
		IdTokenClaims:     request.IdTokenClaims,
		AccessTokenClaims: request.AccessTokenClaims,
	}, nil
}
