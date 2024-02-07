package services

import (
	"reflect"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	services_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/util"
	fluffycore_services_eko_gocache_go_cache "github.com/fluffy-bunny/fluffycore/services/eko_gocache/go_cache"
	services_oidcflowstore "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/oidcflowstore"

)

// put all services you want shared between the echo and grpc servers here
// NOTE: they are NOT the same instance, but they are the same type in context of the server.
func ConfigureServices(builder di.ContainerBuilder) {
	services_util.AddSingletonISomeUtil(builder)
	fluffycore_services_eko_gocache_go_cache.AddISingletonInMemoryCache(builder, reflect.TypeOf((*contracts_eko_gocache.IOIDCFlowCache)(nil)))
	services_oidcflowstore.AddSingletonIOIDCFlowCache(builder)
}
