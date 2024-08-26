package services

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_eko_gocache "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/eko_gocache"
	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	services_AuthorizationCodeClaimsAugmentor "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/AuthorizationCodeClaimsAugmentor"
	services_client_inmemory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/client/inmemory"
	services_codeexchanges_genericoidc "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/codeexchanges/genericoidc"
	services_codeexchanges_github "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/codeexchanges/github"
	services_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/cookies"
	services_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/email"
	services_emailrenderer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/emailrenderer"
	services_identity_passwordhasher "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/identity/passwordhasher"
	services_identity_userid "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/identity/userid"
	services_idp_inmemory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/idp/inmemory"
	services_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/localizer"
	services_localizerbundle "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/localizerbundle"
	services_oauth2factory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/oauth2factory"
	services_oidcproviderfactory "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/oidcproviderfactory"
	services_selfoauth2provider "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/selfoauth2provider"
	services_tokenservice "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/tokenservice"
	services_util "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/util"
	services_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/webauthn"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_contracts_jwtminter "github.com/fluffy-bunny/fluffycore/contracts/jwtminter"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_echo_services_cookies_insecure "github.com/fluffy-bunny/fluffycore/echo/services/cookies/insecure"
	fluffycore_echo_services_cookies_secure "github.com/fluffy-bunny/fluffycore/echo/services/cookies/secure"
	fluffycore_echo_templates "github.com/fluffy-bunny/fluffycore/echo/templates"
	fluffycore_services_eko_gocache_go_cache "github.com/fluffy-bunny/fluffycore/services/eko_gocache/go_cache"
	fluffycore_services_jwtminter "github.com/fluffy-bunny/fluffycore/services/jwtminter"
	fluffycore_services_keymaterial "github.com/fluffy-bunny/fluffycore/services/keymaterial"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	zerolog "github.com/rs/zerolog"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

// put all services you want shared between the echo and grpc servers here
// NOTE: they are NOT the same instance, but they are the same type in context of the server.
func ConfigureServices(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) error {

	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()
	// this has to be added FIRST as it sets up the default inmemory version of the IClient stores
	// it addes an empty *stores_services_client_inmemory.Clients
	services_client_inmemory.AddSingletonIFluffyCoreClientServiceServer(builder)
	services_idp_inmemory.AddSingletonIFluffyCoreIDPServiceServer(builder)
	services_oauth2factory.AddSingletonIOAuth2Factory(builder)
	services_tokenservice.AddSingletonITokenService(builder)
	services_AuthorizationCodeClaimsAugmentor.AddSingletonIClaimsAugmentor(builder)
	services_codeexchanges_github.AddSingletonIGithubCodeExchange(builder)
	services_codeexchanges_genericoidc.AddSingletonIGenericOIDCCodeExchange(builder)
	services_oidcproviderfactory.AddSingletonIOIDCProviderFactory(builder)
	services_util.AddSingletonISomeUtil(builder)
	services_identity_passwordhasher.AddSingletonIPasswordHasher(builder)
	services_identity_userid.AddSingletonIUserIdGenerator(builder)
	switch config.BackingCache.Type {
	case contracts_config.BackingCacheTypeInMemory:
		inMemoryOptions := &fluffycore_services_eko_gocache_go_cache.InMemoryCacheOptions{
			ImplementedInterfaceTypes: []reflect.Type{
				reflect.TypeOf((*contracts_eko_gocache.IAuthorizationRequestStateCache)(nil)),
				reflect.TypeOf((*contracts_eko_gocache.IExternalOAuth2Cache)(nil)),
			},
		}
		durationPtr := func(duration time.Duration) *time.Duration {
			return &duration
		}
		if config.BackingCache.InMemoryCache.DefaultExpirationSeconds > 0 {
			inMemoryOptions.DefaultExpiration = durationPtr(time.Duration(config.BackingCache.InMemoryCache.DefaultExpirationSeconds) * time.Second)
		}
		if config.BackingCache.InMemoryCache.CleanupIntervalSeconds > 0 {
			inMemoryOptions.CleanupInterval = durationPtr(time.Duration(config.BackingCache.InMemoryCache.CleanupIntervalSeconds) * time.Second)
		}
		fluffycore_services_eko_gocache_go_cache.AddISingletonInMemoryCacheWithOptions(builder, inMemoryOptions)
	}

	di.AddInstance[*contracts_config.Config](builder, config)
	di.AddInstance[*contracts_config.OIDCConfig](builder, config.OIDCConfig)
	di.AddInstance[*contracts_config.SelfIDPConfig](builder, config.SelfIDPConfig)
	di.AddInstance[*contracts_email.EmailConfig](builder, config.EmailConfig)
	di.AddInstance[*contracts_config.EchoConfig](builder, config.Echo)
	di.AddInstance[*contracts_config.BackingCacheConfig](builder, config.BackingCache)
	di.AddInstance[*contracts_config.PasswordConfig](builder, config.PasswordConfig)
	di.AddInstance[*contracts_webauthn.WebAuthNConfig](builder, config.WebAuthNConfig)

	if config.CookieConfig == nil {
		config.CookieConfig = &contracts_config.CookieConfig{}
	}
	config.CookieConfig.Domain = config.SystemConfig.Domain
	di.AddInstance[*contracts_config.CookieConfig](builder, config.CookieConfig)
	if config.SessionConfig == nil {
		config.SessionConfig = &contracts_sessions.SessionConfig{}
	}
	config.SessionConfig.Domain = config.SystemConfig.Domain

	di.AddInstance[*contracts_sessions.SessionConfig](builder, config.SessionConfig)
	di.AddInstance[*contracts_webauthn.WebAuthNConfig](builder, config.WebAuthNConfig)

	OnConfigureServicesLoadOIDCClients(ctx, config, builder)
	OnConfigureServicesLoadIDPs(ctx, config, builder)
	addJwtMinter := func() {
		signingKeys := []*fluffycore_contracts_jwtminter.SigningKey{}
		fileContent, err := os.ReadFile(config.ConfigFiles.SigningKeyJsonPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read signing key file")
		}
		err = json.Unmarshal(fileContent, &signingKeys)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to unmarshal signing key file")
		}
		keymaterial := &fluffycore_contracts_jwtminter.KeyMaterial{
			SigningKeys: signingKeys,
		}
		di.AddInstance[*fluffycore_contracts_jwtminter.KeyMaterial](builder, keymaterial)
		fluffycore_services_keymaterial.AddSingletonIKeyMaterial(builder)
		fluffycore_services_jwtminter.AddSingletonIJWTMinter(builder)
	}
	addJwtMinter()
	templateEngine, err := fluffycore_echo_templates.FindAndParseTemplates("./static/templates_email", nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse email templates")
		return err
	}
	config.EmailConfig.TemplateEngine = templateEngine

	services_localizerbundle.AddSingletonILocalizerBundle(builder)
	services_localizer.AddScopedILocalizer(builder)
	services_email.AddScopedIEmailService(builder)
	services_emailrenderer.AddSingletonIEmailRenderer(builder)
	fluffycore_echo_services_cookies_insecure.AddCookies(builder)
	fluffycore_echo_services_cookies_secure.AddSecureCookies(builder, config.Echo.SecureCookies)
	services_cookies.AddSingletonIWellknownCookies(builder)

	services_selfoauth2provider.AddSingletonISelfOAuth2Provider(builder)
	services_webauthn.AddSingletonIWebAuthN(builder)
	return nil
}
func OnConfigureServicesLoadIDPs(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnConfigureServicesLoadIDPs").Logger()
	fileContent, err := os.ReadFile(config.ConfigFiles.IDPsPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read IDPsPath - may not be a problem if idps are comming from a DB")
		return nil
	}
	fixedFileContent := fluffycore_utils.ReplaceEnv(string(fileContent), "${%s}")
	var idps *proto_oidc_models.IDPs = &proto_oidc_models.IDPs{}
	err = protojson.Unmarshal([]byte(fixedFileContent), idps)
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
	fixedFileContent := fluffycore_utils.ReplaceEnv(string(fileContent), "${%s}")

	var oidcClients *proto_oidc_models.Clients = &proto_oidc_models.Clients{}
	err = protojson.Unmarshal([]byte(fixedFileContent), oidcClients)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal OIDCClientPath")
		return err
	}
	di.AddSingleton[*proto_oidc_models.Clients](builder, func() *proto_oidc_models.Clients {
		return oidcClients

	})
	return nil

}
