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

	users := make([]*proto_external_models.User, 0)
	for _, v := range s.userMap {
		users = append(users, s.makeUserCopy(v))
	}
	linq.From(users).WhereT(func(c *proto_external_models.User) bool {
		if c.Metadata == nil {
			c.Metadata = make(map[string]string)
		}
		if request.Filter != nil {
			if request.Filter.State != nil {
				if request.Filter.State.Eq != c.State {
					return false
				}
			}
			if request.Filter.Metadata != nil {
				if !fluffycore_utils.IsEmptyOrNil(request.Filter.Metadata.Key) {
					_, ok := c.Metadata[request.Filter.Metadata.Key]
					if !ok {
						return false
					}
					if request.Filter.Metadata.Value != nil {
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.Metadata.Value.Eq) {
							if c.Metadata[request.Filter.Metadata.Key] != request.Filter.Metadata.Value.Eq {
								return false
							}
						}
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.Metadata.Value.Contains) {
							if !strings.Contains(c.Metadata[request.Filter.Metadata.Key], request.Filter.Metadata.Value.Contains) {
								return false
							}
						}
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.Metadata.Value.In) {
							found := false
							for _, v := range request.Filter.Metadata.Value.In {
								if c.Metadata[request.Filter.Metadata.Key] == v {
									found = true
									break
								}
							}
							if !found {
								return false
							}
						}
						if !fluffycore_utils.IsEmptyOrNil(request.Filter.Metadata.Value.Ne) {
							if c.Metadata[request.Filter.Metadata.Key] == request.Filter.Metadata.Value.Ne {
								return false
							}
						}
					}
				}

			}
			return true
		} else {
			return true
		}

	}).SelectT(func(c *proto_external_models.User) *proto_external_models.User {
		return s.makeUserCopy(c)
	}).ToSlice(&users)

	return &proto_external_user.ListUserResponse{
		Users: users,
	}, nil
}
