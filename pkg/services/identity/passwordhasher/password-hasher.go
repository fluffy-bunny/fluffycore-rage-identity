package passwordhasher

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	passwordvalidator "github.com/wagslane/go-password-validator"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config *contracts_config.PasswordConfig
	}
)

var stemService = (*service)(nil)
var _ contracts_identity.IPasswordHasher = stemService

func (s *service) Ctor(config *contracts_config.PasswordConfig) (contracts_identity.IPasswordHasher, error) {
	// Compile the regex

	return &service{
		config: config,
	}, nil
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
func (s *service) validateIsAcceptablePasswordRequest(request *contracts_identity.IsAcceptablePasswordRequest) error {
	if fluffycore_utils.IsNil(request) {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Password) {
		return status.Error(codes.InvalidArgument, "Password is empty")
	}
	// a password that looks like an email is NOT allowed
	_, ok := echo_utils.IsValidEmailAddress(request.Password)
	if ok {
		// email not allowed
		return status.Error(codes.InvalidArgument, "the password may NOT be an email address")
	}
	return nil
}
func (s *service) IsAcceptablePassword(request *contracts_identity.IsAcceptablePasswordRequest) error {
	err := s.validateIsAcceptablePasswordRequest(request)
	if err != nil {
		return err
	}
	// regex expression check
	err = passwordvalidator.Validate(request.Password, s.config.MinEntropyBits)
	return err
}
