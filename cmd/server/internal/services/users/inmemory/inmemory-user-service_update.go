package inmemory

import (
	"context"

	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

func (s *service) validateUpdateUserRequest(request *proto_external_user.UpdateUserRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User) {
		return status.Error(codes.InvalidArgument, "request.User is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.User.Id) {
		return status.Error(codes.InvalidArgument, "request.User.Id is required")
	}

	return nil

}
func (s *service) UpdateUser(ctx context.Context, request *proto_external_user.UpdateUserRequest) (*proto_external_user.UpdateUserResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateUpdateUserRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateUpdateUserRequest")
		return nil, err
	}
	getUserResp, err := s.GetUser(ctx, &proto_external_user.GetUserRequest{
		Subject: request.User.Id,
	})
	if err != nil {
		return nil, err
	}
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	//--~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-//
	user := getUserResp.User
	if request.User.Metadata != nil {
		if user.Metadata == nil {
			user.Metadata = make(map[string]string)
		}
		switch v := request.User.Metadata.Update.(type) {
		case *proto_types.StringMapUpdate_Replace:
			if v.Replace.Value == nil {
				v.Replace.Value = make(map[string]string)
			}
			user.Metadata = v.Replace.Value
		case *proto_types.StringMapUpdate_Granular_:
			if v.Granular.Remove != nil {
				for _, k := range v.Granular.Remove {
					delete(user.Metadata, k)
				}
			}
			if v.Granular.Add != nil {
				for k, v := range v.Granular.Add {
					user.Metadata[k] = v
				}
			}
		}
	}
	if request.User.State != nil {
		user.State = request.User.State.Value
	}
	if request.User.Profile != nil {
		if user.Profile == nil {
			user.Profile = &proto_external_models.Profile{}
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
				user.Profile.Address = &proto_external_models.Address{}
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
	s.userMap[user.Id] = user

	return &proto_external_user.UpdateUserResponse{
		User: s.makeUserCopy(user),
	}, nil
}
