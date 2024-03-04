package inmemory

import (
	"context"

	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"

	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateCreateRageUserRequest(request *proto_oidc_user.CreateRageUserRequest) error {
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

// CreateRageUser creates a new user
// We pass it to the external service to create the user.  The RageUser is a sub-object of the user.
func (s *service) CreateRageUser(ctx context.Context, request *proto_oidc_user.CreateRageUserRequest) (*proto_oidc_user.CreateRageUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateCreateRageUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateCreateUserRequest")
		return nil, err
	}

	user := request.User
	getUserResponse, err := s.GetRageUser(ctx, &proto_oidc_user.GetRageUserRequest{
		By: &proto_oidc_user.GetRageUserRequest_Subject{
			Subject: user.RootIdentity.Subject,
		},
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != codes.NotFound {
			return nil, err
		}
	} else {
		return &proto_oidc_user.CreateRageUserResponse{
			User: getUserResponse.User,
		}, nil
	}

	createUserResponse, err := s.CreateUser(ctx, &proto_external_user.CreateUserRequest{
		User: &proto_external_models.ExampleUser{
			Id:       user.RootIdentity.Subject,
			RageUser: user,
		},
	})
	if err != nil {
		return nil, err
	}
	// create the user
	return &proto_oidc_user.CreateRageUserResponse{
		User: s.makeRageUserCopy(createUserResponse.User.RageUser),
	}, nil
}
