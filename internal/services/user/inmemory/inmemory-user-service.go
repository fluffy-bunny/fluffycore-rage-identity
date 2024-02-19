package inmemory

import (
	"context"
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		proto_oidc_user.UnimplementedUserServiceServer
		userMap map[string]*proto_oidc_models.User
		users   []*proto_oidc_models.User
		rwLock  sync.RWMutex
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_user.IFluffyCoreUserServiceServer = stemService
}
func (s *service) Ctor(clients *proto_oidc_models.Clients) (proto_oidc_user.IFluffyCoreUserServiceServer, error) {
	return &service{
		userMap: make(map[string]*proto_oidc_models.User),
	}, nil
}

func AddSingletonIFluffyCoreUserServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_user.IFluffyCoreUserServiceServer](cb, stemService.Ctor)
}

func (s *service) validateGetUserRequest(request *proto_oidc_user.GetUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Subject) {
		return status.Error(codes.InvalidArgument, "request.Subject is required")
	}
	return nil

}

func (s *service) GetUser(ctx context.Context, request *proto_oidc_user.GetUserRequest) (*proto_oidc_user.GetUserResponse, error) {
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
		return &proto_oidc_user.GetUserResponse{
			User: user,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "User not found")
}
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
	return &proto_oidc_user.UpdateUserResponse{
		User: getUserResp.User,
	}, nil
}
