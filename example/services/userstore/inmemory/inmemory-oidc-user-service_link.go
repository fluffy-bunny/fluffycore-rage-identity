package inmemory

import (
	"context"

	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *service) validateLinkUserRequest(request *proto_oidc_user.LinkRageUserRequest) error {
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

func (s *service) LinkRageUser(ctx context.Context, request *proto_oidc_user.LinkRageUserRequest) (*proto_oidc_user.LinkRageUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateLinkUserRequest(request)
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

	isUserLinked := false
	// user cannot be linked to any other account
	getRageUserResponse, err := s.GetRageUser(ctx,
		&proto_oidc_user.GetRageUserRequest{
			By: &proto_oidc_user.GetRageUserRequest_ExternalIdentity{
				ExternalIdentity: request.ExternalIdentity,
			},
		})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			// user is not linked
			err = nil
		} else {
			return nil, err
		}
	} else {
		isUserLinked = true
	}

	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	if isUserLinked {
		// is this a problem?  could be the same user
		alreadyLinkedUser := getRageUserResponse.User

		if alreadyLinkedUser.RootIdentity.Subject != user.RageUser.RootIdentity.Subject {
			return nil, status.Error(codes.AlreadyExists, "External Identity already linked to another user")
		}
	} else {
		// add the link
		now := timestamppb.Now()
		if request.ExternalIdentity.CreatedOn == nil {
			request.ExternalIdentity.CreatedOn = now
		}
		request.ExternalIdentity.UpdatedOn = now
		if user.RageUser.LinkedIdentities == nil {
			user.RageUser.LinkedIdentities = &proto_oidc_models.LinkedIdentities{}
		}
		user.RageUser.LinkedIdentities.Identities = append(user.RageUser.LinkedIdentities.Identities, request.ExternalIdentity)
	}
	s.userMap[user.Id] = user
	return &proto_oidc_user.LinkRageUserResponse{
		User: user.RageUser,
	}, nil

}
