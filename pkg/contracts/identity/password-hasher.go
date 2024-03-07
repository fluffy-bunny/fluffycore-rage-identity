package identity

import "context"

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
		HashPassword(ctx context.Context, request *HashPasswordRequest) (*HashPasswordResponse, error)
		VerifyPassword(ctx context.Context, request *VerifyPasswordRequest) error
	}
)
