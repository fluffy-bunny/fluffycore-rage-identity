package types

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	ConfigureServices func(ctx context.Context, builder di.ContainerBuilder)
)
