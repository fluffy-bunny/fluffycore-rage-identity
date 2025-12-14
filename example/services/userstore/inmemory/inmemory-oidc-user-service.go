package inmemory

import (
	"reflect"
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_external_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/models"
	proto_external_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/external/user"
	proto_oidc_user "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/user"
)

type (
	service struct {
		proto_oidc_user.UnimplementedRageUserServiceServer
		proto_external_user.UnimplementedUserServiceServer

		userMap map[string]*proto_external_models.ExampleUser
		rwLock  sync.RWMutex
	}
)

var stemService = (*service)(nil)

var _ proto_oidc_user.IFluffyCoreRageUserServiceServer = stemService
var _ proto_external_user.IFluffyCoreUserServiceServer = stemService

func (s *service) Ctor() (*service, error) {
	return &service{
		userMap: make(map[string]*proto_external_models.ExampleUser),
	}, nil
}

func AddSingletonIFluffyCoreUserServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[*service](
		cb,
		stemService.Ctor,
		reflect.TypeOf((*proto_oidc_user.IFluffyCoreRageUserServiceServer)(nil)),
		reflect.TypeOf((*proto_external_user.IFluffyCoreUserServiceServer)(nil)),
	)
}
