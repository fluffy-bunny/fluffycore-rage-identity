package inmemory

import (
	"context"

	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	proto_types_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types/webauthn"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	uuid "github.com/gofrs/uuid"
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
	updateUser := request.User
	doRageUserUpdate := func() error {
		updateRageUser := updateUser.RageUser
		if updateRageUser == nil {
			return nil
		}
		if user.RageUser == nil {
			user.RageUser = &proto_oidc_models.RageUser{
				RootIdentity: &proto_oidc_models.Identity{
					Subject: user.Id,
				},
			}
		}
		rageUser := user.RageUser
		if updateRageUser.State != nil {
			rageUser.State = updateRageUser.State.Value
		}
		doWebAuthNUpdate := func() error {
			webAuthNUpdate := updateRageUser.WebAuthN
			if webAuthNUpdate == nil || webAuthNUpdate.Credentials == nil {
				// nothing to do
				return nil
			}
			if rageUser.WebAuthN == nil {
				rageUser.WebAuthN = &proto_oidc_models.WebAuthN{}
			}
			switch v := webAuthNUpdate.Credentials.Update.(type) {
			case *proto_types_webauthn.CredentialArrayUpdate_DeleteAll:
				if v.DeleteAll.Value {
					rageUser.WebAuthN.Credentials = make([]*proto_types_webauthn.Credential, 0)
				}
			case *proto_types_webauthn.CredentialArrayUpdate_Granular_:
				mapExisting := make(map[uuid.UUID]*proto_types_webauthn.Credential)
				for _, credential := range rageUser.WebAuthN.Credentials {
					aaguid, _ := uuid.FromBytes(credential.Authenticator.AAGUID)
					mapExisting[aaguid] = credential
				}
				for _, aaGUID := range v.Granular.RemoveAAGUIDs {
					aaguid, _ := uuid.FromBytes(aaGUID)
					delete(mapExisting, aaguid)
				}
				for _, credential := range v.Granular.Add {
					aaguid, _ := uuid.FromBytes(credential.Authenticator.AAGUID)
					mapExisting[aaguid] = credential
				}
				rageUser.WebAuthN.Credentials = make([]*proto_types_webauthn.Credential, 0)
				for _, credential := range mapExisting {
					rageUser.WebAuthN.Credentials = append(rageUser.WebAuthN.Credentials, credential)
				}

			}
			return nil
		}
		err := doWebAuthNUpdate()
		if err != nil {
			return err
		}
		doRecoveryUpdate := func() error {
			recoveryUpdate := updateRageUser.Recovery
			if recoveryUpdate == nil {
				// nothing to do
				return nil
			}
			if rageUser.Recovery == nil {
				rageUser.Recovery = &proto_oidc_models.Recovery{}
			}
			recovery := rageUser.Recovery
			if recovery.Email == nil {
				recovery.Email = &proto_oidc_models.Email{}
			}
			if recoveryUpdate.Email != nil {
				if recoveryUpdate.Email.Email != nil {
					recovery.Email.Email = recoveryUpdate.Email.Email.Value
				}
				if recoveryUpdate.Email.EmailVerified != nil {
					recovery.Email.EmailVerified = recoveryUpdate.Email.EmailVerified.Value
				}
			}
			return nil
		}
		err = doRecoveryUpdate()
		if err != nil {
			return err
		}
		doPasswordUpdate := func() error {
			passwordUpdate := updateRageUser.Password
			if passwordUpdate == nil {
				// nothing to do
				return nil
			}
			if rageUser.Password == nil {
				rageUser.Password = &proto_oidc_models.Password{}
			}
			if passwordUpdate.Hash != nil {
				rageUser.Password = &proto_oidc_models.Password{
					Hash: passwordUpdate.Hash.Value,
				}
			}
			return nil
		}
		err = doPasswordUpdate()
		if err != nil {
			return err
		}

		doRootIdentityUpdate := func() error {
			rootIdentityUpdate := updateRageUser.RootIdentity
			if rootIdentityUpdate == nil {
				// nothing to do
				return nil
			}
			if rageUser.RootIdentity == nil {
				rageUser.RootIdentity = &proto_oidc_models.Identity{}
			}
			rootIdentity := rageUser.RootIdentity
			// set the subject no matter what
			rootIdentity.Subject = user.Id

			if rootIdentityUpdate.EmailVerified != nil {
				rootIdentity.EmailVerified = rootIdentityUpdate.EmailVerified.Value
			}
			if fluffycore_utils.IsNotEmptyOrNil(rootIdentityUpdate.Email) {
				rootIdentity.Email = rootIdentityUpdate.Email.Value
			}
			return nil
		}
		err = doRootIdentityUpdate()
		if err != nil {
			return err
		}
		return nil

	}
	err = doRageUserUpdate()
	if err != nil {
		return nil, err
	}

	doProfileUpdate := func() error {
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
		return nil
	}
	err = doProfileUpdate()
	if err != nil {
		return nil, err
	}

	s.userMap[user.Id] = user

	return &proto_external_user.UpdateUserResponse{
		User: s.makeExampleUserCopy(user),
	}, nil
}
