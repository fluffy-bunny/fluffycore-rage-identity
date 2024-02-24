package inmemory

import (
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
)

type (
	service struct {
		proto_oidc_user.UnimplementedUserServiceServer

		userMap map[string]*proto_oidc_models.User
		rwLock  sync.RWMutex
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_user.IFluffyCoreUserServiceServer = stemService
}
func (s *service) Ctor() (proto_oidc_user.IFluffyCoreUserServiceServer, error) {
	return &service{
		userMap: make(map[string]*proto_oidc_models.User),
	}, nil
}

func AddSingletonIFluffyCoreUserServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_user.IFluffyCoreUserServiceServer](cb, stemService.Ctor)
}
