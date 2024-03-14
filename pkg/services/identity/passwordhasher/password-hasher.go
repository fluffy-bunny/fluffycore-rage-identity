package passwordhasher

import (
	"context"
	"regexp"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	passwordvalidator "github.com/wagslane/go-password-validator"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		config          *contracts_config.PasswordConfig
		regexExpression *regexp.Regexp
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_identity.IPasswordHasher = stemService
}
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
func (s *service) IsAcceptablePassword(user *proto_oidc_models.RageUser, password string) error {
	if user == nil {
		return status.Error(codes.InvalidArgument, "user is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(password) {
		return status.Error(codes.InvalidArgument, "Password is empty")
	}
	// stupidity check
	emailAndPasswordTheSame := strings.EqualFold(user.RootIdentity.Email, password)
	if emailAndPasswordTheSame {
		return status.Error(codes.InvalidArgument, "Password cannot be the same as the email")
	}
	// regex expression check
	err := passwordvalidator.Validate(password, s.config.MinEntropyBits)
	return err

}
