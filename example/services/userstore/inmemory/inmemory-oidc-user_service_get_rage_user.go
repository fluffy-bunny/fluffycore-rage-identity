package inmemory

import (
	"context"

	golinq "github.com/ahmetb/go-linq/v3"
	echo_utils "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/utils"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

func (s *service) validateGetRageUserRequest(request *proto_oidc_user.GetRageUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	switch by := request.By.(type) {
	case *proto_oidc_user.GetRageUserRequest_Subject:
		if fluffycore_utils.IsEmptyOrNil(by.Subject) {
			return status.Error(codes.InvalidArgument, "request.Subject is required")
		}
	case *proto_oidc_user.GetRageUserRequest_Email:
		if fluffycore_utils.IsEmptyOrNil(by.Email) {
			return status.Error(codes.InvalidArgument, "request.Email is required")
		}
		_, ok := echo_utils.IsValidEmailAddress(by.Email)
		if !ok {
			return status.Error(codes.InvalidArgument, "request.Email is not valid")
		}
	case *proto_oidc_user.GetRageUserRequest_ExternalIdentity:
		if fluffycore_utils.IsEmptyOrNil(by.ExternalIdentity.Subject) && fluffycore_utils.IsEmptyOrNil(by.ExternalIdentity.IdpSlug) {
			return status.Error(codes.InvalidArgument, "request.ExternalIdentity.Subject and request.ExternalIdentity.IdpSlug are required")
		}
	}

	return nil

}

func (s *service) makeRageUserCopy(user *proto_oidc_models.RageUser) *proto_oidc_models.RageUser {
	if user == nil {
		return nil
	}
	d, err := protojson.Marshal(user)
	if err != nil {
		return nil
	}
	var newUser proto_oidc_models.RageUser
	err = protojson.Unmarshal(d, &newUser)
	if err != nil {
		return nil
	}
	return &newUser
}
func (s *service) GetRageUser(ctx context.Context, request *proto_oidc_user.GetRageUserRequest) (*proto_oidc_user.GetRageUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetRageUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateGetUserRequest")
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	users := make([]*proto_oidc_models.RageUser, 0)
	for _, v := range s.userMap {
		users = append(users, s.makeRageUserCopy(v.RageUser))
	}
	var rows []*proto_oidc_models.RageUser

	golinq.
		From(users).
		WhereT(func(c *proto_oidc_models.RageUser) bool {
			switch by := request.By.(type) {
			case *proto_oidc_user.GetRageUserRequest_Subject:
				if c.RootIdentity.Subject != by.Subject {
					return false
				}
			case *proto_oidc_user.GetRageUserRequest_Email:
				if c.RootIdentity.Email != by.Email {
					return false
				}
			case *proto_oidc_user.GetRageUserRequest_ExternalIdentity:
				if c.LinkedIdentities == nil ||
					fluffycore_utils.IsEmptyOrNil(c.LinkedIdentities.Identities) {
					return false
				}
				for _, linkedIdentity := range c.LinkedIdentities.Identities {
					if linkedIdentity.Subject != by.ExternalIdentity.Subject ||
						linkedIdentity.IdpSlug != by.ExternalIdentity.IdpSlug {
						return false
					}
				}
			}
			return true
		}).SelectT(func(c *proto_oidc_models.RageUser) *proto_oidc_models.RageUser {
		return c
	}).ToSlice(&rows)
	if len(rows) > 0 {
		return &proto_oidc_user.GetRageUserResponse{
			User: rows[0],
		}, nil
	}
	return nil, status.Error(codes.NotFound, "User not found")
}
