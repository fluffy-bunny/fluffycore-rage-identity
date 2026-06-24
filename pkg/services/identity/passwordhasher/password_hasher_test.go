package passwordhasher_test

import (
	"context"
	"testing"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	passwordhasher "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/identity/passwordhasher"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	require "github.com/stretchr/testify/require"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func buildContainer(t *testing.T, config *contracts_config.PasswordConfig) contracts_identity.IPasswordHasher {
	t.Helper()
	b := di.Builder()
	b.ConfigureOptions(func(o *di.Options) {
		o.ValidateScopes = true
		o.ValidateOnBuild = true
	})
	di.AddInstance[*contracts_config.PasswordConfig](b, config)
	passwordhasher.AddSingletonIPasswordHasher(b)
	container := b.Build()
	hasher := di.Get[contracts_identity.IPasswordHasher](container)
	require.NotNil(t, hasher)
	return hasher
}

func configNoCache() *contracts_config.PasswordConfig {
	return &contracts_config.PasswordConfig{
		MinEntropyBits: 60,
		CacheEnabled:   false,
	}
}

func configWithCache(ttl time.Duration) *contracts_config.PasswordConfig {
	return &contracts_config.PasswordConfig{
		MinEntropyBits: 60,
		CacheEnabled:   true,
		CacheTTL:       fluffycore_utils.Duration(ttl),
	}
}

// grpcCode extracts the gRPC status code from an error, or returns codes.OK for nil.
func grpcCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}
	return codes.Unknown
}

// ---- HashPassword -------------------------------------------------------

func TestHashPassword(t *testing.T) {
	hasher := buildContainer(t, configNoCache())
	ctx := context.Background()

	t.Run("nil request returns InvalidArgument", func(t *testing.T) {
		_, err := hasher.HashPassword(ctx, nil)
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("empty password returns InvalidArgument", func(t *testing.T) {
		_, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: ""})
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("valid password returns non-empty hash", func(t *testing.T) {
		resp, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: "S3cr3tP@ssw0rd!"})
		require.NoError(t, err)
		require.NotEmpty(t, resp.HashedPassword)
		require.Equal(t, "S3cr3tP@ssw0rd!", resp.Password)
	})

	t.Run("two calls produce different hashes (bcrypt salt)", func(t *testing.T) {
		const pw = "S3cr3tP@ssw0rd!"
		r1, err1 := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: pw})
		r2, err2 := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: pw})
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, r1.HashedPassword, r2.HashedPassword, "bcrypt must produce unique salts")
	})
}

// ---- VerifyPassword (no cache) ------------------------------------------

func TestVerifyPassword_NoCache(t *testing.T) {
	hasher := buildContainer(t, configNoCache())
	ctx := context.Background()

	t.Run("nil request returns InvalidArgument", func(t *testing.T) {
		err := hasher.VerifyPassword(ctx, nil)
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("empty password returns InvalidArgument", func(t *testing.T) {
		err := hasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password:       "",
			HashedPassword: "$2a$10$placeholder",
		})
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("empty hashed password returns InvalidArgument", func(t *testing.T) {
		err := hasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password:       "S3cr3tP@ssw0rd!",
			HashedPassword: "",
		})
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("correct password returns nil", func(t *testing.T) {
		const pw = "S3cr3tP@ssw0rd!"
		resp, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: pw})
		require.NoError(t, err)

		err = hasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password:       pw,
			HashedPassword: resp.HashedPassword,
		})
		require.NoError(t, err)
	})

	t.Run("wrong password returns NotFound", func(t *testing.T) {
		resp, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: "S3cr3tP@ssw0rd!"})
		require.NoError(t, err)

		err = hasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password:       "wr0ngP@ssword!",
			HashedPassword: resp.HashedPassword,
		})
		require.Equal(t, codes.NotFound, grpcCode(err))
	})
}

// ---- VerifyPassword (with cache) ----------------------------------------

func TestVerifyPassword_WithCache(t *testing.T) {
	hasher := buildContainer(t, configWithCache(5*time.Minute))
	ctx := context.Background()

	t.Run("correct password cached on second call", func(t *testing.T) {
		const pw = "S3cr3tP@ssw0rd!"
		resp, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: pw})
		require.NoError(t, err)

		req := &contracts_identity.VerifyPasswordRequest{
			Password:       pw,
			HashedPassword: resp.HashedPassword,
		}
		// First call: bcrypt comparison + cache write.
		require.NoError(t, hasher.VerifyPassword(ctx, req))
		// Second call: served from cache.
		require.NoError(t, hasher.VerifyPassword(ctx, req))
	})

	t.Run("wrong password cached on second call", func(t *testing.T) {
		resp, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: "S3cr3tP@ssw0rd!"})
		require.NoError(t, err)

		req := &contracts_identity.VerifyPasswordRequest{
			Password:       "wr0ngP@ssword!",
			HashedPassword: resp.HashedPassword,
		}
		// First call: bcrypt comparison + cache write.
		require.Equal(t, codes.NotFound, grpcCode(hasher.VerifyPassword(ctx, req)))
		// Second call: served from cache, same result.
		require.Equal(t, codes.NotFound, grpcCode(hasher.VerifyPassword(ctx, req)))
	})

	t.Run("different passwords do not share cache entries", func(t *testing.T) {
		const pw = "S3cr3tP@ssw0rd!"
		resp, err := hasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{Password: pw})
		require.NoError(t, err)

		require.NoError(t, hasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password: pw, HashedPassword: resp.HashedPassword,
		}))
		require.Equal(t, codes.NotFound, grpcCode(hasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			Password: "differentP@ss1!", HashedPassword: resp.HashedPassword,
		})))
	})
}

// ---- IsAcceptablePassword -----------------------------------------------

func TestIsAcceptablePassword(t *testing.T) {
	hasher := buildContainer(t, configNoCache())

	t.Run("nil request returns InvalidArgument", func(t *testing.T) {
		err := hasher.IsAcceptablePassword(nil)
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("empty password returns InvalidArgument", func(t *testing.T) {
		err := hasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{Password: ""})
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("email-shaped password is rejected", func(t *testing.T) {
		err := hasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{Password: "user@example.com"})
		require.Equal(t, codes.InvalidArgument, grpcCode(err))
	})

	t.Run("low-entropy password is rejected", func(t *testing.T) {
		err := hasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{Password: "password"})
		require.Error(t, err)
	})

	t.Run("high-entropy password is accepted", func(t *testing.T) {
		err := hasher.IsAcceptablePassword(&contracts_identity.IsAcceptablePasswordRequest{Password: "Tr0ub4dor#3-h0rse-B4tt3ry!"})
		require.NoError(t, err)
	})
}
