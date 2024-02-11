package services

import (
	"context"
	"os"
	"reflect"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/config"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/eko_gocache"
	services_client_inmemory "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/client/inmemory"
	services_codeexchanges_genericoidc "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/codeexchanges/genericoidc"
	services_codeexchanges_github "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/codeexchanges/github"
	services_idp_inmemory "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/idp/inmemory"
	services_oauth2factory "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/oauth2factory"
	services_oauth2flowstore "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/oauth2flowstore"
	services_oidcflowstore "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/oidcflowstore"
	services_oidcproviderfactory "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/services/oidcproviderfactory"
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
	services_idp_inmemory.AddSingletonIFluffyCoreIDPServiceServer(builder)
	services_oauth2factory.AddSingletonIOAuth2Factory(builder)
	services_tokenservice.AddSingletonITokenService(builder)
	services_codeexchanges_github.AddSingletonIGithubCodeExchange(builder)
	services_codeexchanges_genericoidc.AddSingletonIGenericOIDCCodeExchange(builder)
	services_oidcproviderfactory.AddSingletonIOIDCProviderFactory(builder)
	services_util.AddSingletonISomeUtil(builder)
	fluffycore_services_eko_gocache_go_cache.AddISingletonInMemoryCache(builder,
		reflect.TypeOf((*contracts_eko_gocache.IOIDCFlowCache)(nil)),
		reflect.TypeOf((*contracts_eko_gocache.IExternalOAuth2Cache)(nil)),
	)
	services_oidcflowstore.AddSingletonIOIDCFlowCache(builder)
	services_oauth2flowstore.AddSingletonIExternalOauth2FlowStore(builder)
	di.AddInstance[*contracts_config.Config](builder, config)
	OnConfigureServicesLoadOIDCClients(ctx, config, builder)
	OnConfigureServicesLoadIDPs(ctx, config, builder)
}
func OnConfigureServicesLoadIDPs(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnConfigureServicesLoadIDPs").Logger()
	fileContent, err := os.ReadFile(config.ConfigFiles.IDPsPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read IDPsPath - may not be a problem if idps are comming from a DB")
		return nil
	}
	var idps *proto_oidc_models.IDPs = &proto_oidc_models.IDPs{}
	err = protojson.Unmarshal(fileContent, idps)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OIDCClientPath")
		return err
	}
	di.AddSingleton[*proto_oidc_models.IDPs](builder, func() *proto_oidc_models.IDPs {
		return idps
	})
	return nil

}

func OnConfigureServicesLoadOIDCClients(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnConfigureServicesLoadOIDCClients").Logger()
	fileContent, err := os.ReadFile(config.ConfigFiles.OIDCClientPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read OIDCClientPath - may not be a problem if clients are comming from a DB")
		return nil
	}
	var oidcClients *proto_oidc_models.Clients = &proto_oidc_models.Clients{}
	err = protojson.Unmarshal(fileContent, oidcClients)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OIDCClientPath")
		return err
	}
	di.AddSingleton[*proto_oidc_models.Clients](builder, func() *proto_oidc_models.Clients {
		return oidcClients

	})
	return nil

}
