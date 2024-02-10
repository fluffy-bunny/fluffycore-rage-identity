package eko_cache

import (
	"context"

	models "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/models"
	fluffycore_contracts_eko_gocache "github.com/fluffy-bunny/fluffycore/contracts/eko_gocache"
)

type (
	IOIDCFlowCache interface {
		fluffycore_contracts_eko_gocache.IGoCache
	}
	IOIDCFlowStore interface {
		StoreAuthorizationFinal(ctx context.Context, code string, value *models.AuthorizationFinal) error
		GetAuthorizationFinal(ctx context.Context, code string) (*models.AuthorizationFinal, error)
		DeleteAuthorizationFinal(ctx context.Context, code string) error
	}
	IExternalOAuth2Cache interface {
		fluffycore_contracts_eko_gocache.IGoCache
	}
	IExternalOauth2FlowStore interface {
		StoreExternalOauth2Final(ctx context.Context, state string, value *models.ExternalOauth2Final) error
		GetExternalOauth2Final(ctx context.Context, state string) (*models.ExternalOauth2Final, error)
		DeleteExternalOauth2Final(ctx context.Context, state string) error
	}
)
