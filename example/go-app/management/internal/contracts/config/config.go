package config

import (
	"context"

	management_contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/config"
)

type (
	AuthorizationStateCookie struct {
		State string `json:"state"`
	}
	IAppConfigAccessor interface {
		GetAppConfig(ctx context.Context) *management_contracts_config.AppConfig
	}
)
