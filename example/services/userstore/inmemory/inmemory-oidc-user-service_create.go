package inmemory

import (
	"context"

	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
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

	if fluffycore_utils.IsEmptyOrNil(request.User.Id) {
		return status.Error(codes.InvalidArgument, "request.User.Id is required")
	}

	if request.User.RageUser == nil {
		return status.Error(codes.InvalidArgument, "request.User.RageUser is required")
	}
	rageUser := request.User.RageUser
	if rageUser.RootIdentity == nil {
		return status.Error(codes.InvalidArgument, "request.User.RageUser.RootIdentity is required")
	}
	rageUser.RootIdentity.Subject = request.User.Id

	if fluffycore_utils.IsEmptyOrNil(rageUser.RootIdentity.Subject) {
		return status.Error(codes.InvalidArgument, "request.User.RageUser.RootIdentity.Subject is required")
	}
	if fluffycore_utils.IsEmptyOrNil(rageUser.RootIdentity.Email) {
		return status.Error(codes.InvalidArgument, "request.User.RageUser.RootIdentity.Email is required")
	}
	if fluffycore_utils.IsEmptyOrNil(rageUser.RootIdentity.IdpSlug) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity.IdpSlug is required")
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

	user := request.User
	getUserResponse, err := s.GetUser(ctx, &proto_external_user.GetUserRequest{
		Subject: user.Id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.NotFound {
			return nil, err
		}
	} else {
		return &proto_external_user.CreateUserResponse{
			User: getUserResponse.User,
		}, nil
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//

	// create the user
	s.userMap[user.Id] = user
	return &proto_external_user.CreateUserResponse{
		User: s.makeExampleUserCopy(user),
	}, nil
}
