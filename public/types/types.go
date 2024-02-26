package types

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	internal_runtime "github.com/fluffy-bunny/fluffycore-rage-identity/internal/runtime"
)

type (
	ConfigureServices func(ctx context.Context, builder di.ContainerBuilder)
	WithOption        internal_runtime.WithOption
)
