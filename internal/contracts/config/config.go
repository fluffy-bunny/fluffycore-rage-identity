package config

import (
	"strings"

	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	fluffycore_contracts_config "github.com/fluffy-bunny/fluffycore/contracts/config"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
)

type (
	JWTValidators struct {
		Issuers  []string `json:"issuers"`
		JWKSURLS []string `json:"jwksUrls"`
	}
	ConfigFiles struct {
		MockOAuth2ClientPath string `json:"mockOAuth2ClientPath"`
		OIDCClientPath       string `json:"oidcClientPath"`
		IDPsPath             string `json:"idpsPath"`
	}
	InMemoryClient struct {
		Secret   string `json:"secret"`
		ClientId string `json:"clientId"`
	}
	InMemoryClients struct {
		Clients []*proto_oidc_models.Client `json:"clients"`
	}
)
type EchoConfig struct {
	Port int `json:"port"`
}

const (
	BackingCacheTypeInMemory = "in-memory"
	BackingCacheTypeRedis    = "redis"
)

type (
	InMemoryCacheConfig struct {
		DefaultExpirationSeconds int `json:"defaultExpirationSeconds"`
		CleanupIntervalSeconds   int `json:"cleanupIntervalSeconds"`
	}
	BackingCacheConfig struct {
		Type          string              `json:"type"`
		InMemoryCache InMemoryCacheConfig `json:"inMemoryCache"`
	}
)

type Config struct {
	fluffycore_contracts_config.CoreConfig `mapstructure:",squash"`

	ConfigFiles      ConfigFiles                             `json:"configFiles"`
	CustomString     string                                  `json:"customString"`
	SomeSecret       string                                  `json:"someSecret" redact:"true"`
	OAuth2Port       int                                     `json:"oauth2Port"`
	JWTValidators    JWTValidators                           `json:"jwtValidators"`
	DDProfilerConfig *fluffycore_contracts_ddprofiler.Config `json:"ddProfilerConfig"`
	Echo             EchoConfig                              `json:"echo"`
	InMemoryClients  InMemoryClients                         `json:"inMemoryClients"`
	// BaseUrl is the base url for the application.  Hardened as opposed to getting it from the request
	BaseUrl                   string             `json:"baseUrl"`
	BackingCache              BackingCacheConfig `json:"backingCache"`
	AutolinkOnEmailMatch      bool               `json:"autolinkOnEmailMatch"`
	EmailVerificationRequired bool               `json:"emailVerificationRequired"`
}

// ConfigDefaultJSON default json
const configDefaultJSONTemplate = `
{
	"APPLICATION_NAME": "in-environment",
	"APPLICATION_ENVIRONMENT": "in-environment",
	"PRETTY_LOG": false,
	"LOG_LEVEL": "info",
	"PORT": 50051,
	"REST_PORT": 50052,
	"oauth2Port": 50053,
	"customString": "some default value",
	"someSecret": "password",
	"GRPC_GATEWAY_ENABLED": true,
	"baseUrl": "[in-environment]",
	"jwtValidators": {},
	"autolinkOnEmailMatch": true,
	"emailVerificationRequired": true,
	"backingCache": {
		"type": "${{BACKING_CACHE_TYPE}}",
		"inMemoryCache": {
			"defaultExpirationSeconds": -1,
			"cleanupIntervalSeconds": 60
		}
	},
	"configFiles": {
		"mockOAuth2ClientPath": "./config/mockOAuth2Clients.json",
		"oidcClientPath": "./config/oidcClients.json",
		"idpsPath": "./config/idps.json"

	},
	"ddProfilerConfig": {
		"ENABLED": false,
		"SERVICE_NAME": "in-environment",
		"APPLICATION_ENVIRONMENT": "in-environment",
		"VERSION": "1.0.0"
	},
	"echo": {
		"port": 9044 
	},
	"inMemoryClients": {
		"clients": []
	}


  }
`

var ConfigDefaultJSON = []byte(``)

func init() {
	replaceMap := map[string]string{
		"${{BACKING_CACHE_TYPE}}": BackingCacheTypeInMemory,
	}
	fixed := configDefaultJSONTemplate
	for k, v := range replaceMap {
		fixed = strings.Replace(fixed, k, v, -1)
	}
	ConfigDefaultJSON = []byte(fixed)
}
