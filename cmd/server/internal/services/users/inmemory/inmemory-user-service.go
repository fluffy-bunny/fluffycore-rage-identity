package inmemory

import (
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
)

type (
	service struct {
		proto_external_user.UnimplementedUserServiceServer

		userMap map[string]*proto_external_models.User
		rwLock  sync.RWMutex
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_external_user.IFluffyCoreUserServiceServer = stemService
}
func (s *service) Ctor() (proto_external_user.IFluffyCoreUserServiceServer, error) {
	return &service{
		userMap: make(map[string]*proto_external_models.User),
	}, nil
}

func AddSingletonIFluffyCoreUserServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_external_user.IFluffyCoreUserServiceServer](cb, stemService.Ctor)
}
