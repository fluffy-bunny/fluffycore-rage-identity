package inmemory

import (
	"context"
	"strings"

	linq "github.com/ahmetb/go-linq/v3"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/types"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateListUserRequest(request *proto_oidc_user.ListUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if request.Pagination == nil {
		request.Pagination = &proto_types.Pagination{
			Limit:    100,
			Iterator: "",
			Order:    proto_types.Order_ASC,
		}
	}
	return nil

}

func (s *service) ListUser(ctx context.Context, request *proto_oidc_user.ListUserRequest) (*proto_oidc_user.ListUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateListUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateListUserRequest")
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	users := make([]*proto_oidc_models.User, 0)

	linq.From(s.users).WhereT(func(c *proto_oidc_models.User) bool {
		if request.Filter != nil {
			if request.Filter.State != nil {
				if request.Filter.State.Eq != c.State {
					return false
				}
			}
			if request.Filter.RootIdentity != nil {
				if request.Filter.RootIdentity.Email != nil {
					eqEmail := strings.ToLower(request.Filter.RootIdentity.Email.Eq)
					if eqEmail != c.RootIdentity.Email {
						return false
					}
				}
				if request.Filter.RootIdentity.IdpSlug != nil {
					if request.Filter.RootIdentity.IdpSlug.Eq != c.RootIdentity.IdpSlug {
						return false
					}
				}
			}
			if request.Filter.LinkedIdentity != nil {
				if c.LinkedIdentities == nil {
					return false
				}
				for _, v := range c.LinkedIdentities.Identities {
					if request.Filter.LinkedIdentity.Subject != nil {
						if request.Filter.LinkedIdentity.Subject.Eq != v.Subject {
							return false
						}
					}
					if request.Filter.LinkedIdentity.IdpSlug != nil {
						if request.Filter.LinkedIdentity.IdpSlug.Eq != v.IdpSlug {
							return false
						}
					}
					if request.Filter.LinkedIdentity.Email != nil {
						if request.Filter.LinkedIdentity.Email.Eq != v.Email {
							return false
						}
					}
				}

			}
			return true
		} else {
			return true
		}

	}).SelectT(func(c *proto_oidc_models.User) *proto_oidc_models.User {
		return c
	}).ToSlice(&users)

	return &proto_oidc_user.ListUserResponse{
		Users: users,
	}, nil
}
