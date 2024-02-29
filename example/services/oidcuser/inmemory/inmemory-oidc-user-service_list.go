package inmemory

import (
	"context"
	"strings"

	linq "github.com/ahmetb/go-linq/v3"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateListUserRequest(request *proto_external_user.ListUserRequest) error {
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

func (s *service) ListUser(ctx context.Context, request *proto_external_user.ListUserRequest) (*proto_external_user.ListUserResponse, error) {
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

	users := make([]*proto_external_models.ExampleUser, 0)
	for _, v := range s.userMap {
		users = append(users, s.makeExampleUserCopy(v))
	}
	linq.From(users).WhereT(func(c *proto_external_models.ExampleUser) bool {
		if request.Filter != nil {
			if request.Filter.Id != nil {
				if fluffycore_utils.IsEmptyOrNil(request.Filter.Id.Eq) {
					if request.Filter.Id.Eq != c.Id {
						return false
					}
				}
				if fluffycore_utils.IsEmptyOrNil(request.Filter.Id.In) {
					found := false
					for _, v := range request.Filter.Id.In {
						if v == c.Id {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}

			}
			if request.Filter.RageUser != nil {
				if c.RageUser == nil {
					return false
				}
				if request.Filter.RageUser.State != nil {
					if request.Filter.RageUser.State.Eq != c.RageUser.State {
						return false
					}
				}

				if request.Filter.RageUser.RootIdentity != nil {
					if c.RageUser.RootIdentity == nil {
						return false
					}
					c.RageUser.RootIdentity.Email = strings.ToLower(c.RageUser.RootIdentity.Email)
					if request.Filter.RageUser.RootIdentity.Email != nil {
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.RootIdentity.Email.Eq) {
							if strings.ToLower(request.Filter.RageUser.RootIdentity.Email.Eq) != c.RageUser.RootIdentity.Email {
								return false
							}
						}
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.RootIdentity.Email.Contains) {
							if !strings.Contains(c.RageUser.RootIdentity.Email, strings.ToLower(request.Filter.RageUser.RootIdentity.Email.Contains)) {
								return false
							}
						}
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.RootIdentity.Email.Ne) {
							if request.Filter.RageUser.RootIdentity.Email.Ne == c.RageUser.RootIdentity.Email {
								return false
							}
						}
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.RootIdentity.Email.In) {
							found := false
							for _, v := range request.Filter.RageUser.RootIdentity.Email.In {
								if strings.ToLower(v) == c.RageUser.RootIdentity.Email {
									found = true
									break
								}
							}
							if !found {
								return false
							}
						}

					}
				}
				if request.Filter.RageUser.LinkedIdentity != nil {
					if c.RageUser.LinkedIdentities == nil || fluffycore_utils.IsEmptyOrNil(c.RageUser.LinkedIdentities.Identities) {
						return false
					}

					for _, v := range c.RageUser.LinkedIdentities.Identities {
						v.Email = strings.ToLower(v.Email)
						v.IdpSlug = strings.ToLower(v.IdpSlug)

						if request.Filter.RageUser.LinkedIdentity.IdpSlug != nil {
							if strings.ToLower(request.Filter.RageUser.LinkedIdentity.IdpSlug.Eq) != v.IdpSlug {
								return false
							}
						}
						if request.Filter.RageUser.LinkedIdentity.Subject != nil {
							if request.Filter.RageUser.LinkedIdentity.Subject.Eq != v.Subject {
								return false
							}
						}
						if request.Filter.RageUser.LinkedIdentity.Email != nil {
							if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.LinkedIdentity.Email.Eq) {
								if strings.ToLower(request.Filter.RageUser.LinkedIdentity.Email.Eq) != v.Email {
									return false
								}
							}
							if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.LinkedIdentity.Email.Contains) {
								if !strings.Contains(v.Email, strings.ToLower(request.Filter.RageUser.LinkedIdentity.Email.Contains)) {
									return false
								}
							}
							if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.LinkedIdentity.Email.Ne) {
								if request.Filter.RageUser.LinkedIdentity.Email.Ne == v.Email {
									return false
								}
							}
							if !fluffycore_utils.IsEmptyOrNil(request.Filter.RageUser.LinkedIdentity.Email.In) {
								found := false
								for _, v := range request.Filter.RageUser.LinkedIdentity.Email.In {
									if strings.ToLower(v) == c.RageUser.RootIdentity.Email {
										found = true
										break
									}
								}
								if !found {
									return false
								}
							}
						}
					}
				}
			}
		}
		return true

	}).SelectT(func(c *proto_external_models.ExampleUser) *proto_external_models.ExampleUser {
		return s.makeExampleUserCopy(c)
	}).ToSlice(&users)

	return &proto_external_user.ListUserResponse{
		Users: users,
	}, nil
}
