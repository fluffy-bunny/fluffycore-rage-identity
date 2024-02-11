package inmemory

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/models"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/oidc/user"
)

type (
	service struct {
		proto_oidc_client.UnimplementedUserServiceServer
		userMap map[string]*proto_oidc_models.User
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_client.IFluffyCoreUserServiceServer = stemService
}
func (s *service) Ctor(clients *proto_oidc_models.Clients) (proto_oidc_client.IFluffyCoreUserServiceServer, error) {
	return &service{
		userMap: make(map[string]*proto_oidc_models.User),
	}, nil
}

func AddSingletonIFluffyCoreUserServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_client.IFluffyCoreUserServiceServer](cb, stemService.Ctor)
}
