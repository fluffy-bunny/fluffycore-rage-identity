package types

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	echo "github.com/labstack/echo/v4"
)

type (
	ConfigureServices             func(ctx context.Context, config *contracts_config.Config, builder di.ContainerBuilder)
	ConfigureManagementMiddleware func(ctn di.Container) echo.MiddlewareFunc
)
