package config

import (
	management_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
	oidc_login_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/contracts/config"
	fluffycore_contracts_config "github.com/fluffy-bunny/fluffycore/contracts/config"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
)

type (
	JWTValidators struct {
		Issuers  []string `json:"issuers"`
		JWKSURLS []string `json:"jwksUrls"`
	}
	InitialConfig struct {
		ConfigFiles ConfigFiles `json:"configFiles"`
	}
	ConfigFiles struct {
		MyAppPath string `json:"myAppPath"`
	}
	Config struct {
		fluffycore_contracts_config.CoreConfig `mapstructure:",squash"`

		ConfigFiles         ConfigFiles                             `json:"configFiles"`
		DDProfilerConfig    *fluffycore_contracts_ddprofiler.Config `json:"ddProfilerConfig" mapstructure:"DD_PROFILER_CONFIG"`
		JWTValidators       JWTValidators                           `json:"jwtValidators"`
		ManagementAppConfig *management_contracts_config.AppConfig  `json:"managementAppConfig"`
		OIDCLoginAppConfig  *oidc_login_contracts_config.AppConfig  `json:"oidcLoginAppConfig"`
	}
)

// ConfigDefaultJSON default json
var ConfigDefaultJSON = []byte(`
{
    "APPLICATION_NAME": "in-environment",
    "APPLICATION_ENVIRONMENT": "in-environment",
    "PORT": 50051,
    "GRPC_GATEWAY_ENABLED": true,
    "REST_PORT": 50052,
    "PRETTY_LOG": false,
    "LOG_LEVEL": "info",
    "emailTemplateConfig": {
        "homeUrl": "in-environment",
        "supportEmail": "support@rage.com"
    },
    "DD_PROFILER_CONFIG": {
        "ENABLED": false,
        "SERVICE_NAME": "in-environment",
        "APPLICATION_ENVIRONMENT": "in-environment",
        "VERSION": "1.0.0"
    },
    "jwtValidators": {
        "issuers": [],
        "jwksUrls": []
    },
    "configFiles": {
        "myAppPath": "./config/myapp.json"
    },
    "oidcFlowCollectionConfig": {
        "databaseName": "mastodon",
        "collectionName": "mastodon_identity_state"
    },
    "mongoMastodonIdentityCollectionConfig": {
        "databaseName": "mastodon",
        "collectionName": "mastodon_identity"
    },
    "orgCacheTTLSeconds": 3600,
    "managementAppConfig": {
        "basehref": "management",
        "bannerBranding": {
            "title": "RAGE Identity Management",
            "logoUrl": "web/m_logo.svg",
            "showBannerVersion": false
        },
        "returnUrl": "http://{environment}/management/",
        "rageBaseUrl": "http://{environment}",
        "accountManagementBaseUrl": "http://{environment}"
    },
    "oidcLoginAppConfig": {
	    "basehref": "oidc-login",
        "bannerBranding": {
            "title": "RAGE Identity",
            "logoUrl": "web/m_logo.svg",
            "showBannerVersion": false
        },
        "rageBaseUrl": "http://{environment}"
    }
}
`)
