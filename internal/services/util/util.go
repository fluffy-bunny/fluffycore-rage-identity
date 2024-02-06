package somedisposable

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_util "github.com/fluffy-bunny/fluffycore-hanko-oidc/internal/contracts/util"
)

type (
	service struct{}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_util.ISomeUtil = stemService
}
func (s *service) Ctor() (contracts_util.ISomeUtil, error) {
	return &service{}, nil
}

func AddSingletonISomeUtil(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_util.ISomeUtil](cb, stemService.Ctor)
}

func (s *service) DoSomething(ctx context.Context) (string, error) {
	return "Hello World", nil
}
