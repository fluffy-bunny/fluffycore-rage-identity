package inmemory

import (
	"context"

	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	"github.com/gogo/status"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
)

func (s *service) validateCreateUserRequest(request *proto_oidc_user.CreateUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if request.User == nil {
		return status.Error(codes.InvalidArgument, "request.User is required")
	}
	if request.User.RootIdentity == nil {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity.Subject) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity.Subject is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity.Email) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity.Email is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity.IdpSlug) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity.IdpSlug is required")
	}
	return nil
}
func (s *service) CreateUser(ctx context.Context, request *proto_oidc_user.CreateUserRequest) (*proto_oidc_user.CreateUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateCreateUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateCreateUserRequest")
		return nil, err
	}

	user := request.User
	getUserResponse, err := s.GetUser(ctx, &proto_oidc_user.GetUserRequest{
		Subject: user.RootIdentity.Subject,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.NotFound {
			return nil, err
		}
	} else {
		return &proto_oidc_user.CreateUserResponse{
			User: getUserResponse.User,
		}, nil
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//

	// create the user
	s.userMap[user.RootIdentity.Subject] = user
	return &proto_oidc_user.CreateUserResponse{
		User: s.makeUserCopy(user),
	}, nil
}
