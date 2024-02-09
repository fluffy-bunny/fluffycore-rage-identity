package services

import (
	"context"
	"os"
	"reflect"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/config"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	services_client_inmemory "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/client/inmemory"
	services_oidcflowstore "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/oidcflowstore"
	services_tokenservice "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/tokenservice"
	services_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/util"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-hanko-oidc/proto/oidc/models"
	fluffycore_services_eko_gocache_go_cache "github.com/fluffy-bunny/fluffycore/services/eko_gocache/go_cache"
	zerolog "github.com/rs/zerolog"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

// put all services you want shared between the echo and grpc servers here
// NOTE: they are NOT the same instance, but they are the same type in context of the server.
func ConfigureServices(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) {
	// this has to be added FIRST as it sets up the default inmemory version of the IClient stores
	// it addes an empty *stores_services_client_inmemory.Clients
	services_client_inmemory.AddSingletonIFluffyCoreClientServiceServer(builder)

	services_tokenservice.AddSingletonITokenService(builder)
	services_util.AddSingletonISomeUtil(builder)
	fluffycore_services_eko_gocache_go_cache.AddISingletonInMemoryCache(builder, reflect.TypeOf((*contracts_eko_gocache.IOIDCFlowCache)(nil)))
	services_oidcflowstore.AddSingletonIOIDCFlowCache(builder)
	OnConfigureServicesLoadOIDCClients(ctx, config, builder)
}

func OnConfigureServicesLoadOIDCClients(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnPreServerStartupLoadOIDCClients").Logger()
	oidcClientClientsJSON, err := os.ReadFile(config.ConfigFiles.OIDCClientPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read OIDCClientPath - may not be a problem if clients are comming from a DB")
		return nil
	}
	var oidcClients *proto_oidc_models.Clients = &proto_oidc_models.Clients{}
	err = protojson.Unmarshal(oidcClientClientsJSON, oidcClients)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OIDCClientPath")
		return err
	}
	di.AddSingleton[*proto_oidc_models.Clients](builder, func() *proto_oidc_models.Clients {
		return oidcClients

	})
	return nil

}
