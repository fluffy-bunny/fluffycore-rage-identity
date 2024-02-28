package inmemory

import (
	"context"

	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	status "github.com/gogo/status"
	xid "github.com/rs/xid"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateCreateUserRequest(request *proto_external_user.CreateUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if request.User == nil {
		return status.Error(codes.InvalidArgument, "request.User is required")
	}
	if request.User.Metadata == nil {
		request.User.Metadata = make(map[string]string)
	}
	return nil
}
func (s *service) CreateUser(ctx context.Context, request *proto_external_user.CreateUserRequest) (*proto_external_user.CreateUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateCreateUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateCreateUserRequest")
		return nil, err
	}

	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	id := xid.New().String()
	request.User.Id = id
	// create the user
	s.userMap[id] = request.User

	return &proto_external_user.CreateUserResponse{
		User: s.makeUserCopy(request.User),
	}, nil
}
