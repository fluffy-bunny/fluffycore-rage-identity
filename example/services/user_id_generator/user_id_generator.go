package user_id_generator

import (
	"fmt"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	xid "github.com/rs/xid"
)

type (
	service struct{}
)

var stemService = (*service)(nil)
var _ contracts_identity.IUserIdGenerator = stemService

func (s *service) Ctor() (contracts_identity.IUserIdGenerator, error) {
	return &service{}, nil
}

func AddSingletonIUserIdGenerator(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_identity.IUserIdGenerator](cb, stemService.Ctor)
}

func (s *service) GenerateUserId() string {
	return fmt.Sprintf("rage_%s", xid.New().String())
}
