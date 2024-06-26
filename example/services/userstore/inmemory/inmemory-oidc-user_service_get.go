package inmemory

import (
	"context"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

func (s *service) validateGetUserRequest(request *proto_external_user.GetUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Subject) {
		return status.Error(codes.InvalidArgument, "request.Subject is required")
	}
	return nil

}

func (s *service) makeExampleUserCopy(user *proto_oidc_models.ExampleUser) *proto_oidc_models.ExampleUser {
	if user == nil {
		return nil
	}
	d, err := protojson.Marshal(user)
	if err != nil {
		return nil
	}
	var newUser proto_oidc_models.ExampleUser
	err = protojson.Unmarshal(d, &newUser)
	if err != nil {
		return nil
	}
	return &newUser
}
func (s *service) GetUser(ctx context.Context, request *proto_external_user.GetUserRequest) (*proto_external_user.GetUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateGetUserRequest")
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	user, ok := s.userMap[request.Subject]
	if ok {
		return &proto_external_user.GetUserResponse{
			User: s.makeExampleUserCopy(user),
		}, nil
	}
	// TODO: THIS MUST BE ENCRYPTED AT REST
	// if there is a TOTP secret, we need to return it
	// it better be encrypted at rest
	return nil, status.Error(codes.NotFound, "User not found")
}
