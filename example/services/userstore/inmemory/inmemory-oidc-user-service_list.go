package inmemory

import (
	"context"
	"fmt"
	"strings"

	linq "github.com/ahmetb/go-linq/v3"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
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
		if c.RageUser == nil {
			return false
		}
		if c.RageUser.LinkedIdentities == nil {
			c.RageUser.LinkedIdentities = &proto_oidc_models.LinkedIdentities{}
		}
		if request.Filter != nil {
			if request.Filter.Id != nil {
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Id.Eq) {
					if request.Filter.Id.Eq != c.Id {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Id.In) {
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

			if request.Filter.Email != nil {
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Email.Eq) {
					if request.Filter.Email.Eq != c.RageUser.RootIdentity.Email {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Email.In) {
					found := false
					for _, v := range request.Filter.Email.In {
						if v == c.RageUser.RootIdentity.Email {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Email.Contains) {
					if !strings.Contains(c.RageUser.RootIdentity.Email, request.Filter.Email.Contains) {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.Email.Ne) {
					if request.Filter.Email.Ne == c.RageUser.RootIdentity.Email {
						return false
					}
				}
			}
			if request.Filter.LinkedIdentityIdpSlug != nil && request.Filter.LinkedIdentitySubject == nil {
				// no subject idp combo
				return false
			}
			if request.Filter.LinkedIdentitySubject != nil && request.Filter.LinkedIdentityIdpSlug == nil {
				// no subject idp combo
				return false
			}
			// build a fast lookup
			mapIdpSlugSubjectMap := make(map[string]*proto_oidc_models.Identity)
			mapEmailMap := make(map[string]*proto_oidc_models.Identity)
			makeSubjectIpdSlugKey := func(subject string, idpSlug string) string {
				return fmt.Sprintf("%s:%s", subject, idpSlug)
			}
			if (request.Filter.LinkedIdentityIdpSlug != nil && request.Filter.LinkedIdentitySubject != nil) ||
				request.Filter.LinkedIdentityEmail != nil {
				for _, linkedIdentity := range c.RageUser.LinkedIdentities.Identities {
					subjectIdpSlugKey := makeSubjectIpdSlugKey(linkedIdentity.Subject, linkedIdentity.IdpSlug)
					mapIdpSlugSubjectMap[subjectIdpSlugKey] = linkedIdentity
					mapEmailMap[linkedIdentity.Email] = linkedIdentity
				}
			}

			if request.Filter.LinkedIdentityIdpSlug != nil && request.Filter.LinkedIdentitySubject != nil {
				key := makeSubjectIpdSlugKey(request.Filter.LinkedIdentitySubject.Eq, request.Filter.LinkedIdentityIdpSlug.Eq)
				_, found := mapIdpSlugSubjectMap[key]
				if !found {
					return false
				}
			}

			if request.Filter.LinkedIdentityEmail != nil {
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.LinkedIdentityEmail.Eq) {
					_, found := mapEmailMap[request.Filter.LinkedIdentityEmail.Eq]
					if !found {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.LinkedIdentityEmail.In) {
					found := false
					for _, v := range request.Filter.LinkedIdentityEmail.In {
						_, found = mapEmailMap[v]
						if found {
							break
						}
					}
					if !found {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.LinkedIdentityEmail.Contains) {
					found := false
					for k := range mapEmailMap {
						if strings.Contains(k, request.Filter.LinkedIdentityEmail.Contains) {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
				if fluffycore_utils.IsNotEmptyOrNil(request.Filter.LinkedIdentityEmail.Ne) {
					_, found := mapEmailMap[request.Filter.LinkedIdentityEmail.Ne]
					if found {
						return false
					}
				}
			}

		}
		return true

	}).SelectT(func(c *proto_external_models.ExampleUser) *proto_external_models.ExampleUser {
		if c.RageUser.TOTP != nil {
			// TODO: THIS MUST BE ENCRYPTED AT REST
			// if there is a TOTP secret, we need to return it
			// it better be encrypted at rest
		}
		return s.makeExampleUserCopy(c)
	}).ToSlice(&users)

	return &proto_external_user.ListUserResponse{
		Users: users,
	}, nil
}
