package utils

import (
	"github.com/alexedwards/argon2id"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	"github.com/rs/xid"
)

type (
	PasswordHashSet struct {
		Hash     string `json:"hash"`
		Password string `json:"password"`
	}
)

func GeneratePasswordHash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func ComparePasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

const special string = "!@#$%&*"

func GeneratePassword() (string, error) {
	guid := xid.New().String()
	return guid, nil

}

func GeneratePasswordHashSet(secret *string) (*PasswordHashSet, error) {
	var pass string
	var err error
	if fluffycore_utils.IsNotEmptyOrNil(secret) && fluffycore_utils.IsNotEmptyOrNil(*secret) {
		pass = *secret
	} else {
		pass, err = GeneratePassword()
		if err != nil {
			return nil, err
		}
	}

	hash, err := GeneratePasswordHash(pass)
	if err != nil {
		return nil, err
	}
	return &PasswordHashSet{
		Hash:     hash,
		Password: pass,
	}, nil

}
