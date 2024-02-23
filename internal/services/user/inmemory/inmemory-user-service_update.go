package inmemory

import (
	"context"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateUpdateUserRequest(request *proto_oidc_user.UpdateUserRequest) error {
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
func (s *service) UpdateUser(ctx context.Context, request *proto_oidc_user.UpdateUserRequest) (*proto_oidc_user.UpdateUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateUpdateUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateUpdateUserRequest")
		return nil, err
	}
	getUserResp, err := s.GetUser(ctx, &proto_oidc_user.GetUserRequest{
		Subject: request.User.RootIdentity.Subject,
	})
	if err != nil {
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	if request.User.RootIdentity.EmailVerified != nil {
		getUserResp.User.RootIdentity.EmailVerified = request.User.RootIdentity.EmailVerified.Value
	}
	if request.User.Password != nil {
		if request.User.Password.Hash != nil {

			getUserResp.User.Password = &proto_oidc_models.Password{
				Hash: request.User.Password.Hash.Value,
			}
		}
	}
	return &proto_oidc_user.UpdateUserResponse{
		User: getUserResp.User,
	}, nil
}
