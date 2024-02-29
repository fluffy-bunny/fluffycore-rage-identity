package inmemory

import (
	"context"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

func (s *service) validateGetRageUserRequest(request *proto_oidc_user.GetRageUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Subject) {
		return status.Error(codes.InvalidArgument, "request.Subject is required")
	}
	return nil

}

func (s *service) makeRageUserCopy(user *proto_oidc_models.RageUser) *proto_oidc_models.RageUser {
	if user == nil {
		return nil
	}
	d, err := protojson.Marshal(user)
	if err != nil {
		return nil
	}
	var newUser proto_oidc_models.RageUser
	err = protojson.Unmarshal(d, &newUser)
	if err != nil {
		return nil
	}
	return &newUser
}
func (s *service) GetRageUser(ctx context.Context, request *proto_oidc_user.GetRageUserRequest) (*proto_oidc_user.GetRageUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetRageUserRequest(request)
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
		return &proto_oidc_user.GetRageUserResponse{
			User: s.makeRageUserCopy(user.RageUser),
		}, nil
	}
	return nil, status.Error(codes.NotFound, "User not found")
}
