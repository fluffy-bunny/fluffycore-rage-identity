package identity

import (
	"context"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
)

type (
	HashPasswordRequest struct {
		Password string `json:"password" validate:"required"`
	}
	HashPasswordResponse struct {
		Password       string `json:"password"`
		HashedPassword string `json:"hashedPassword"`
	}
	VerifyPasswordRequest struct {
		HashedPassword string `json:"hashedPassword" validate:"required"`
		Password       string `json:"password" validate:"required"`
	}
	IPasswordHasher interface {
		// IsAcceptablePassword checks if the password is acceptable.  i.e. not the same as the username, and meets the minimum requirements
		IsAcceptablePassword(user *proto_oidc_models.RageUser, password string) error
		// HashPassword hashes the password
		HashPassword(ctx context.Context, request *HashPasswordRequest) (*HashPasswordResponse, error)
		// VerifyPassword verifies the password
		VerifyPassword(ctx context.Context, request *VerifyPasswordRequest) error
	}
)
