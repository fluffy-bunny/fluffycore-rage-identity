package inmemory

import (
	"context"

	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateLinkUsersRequest(request *proto_oidc_user.LinkRageUsersRequest) error {
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

func (s *service) LinkUsers(ctx context.Context, request *proto_oidc_user.LinkRageUsersRequest) (*proto_oidc_user.LinkRageUsersResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateLinkUsersRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateLinkUsersRequest")
		return nil, err
	}

	getUserResponse, err := s.GetUser(ctx, &proto_external_user.GetUserRequest{
		Subject: request.RootSubject,
	})
	if err != nil {
		return nil, err
	}
	user := getUserResponse.User

	// user cannot be linked to any other account
	listUserResponse, err := s.ListRageUser(ctx, &proto_oidc_user.ListRageUserRequest{
		Filter: &proto_oidc_models.RageUserFilter{
			LinkedIdentity: &proto_oidc_models.IdentityFilter{
				Subject: &proto_types.IDFilterExpression{
					Eq: request.ExternalIdentity.Subject,
				},
				IdpSlug: &proto_types.IDFilterExpression{
					Eq: request.ExternalIdentity.IdpSlug,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	if len(listUserResponse.Users) > 0 {
		// is this a problem?  could be the same user
		if listUserResponse.Users[0].RootIdentity.Subject != user.RageUser.RootIdentity.Subject {
			return nil, status.Error(codes.AlreadyExists, "External Identity already linked to another user")
		}
	} else {
		// add the link
		if user.RageUser.LinkedIdentities == nil {
			user.RageUser.LinkedIdentities = &proto_oidc_models.LinkedIdentities{}
		}
		user.RageUser.LinkedIdentities.Identities = append(user.RageUser.LinkedIdentities.Identities, request.ExternalIdentity)
	}
	s.userMap[user.Id] = user
	return &proto_oidc_user.LinkRageUsersResponse{
		User: user.RageUser,
	}, nil

}
