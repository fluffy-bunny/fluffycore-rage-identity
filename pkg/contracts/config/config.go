package config

import (
	"strings"

	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_contracts_config "github.com/fluffy-bunny/fluffycore/contracts/config"
	fluffycore_echo_contracts_cookies "github.com/fluffy-bunny/fluffycore/echo/contracts/cookies"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
)

type (
	JWTValidators struct {
		Issuers  []string `json:"issuers"`
		JWKSURLS []string `json:"jwksUrls"`
	}
	ConfigFiles struct {
		OIDCClientPath     string `json:"oidcClientPath"`
		IDPsPath           string `json:"idpsPath"`
		SigningKeyJsonPath string `json:"signingKeyJsonPath"`
		RagePath           string `json:"ragePath"`
		SeedUsersPath      string `json:"seedUsersPath"`
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
	Port                 int                                                    `json:"port"`
	SecureCookies        *fluffycore_echo_contracts_cookies.SecureCookiesConfig `json:"secureCookies"`
	DisableSecureCookies bool                                                   `json:"disableSecureCookies"`
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

type SelfIDPConfig struct {
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	RedirectURL  string   `json:"redirectUrl"`
	Authority    string   `json:"authority"`
	Scopes       []string `json:"scopes"`
}

type OIDCConfig struct {
	BaseUrl            string `json:"baseUrl"`
	OAuth2CallbackPath string `json:"oauth2CallbackPath"`
}
type CookieConfig struct {
	Domain string `json:"domain"`
}
type SystemConfig struct {
	DeveloperMode bool   `json:"developerMode"`
	Domain        string `json:"domain"`
}
type InitialConfig struct {
	ConfigFiles ConfigFiles `json:"configFiles"`
}

type Config struct {
	fluffycore_contracts_config.CoreConfig `mapstructure:",squash"`

	ConfigFiles               ConfigFiles                       `json:"configFiles"`
	Echo                      *EchoConfig                       `json:"echo"`
	InMemoryClients           InMemoryClients                   `json:"inMemoryClients"`
	OIDCConfig                *OIDCConfig                       `json:"oidcConfig"`
	BackingCache              *BackingCacheConfig               `json:"backingCache"`
	AutolinkOnEmailMatch      bool                              `json:"autolinkOnEmailMatch"`
	EmailVerificationRequired bool                              `json:"emailVerificationRequired"`
	EmailConfig               *contracts_email.EmailConfig      `json:"emailConfig"`
	SelfIDPConfig             *SelfIDPConfig                    `json:"selfIDPConfig"`
	CookieConfig              *CookieConfig                     `json:"cookieConfig"`
	SystemConfig              *SystemConfig                     `json:"systemConfig"`
	SessionConfig             *contracts_sessions.SessionConfig `json:"sessionConfig"`
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
	"GRPC_GATEWAY_ENABLED": true,
	"cookieConfig": {
 	},
	"jwtValidators": {},
	"autolinkOnEmailMatch": true,
	"emailVerificationRequired": true,
	"selfIDPConfig": {
		"clientId": "self-client",
		"clientSecret": "secret",
		"redirectUrl": "http://localhost:9044/auth/callback",
		"authority": "http://localhost:9044",
		"scopes": ["openid", "profile", "email"]
	},
	"emailConfig": {
		"justLogIt": false,
		"fromName": "IN_ENVIRONMENT",
		"fromEmail": "IN_ENVIRONMENT@example.com",
		"host": "localhost:25",
		"auth": {
			"plainAuth": {
				"identity": "",
				"username": "",
				"password": "",
				"host": "localhost"
			}
		}
	},
	"backingCache": {
		"type": "${{BACKING_CACHE_TYPE}}",
		"inMemoryCache": {
			"defaultExpirationSeconds": -1,
			"cleanupIntervalSeconds": 60
		}
	},
	"configFiles": {
		"ragePath": "./config/rage.json",
		"oidcClientPath": "./config/oidc-clients.json",
		"idpsPath": "./config/idps.json",
		"signingKeyJsonPath": "./config/signing-keys.json",
		"seedUsersPath": "./config/seed-users.json"
	},
	"ddProfilerConfig": {
		"ENABLED": false,
		"SERVICE_NAME": "in-environment",
		"APPLICATION_ENVIRONMENT": "in-environment",
		"VERSION": "1.0.0"
	},
	"echo": {
		"port": 0,
		"disableSecureCookies": false,
		"secureCookies": {
			"hashKey": "7f6a8b9c0d1e2f3a4b5c6d7e8f9a0b1c",
			"blockKey": "1234567890abcdef1234567890abcdef"
		}
	},
	"inMemoryClients": {
		"clients": []
	},	
	"oidcConfig": {
		"baseUrl": "IN_ENVIRONMENT",
		"oauth2CallbackPath": "/oauth2/callback"
	},
	"sessionConfig": {
		"sessionName": "_session",
		"authenticationKey": "7f6a8b9c0d1e2f3a4b5c6d7e8f9a0b1c",
		"encryptionKey": "1234567890abcdef1234567890abcdef",
		"maxAge": 1800,
		"domain": "@@C"
	},
	"systemConfig": {
		"domain": "@@CHANGEME@@",
		"developerMode": false
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
