package inmemory

import (
	"context"

	linq "github.com/ahmetb/go-linq/v3"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateUnlinkUsersRequest(request *proto_oidc_user.UnlinkUsersRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.RootSubject) {
		return status.Error(codes.InvalidArgument, "request.Subject is required")
	}
	if request.ExternalIdentity == nil {
		return status.Error(codes.InvalidArgument, "request.ExternalIdentity is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ExternalIdentity.Subject) {
		return status.Error(codes.InvalidArgument, "request.ExternalIdentity.Subject is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ExternalIdentity.IdpSlug) {
		return status.Error(codes.InvalidArgument, "request.ExternalIdentity.IdpSlug is required")
	}
	return nil
}
func (s *service) UnlinkUsers(ctx context.Context, request *proto_oidc_user.UnlinkUsersRequest) (*proto_oidc_user.UnlinkUsersResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateUnlinkUsersRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateLinkUsersRequest")
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	getUserResponse, err := s.GetUser(ctx, &proto_oidc_user.GetUserRequest{
		Subject: request.RootSubject,
	})
	if err != nil {
		return nil, err
	}
	user := getUserResponse.User

	linkedIdentities := make([]*proto_oidc_models.Identity, 0)

	linq.From(s.userMap).WhereT(func(c *proto_oidc_models.Identity) bool {
		if c.Subject == request.ExternalIdentity.Subject &&
			c.IdpSlug == request.ExternalIdentity.IdpSlug {
			// cull it
			return false
		}
		return true
	}).SelectT(func(c *proto_oidc_models.Identity) *proto_oidc_models.Identity {
		return c
	}).ToSlice(&linkedIdentities)
	user.LinkedIdentities.Identities = linkedIdentities
	s.userMap[user.RootIdentity.Subject] = user
	return &proto_oidc_user.UnlinkUsersResponse{
		User: s.makeUserCopy(user),
	}, nil

}
