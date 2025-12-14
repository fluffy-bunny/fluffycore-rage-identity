package config

import (
	"context"
	"encoding/json"
	"os"

	oidc_login_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/contracts/config"
	contracts_OIDCFlowAppConfig "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/OIDCFlowAppConfig"
)

type (
	IDP struct {
		Slug string `json:"slug"`
	}

	AuthorizationStateCookie struct {
		State string `json:"state"`
	}
	IAppConfigAccessor interface {
		GetAppConfig(ctx context.Context) *oidc_login_contracts_config.AppConfig
		GetOIDCFlowAppConfig(ctx context.Context) (*contracts_OIDCFlowAppConfig.OIDCFlowAppConfig, error)
		GetAuthorizationStateCookie(ctx context.Context) (*AuthorizationStateCookie, error)
	}
)

const DefaultConfigJSON = `{
	"version": "1.0.0",
	"socialIDPs": [],
	"developmentMode": false,
	"disableLocalAccountCreation": false,
	"disableSocialAccounts": false,
	"bannerBranding": {
		"title": "Rage Identity",
		"logoUrl": "/web/apple-touch-icon-192x192.png",
		"showBannerVersion": true
	}
}`

func LoadConfigFromFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg T
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
