package types

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
)

type (
	ConfigureServices func(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder)
)
