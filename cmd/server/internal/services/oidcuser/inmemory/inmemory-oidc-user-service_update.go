package inmemory

import (
	"context"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateUpdateUserRequest(request *proto_oidc_user.UpdateUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User) {
		return status.Error(codes.InvalidArgument, "request.User is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.RootIdentity.Subject) {
		return status.Error(codes.InvalidArgument, "request.User.RootIdentity.Subject is required")
	}
	return nil

}
func (s *service) UpdateUser(ctx context.Context, request *proto_oidc_user.UpdateUserRequest) (*proto_oidc_user.UpdateUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateUpdateUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateUpdateUserRequest")
		return nil, err
	}
	getUserResp, err := s.GetUser(ctx, &proto_oidc_user.GetUserRequest{
		Subject: request.User.RootIdentity.Subject,
	})
	if err != nil {
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	user := getUserResp.User
	if request.User.RootIdentity.EmailVerified != nil {
		user.RootIdentity.EmailVerified = request.User.RootIdentity.EmailVerified.Value
	}
	if !fluffycore_utils.IsEmptyOrNil(request.User.State) {
		user.State = request.User.State.Value
	}
	if request.User.Password != nil {
		if request.User.Password.Hash != nil {
			user.Password = &proto_oidc_models.Password{
				Hash: request.User.Password.Hash.Value,
			}
		}
	}
	if request.User.Profile != nil {
		if user.Profile == nil {
			user.Profile = &proto_oidc_models.Profile{}
		}
		profile := request.User.Profile
		if profile.GivenName != nil {
			user.Profile.GivenName = profile.GivenName.Value
		}
		if profile.FamilyName != nil {
			user.Profile.FamilyName = profile.FamilyName.Value
		}
		if profile.Address != nil {
			if user.Profile.Address == nil {
				user.Profile.Address = &proto_oidc_models.Address{}
			}
			if profile.Address.City != nil {
				user.Profile.Address.City = profile.Address.City.Value
			}
			if profile.Address.Country != nil {
				user.Profile.Address.Country = profile.Address.Country.Value
			}
			if profile.Address.PostalCode != nil {
				user.Profile.Address.PostalCode = profile.Address.PostalCode.Value
			}
			if profile.Address.State != nil {
				user.Profile.Address.State = profile.Address.State.Value
			}
			if profile.Address.Street != nil {
				user.Profile.Address.Street = profile.Address.Street.Value
			}
		}
		if profile.PhoneNumbers != nil {
			if user.Profile.PhoneNumbers == nil {
				user.Profile.PhoneNumbers = make([]*proto_types.PhoneNumberDTO, 0)
			}
			existingPhoneNumbers := make(map[string]*proto_types.PhoneNumberDTO)
			for _, v := range user.Profile.PhoneNumbers {
				existingPhoneNumbers[v.Id] = v
			}
			toAdd := make([]*proto_types.PhoneNumberDTO, 0)
			for _, v := range profile.PhoneNumbers {
				original, ok := existingPhoneNumbers[v.Id]
				if ok {
					if v.CountryCode != nil {
						original.CountryCode = v.CountryCode.Value
					}
					if v.Number != nil {
						original.Number = v.Number.Value
					}
					if v.Type != nil {
						original.Type = v.Type.Value
					}
				} else {
					// add it.
					newPhoneNumber := &proto_types.PhoneNumberDTO{
						Id: v.Id,
					}
					if v.CountryCode != nil {
						newPhoneNumber.CountryCode = v.CountryCode.Value
					}
					if v.Number != nil {
						newPhoneNumber.Number = v.Number.Value
					}
					if v.Type != nil {
						newPhoneNumber.Type = v.Type.Value
					}
					toAdd = append(toAdd, newPhoneNumber)
				}
			}
			user.Profile.PhoneNumbers = append(user.Profile.PhoneNumbers, toAdd...)
		}
	}
	s.userMap[user.RootIdentity.Subject] = user

	return &proto_oidc_user.UpdateUserResponse{
		User: s.makeUserCopy(user),
	}, nil
}
