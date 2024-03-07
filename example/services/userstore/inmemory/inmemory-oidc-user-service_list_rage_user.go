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

func (s *service) validateListRageUsersRequest(request *proto_oidc_user.ListRageUsersRequest) error {
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

func (s *service) ListRageUsers(ctx context.Context, request *proto_oidc_user.ListRageUsersRequest) (*proto_oidc_user.ListRageUsersResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateListRageUsersRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateListRageUserRequest")
		return nil, err
	}
	filter := &proto_external_models.ExampleUserFilter{}
	if request.Filter != nil {
		if request.Filter.RootSubject != nil {
			filter.Id = &proto_types.IDFilterExpression{
				Eq: request.Filter.RootSubject.Eq,
				In: request.Filter.RootSubject.In,
			}
		}
		if request.Filter.RootEmail != nil {
			filter.Email = &proto_types.StringFilterExpression{
				Eq:       request.Filter.RootEmail.Eq,
				In:       request.Filter.RootEmail.In,
				Contains: request.Filter.RootEmail.Contains,
				Ne:       request.Filter.RootEmail.Ne,
			}
		}
		if request.Filter.LinkedIdentitySubject != nil {
			filter.LinkedIdentitySubject = &proto_types.IDFilterExpression{
				Eq: request.Filter.LinkedIdentitySubject.Eq,
				In: request.Filter.LinkedIdentitySubject.In,
			}
		}
		if request.Filter.LinkedIdentityIdpSlug != nil {
			filter.LinkedIdentityIdpSlug = &proto_types.IDFilterExpression{
				Eq: request.Filter.LinkedIdentityIdpSlug.Eq,
				In: request.Filter.LinkedIdentityIdpSlug.In,
			}
		}
		if request.Filter.LinkedIdentityEmail != nil {
			filter.LinkedIdentityEmail = &proto_types.StringFilterExpression{
				Eq:       request.Filter.LinkedIdentityEmail.Eq,
				In:       request.Filter.LinkedIdentityEmail.In,
				Contains: request.Filter.LinkedIdentityEmail.Contains,
				Ne:       request.Filter.LinkedIdentityEmail.Ne,
			}
		}
	}
	listUserResponse, err := s.ListUser(ctx, &proto_external_user.ListUserRequest{
		Pagination: request.Pagination,
		Filter:     filter,
	})
	if err != nil {
		return nil, err
	}

	users := make([]*proto_oidc_models.RageUser, 0)
	for _, v := range listUserResponse.Users {
		users = append(users, s.makeRageUserCopy(v.RageUser))
	}

	return &proto_oidc_user.ListRageUsersResponse{
		Users: users,
	}, nil
}
