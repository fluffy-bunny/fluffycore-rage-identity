package runtime

import (
	internal_runtime "github.com/fluffy-bunny/fluffycore-rage-identity/internal/runtime"
	internal_types "github.com/fluffy-bunny/fluffycore-rage-identity/internal/types"
	public_types "github.com/fluffy-bunny/fluffycore-rage-identity/public/types"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
)

func NewStartup(options ...public_types.WithOption) fluffycore_contracts_runtime.IStartup {
	internalOptions := make([]internal_runtime.WithOption, len(options))
	for i, option := range options {
		internalOptions[i] = internal_runtime.WithOption(option)
	}
	return internal_runtime.NewStartup(internalOptions...)
}
func WithConfigureServices(extConfigureServices public_types.ConfigureServices) public_types.WithOption {
	intOpt := internal_runtime.WithConfigureServices(internal_types.ConfigureServices(extConfigureServices))
	return public_types.WithOption(intOpt)
}
