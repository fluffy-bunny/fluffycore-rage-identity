package passwordhasher

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct{}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_identity.IPasswordHasher = stemService
}
func (s *service) Ctor() (contracts_identity.IPasswordHasher, error) {
	return &service{}, nil
}

func AddSingletonIPasswordHasher(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_identity.IPasswordHasher](cb, stemService.Ctor)
}
func (s *service) validateHashPasswordRequest(request *contracts_identity.HashPasswordRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")

	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		return status.Error(codes.InvalidArgument, "Password is empty")

	}
	return nil
}

func (s *service) HashPassword(ctx context.Context, request *contracts_identity.HashPasswordRequest) (*contracts_identity.HashPasswordResponse, error) {
	if err := s.validateHashPasswordRequest(request); err != nil {
		return nil, err
	}
	hash, err := utils.GeneratePasswordHash(request.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contracts_identity.HashPasswordResponse{
		Password:       request.Password,
		HashedPassword: hash,
	}, nil
}
func (s *service) validateVerifyPasswordRequest(request *contracts_identity.VerifyPasswordRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		return status.Error(codes.InvalidArgument, "Password is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.HashedPassword) {
		return status.Error(codes.InvalidArgument, "HashedPassword is empty")
	}
	return nil
}
func (s *service) VerifyPassword(ctx context.Context, request *contracts_identity.VerifyPasswordRequest) error {
	if err := s.validateVerifyPasswordRequest(request); err != nil {
		return err
	}
	ok, err := utils.ComparePasswordHash(request.Password, request.HashedPassword)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !ok {
		return status.Error(codes.NotFound, "Password does not match")
	}
	return nil
}
