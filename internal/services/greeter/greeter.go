package greeter

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	proto_helloworld "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/helloworld"
	endpoint "github.com/fluffy-bunny/fluffycore/contracts/endpoint"
	grpc_gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	zerolog "github.com/rs/zerolog"
	grpc "google.golang.org/grpc"
)

type (
	service struct {
		proto_helloworld.GreeterFluffyCoreServer

		config *contracts_config.Config
	}
)

var (
	stemService = (*service)(nil)
)

func init() {
	var _ proto_helloworld.IFluffyCoreGreeterServer = (*service)(nil)
	var _ endpoint.IEndpointRegistration = (*service)(nil)
}

func (s *service) Ctor(
	config *contracts_config.Config,
) proto_helloworld.IFluffyCoreGreeterServer {
	return &service{
		config: config,
	}
}
func (s *service) RegisterFluffyCoreHandler(gwmux *grpc_gateway_runtime.ServeMux, conn *grpc.ClientConn) {
	proto_helloworld.RegisterGreeterHandler(context.Background(), gwmux, conn)
}

func AddGreeterService(builder di.ContainerBuilder) {
	proto_helloworld.AddGreeterServerWithExternalRegistration(builder,
		stemService.Ctor,
		func() endpoint.IEndpointRegistration {
			return &service{}
		})
}

func (s *service) SayHello(ctx context.Context, request *proto_helloworld.HelloRequest) (*proto_helloworld.HelloReply, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("SayHello")
	return &proto_helloworld.HelloReply{
		Message: "Hello " + request.Name,
	}, nil
}
