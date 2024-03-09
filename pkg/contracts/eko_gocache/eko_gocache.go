package eko_cache

import (
	fluffycore_contracts_eko_gocache "github.com/fluffy-bunny/fluffycore/contracts/eko_gocache"
)

type (
	IAuthorizationRequestStateCache interface {
		fluffycore_contracts_eko_gocache.IGoCache
	}

	IExternalOAuth2Cache interface {
		fluffycore_contracts_eko_gocache.IGoCache
	}
)
