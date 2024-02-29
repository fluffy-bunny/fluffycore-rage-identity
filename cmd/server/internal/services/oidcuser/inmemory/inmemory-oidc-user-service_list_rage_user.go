package inmemory

import (
	"context"

	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateListRageUserRequest(request *proto_oidc_user.ListRageUserRequest) error {
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

func (s *service) ListRageUser(ctx context.Context, request *proto_oidc_user.ListRageUserRequest) (*proto_oidc_user.ListRageUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateListRageUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateListRageUserRequest")
		return nil, err
	}

	listUserResponse, err := s.ListUser(ctx, &proto_external_user.ListUserRequest{
		Pagination: request.Pagination,
		Filter: &proto_external_models.ExampleUserFilter{
			RageUser: request.Filter,
		},
	})
	if err != nil {
		return nil, err
	}

	users := make([]*proto_oidc_models.RageUser, 0)
	for _, v := range listUserResponse.Users {
		users = append(users, s.makeRageUserCopy(v.RageUser))
	}

	return &proto_oidc_user.ListRageUserResponse{
		Users: users,
	}, nil
}
