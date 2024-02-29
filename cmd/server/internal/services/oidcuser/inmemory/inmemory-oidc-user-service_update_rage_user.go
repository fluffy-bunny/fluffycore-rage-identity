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

func (s *service) validateUpdateRageUserRequest(request *proto_oidc_user.UpdateRageUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User) {
		return status.Error(codes.InvalidArgument, "request.User is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity.Subject) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity.Subject is required")
	}
	return nil

}
func (s *service) UpdateRageUser(ctx context.Context, request *proto_oidc_user.UpdateRageUserRequest) (*proto_oidc_user.UpdateRageUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateUpdateRageUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateUpdateUserRequest")
		return nil, err
	}
	updateUserResponse, err := s.UpdateUser(ctx, &proto_external_user.UpdateUserRequest{
		User: &proto_external_models.ExampleUserUpdate{
			Id:       request.User.RootIdentity.Subject,
			RageUser: request.User,
		},
	})
	if err != nil {
		return nil, err
	}

	return &proto_oidc_user.UpdateRageUserResponse{
		User: updateUserResponse.User.RageUser,
	}, nil
}
