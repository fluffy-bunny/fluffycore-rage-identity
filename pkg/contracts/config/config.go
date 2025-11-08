package config

import (
	"strings"

	contracts_email "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/email"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	models_api_appsettings "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/appsettings"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_contracts_config "github.com/fluffy-bunny/fluffycore/contracts/config"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
	fluffycore_contracts_otel "github.com/fluffy-bunny/fluffycore/contracts/otel"
	fluffycore_echo_contracts_cookies "github.com/fluffy-bunny/fluffycore/echo/contracts/cookies"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
)

type (
	JWTValidators struct {
		Issuers  []string `json:"issuers"`
		JWKSURLS []string `json:"jwksUrls"`
	}
	CORSConfig struct {
		Enabled                                  bool     `json:"enabled"`
		AllowedOrigins                           []string `json:"allowedOrigins"`
		AllowedMethods                           []string `json:"allowedMethods"`
		AllowedHeaders                           []string `json:"allowedHeaders"`
		AllowCredentials                         bool     `json:"allowCredentials"`
		UnsafeWildcardOriginWithAllowCredentials bool     `json:"unsafeWildcardOriginWithAllowCredentials"`
		ExposeHeaders                            []string `json:"exposeHeaders"`
		MaxAge                                   int      `json:"maxAge"`
	}
	CSRFConfig struct {
		SkipApi bool `json:"skipApi"`
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
	EchoConfig struct {
		Port                 int                                                    `json:"port"`
		SecureCookies        *fluffycore_echo_contracts_cookies.SecureCookiesConfig `json:"secureCookies"`
		DisableSecureCookies bool                                                   `json:"disableSecureCookies"`
	}
	PasswordConfig struct {
		MinEntropyBits float64 `json:"minEntropyBits"`
	}
	InMemoryCacheConfig struct {
		DefaultExpirationSeconds int `json:"defaultExpirationSeconds"`
		CleanupIntervalSeconds   int `json:"cleanupIntervalSeconds"`
	}
	BackingCacheConfig struct {
		Type          string              `json:"type"`
		InMemoryCache InMemoryCacheConfig `json:"inMemoryCache"`
	}

	SelfIDPConfig struct {
		ClientID     string   `json:"clientId"`
		ClientSecret string   `json:"clientSecret"`
		RedirectURL  string   `json:"redirectUrl"`
		Authority    string   `json:"authority"`
		Scopes       []string `json:"scopes"`
	}

	OIDCConfig struct {
		BaseUrl            string `json:"baseUrl"`
		OAuth2CallbackPath string `json:"oauth2CallbackPath"`
	}
	CookieConfig struct {
		Domain string `json:"domain"`
	}
	SystemConfig struct {
		DeveloperMode bool   `json:"developerMode"`
		Domain        string `json:"domain"`
	}
	InitialConfig struct {
		ConfigFiles ConfigFiles `json:"configFiles"`
	}
	TOTPConfig struct {
		Enabled    bool   `json:"enabled"`
		IssuerName string `json:"issuerName"`
	}
	KeyValuePair struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	CacheBustingHTMLConfig struct {
		FilePath      string          `json:"filePath"`
		EchoPath      string          `json:"echoPath"`
		StaticPath    string          `json:"staticPath"`
		RootPath      string          `json:"rootPath"`
		ReplaceParams []*KeyValuePair `json:"replaceParams"`
	}
	OIDCUIConfig struct {
		AppSettings        *models_api_appsettings.OIDCUIAppSettings `json:"appSettings"`
		StaticFilePath     string                                    `json:"staticFilePath"`
		CacheBustingConfig *CacheBustingHTMLConfig                   `json:"cacheBustingConfig"`
		URIEntryPath       string                                    `json:"uriEntryPath"`
	}
	AccountUIConfig struct {
		AppSettings        *models_api_appsettings.AccountAppSettings `json:"appSettings"`
		StaticFilePath     string                                     `json:"staticFilePath"`
		CacheBustingConfig *CacheBustingHTMLConfig                    `json:"cacheBustingConfig"`
		URIEntryPath       string                                     `json:"uriEntryPath"`
	}
	Config struct {
		fluffycore_contracts_config.CoreConfig `mapstructure:",squash"`

		ConfigFiles                    ConfigFiles                             `json:"configFiles"`
		Echo                           *EchoConfig                             `json:"echo"`
		EchoOIDCUI                     *EchoConfig                             `json:"echoOIDCUI"`
		EchoAccount                    *EchoConfig                             `json:"echoAccount"`
		InMemoryClients                InMemoryClients                         `json:"inMemoryClients"`
		OIDCConfig                     *OIDCConfig                             `json:"oidcConfig"`
		BackingCache                   *BackingCacheConfig                     `json:"backingCache"`
		AutolinkOnEmailMatch           bool                                    `json:"autolinkOnEmailMatch"`
		EmailVerificationRequired      bool                                    `json:"emailVerificationRequired"`
		MultiFactorRequired            bool                                    `json:"multiFactorRequired"`
		MultiFactorRequiredByEmailCode bool                                    `json:"multiFactorRequiredByEmailCode"`
		DisableLocalAccountCreation    bool                                    `json:"disableLocalAccountCreation"`
		DisableSocialAccounts          bool                                    `json:"disableSocialAccounts"`
		TOTP                           *TOTPConfig                             `json:"totp"`
		EmailConfig                    *contracts_email.EmailConfig            `json:"emailConfig"`
		SelfIDPConfig                  *SelfIDPConfig                          `json:"selfIDPConfig"`
		CookieConfig                   *CookieConfig                           `json:"cookieConfig"`
		SystemConfig                   *SystemConfig                           `json:"systemConfig"`
		SessionConfig                  *contracts_sessions.SessionConfig       `json:"sessionConfig"`
		WebAuthNConfig                 *contracts_webauthn.WebAuthNConfig      `json:"webAuthNConfig"`
		PasswordConfig                 *PasswordConfig                         `json:"passwordConfig"`
		CORSConfig                     *CORSConfig                             `json:"corsConfig"`
		CSRFConfig                     *CSRFConfig                             `json:"csrfConfig"`
		OTELConfig                     *fluffycore_contracts_otel.OTELConfig   `json:"otelConfig"`
		DDConfig                       *fluffycore_contracts_ddprofiler.Config `json:"ddConfig"`

		OIDCUIConfig       *OIDCUIConfig                              `json:"oidcUIConfig"`
		AccountUIConfig    *AccountUIConfig                           `json:"accountUIConfig"`
		AccountAppSettings *models_api_appsettings.AccountAppSettings `json:"accountAppSettings"`
		ApiAppSettings     *models_api_appsettings.ApiAppSettings     `json:"apiAppSettings"`
	}
)

const (
	BackingCacheTypeInMemory = "in-memory"
	BackingCacheTypeRedis    = "redis"
)

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

	"apiAppSettings": {
	  "ApplicationEnvironment": "IN_ENVIRONMENT",
      "BaseApiUrl": "",
	  "PrivacyPolicyUrl": "",
	  "CookiePolicyUrl": ""
	},
    "accountUIConfig": {
    	"appSettings": {
			"applicationEnvironment": "IN_ENVIRONMENT",
			"baseApiUrl": ""
		},
		"staticFilePath": "./static",
		"uriEntryPath": "/",
		"cacheBustingConfig": {
			"filePath": "IN_ENVIRONMENT",
            "staticPath": "IN_ENVIRONMENT",
            "rootPath": "/",
			"echoPath": "/*"
		}
    },
	"oidcUIConfig": {
		"appSettings": {
			"applicationEnvironment": "IN_ENVIRONMENT",
			"baseApiUrl": ""
		},
		"staticFilePath": "./static",
		"uriEntryPath": "/oidc-login-ui/",
		"cacheBustingConfig": {
			"filePath": "IN_ENVIRONMENT",
            "staticPath": "IN_ENVIRONMENT",
            "rootPath": "/oidc-login-ui/",
			"echoPath": "/oidc-login-ui/*"
		}
	},
    "oidcUIAppSettings": {
        "applicationEnvironment": "IN_ENVIRONMENT",
        "baseApiUrl": ""
    },
    "accountAppSettings": {
        "applicationEnvironment": "IN_ENVIRONMENT",
        "baseApiUrl": ""
    },
    "cookieConfig": {
        "domain": "@@CHANGEME@@"
    },
    "csrfConfig": {
        "skipApi": false
    },
    "corsConfig": {
        "enabled": true,
        "allowedOrigins": [
            "*"
        ],
        "allowedMethods": [
            "GET",
            "POST",
            "PUT",
            "DELETE"
        ],
        "allowedHeaders": [
            "Authorization",
            "Content-Type",
            "x-csrf-token"
        ],
        "allowCredentials": true,
        "unsafeWildcardOriginWithAllowCredentials": false,
        "exposeHeaders": [],
        "maxAge": 0
    },
    "jwtValidators": {},
    "autolinkOnEmailMatch": true,
    "emailVerificationRequired": true,
    "multiFactorRequired": false,
    "multiFactorRequiredByEmailCode": false,
    "selfIDPConfig": {
        "clientId": "self-client",
        "clientSecret": "secret",
        "redirectUrl": "http://localhost:9044/auth/callback",
        "authority": "http://localhost:9044",
        "scopes": [
            "openid",
            "profile",
            "email"
        ]
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
    "echoOIDCUI": {
        "port": 0,
        "disableSecureCookies": false,
        "secureCookies": {
            "hashKey": "7f6a8b9c0d1e2f3a4b5c6d7e8f9a0b1c",
            "blockKey": "1234567890abcdef1234567890abcdef"
        }
    },
    "echoAccount": {
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
        "domain": "@@CHANGEME@@",
        "insecure": false
    },
    "systemConfig": {
        "domain": "@@CHANGEME@@",
        "developerMode": false
    },
    "totp": {
        "enabled": false,
        "issuerName": "RAGE.IDENTITY"
    },
    "webAuthNConfig": {
        "enabled": false,
        "rpDisplayName": "RAGE",
        "rpID": "[the domain]",
        "rpOrigins": []
    },
    "passwordConfig": {
        "minEntropyBits": 60
    },
	"ddConfig": {
        "ddProfilerConfig": {
            "enabled": false
        },
        "tracingEnabled": false,
        "serviceName": "in-environment",
        "applicationEnvironment": "in-environment",
        "version": "1.0.0"
    },
    "otelConfig": {
        "serviceName": "in-environment",
        "tracingConfig": {
            "enabled": false,
            "endpointType": "stdout",
            "endpoint": "localhost:4318"
        },
        "metricConfig": {
            "enabled": false,
            "endpointType": "stdout",
            "intervalSeconds": 10,
            "endpoint": "localhost:4318",
            "runtimeEnabled": false,
            "hostEnabled": false
        }
    }
}
`

/*
	Minimum length of 8 characters
	At least 2 uppercase letters
	At least 1 special character (such as !, @, #, $, &, )
	At least 2 digits
	At least 3 lowercase letters

	pattern := `^(?=.*[A-Z].*[A-Z])(?=.*[!@#$&*])(?=.*[0-9].*[0-9])(?=.*[a-z].*[a-z].*[a-z]).{8}$`

*/

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
