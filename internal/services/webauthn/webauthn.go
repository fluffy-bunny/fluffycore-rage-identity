package greeter

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/contracts/config"
	proto_auth_webauthn "github.com/fluffy-bunny/fluffycore-rage-oidc/proto/auth/webauthn"
	endpoint "github.com/fluffy-bunny/fluffycore/contracts/endpoint"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	grpc_gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	zerolog "github.com/rs/zerolog"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		proto_auth_webauthn.WebAuthNServiceFluffyCoreServer

		config *contracts_config.Config
	}
)

var (
	stemService = (*service)(nil)
)

func init() {
	var _ proto_auth_webauthn.IFluffyCoreWebAuthNServiceServer = (*service)(nil)
	var _ endpoint.IEndpointRegistration = (*service)(nil)
}

func (s *service) Ctor(
	config *contracts_config.Config,
) proto_auth_webauthn.IFluffyCoreWebAuthNServiceServer {
	return &service{
		config: config,
	}
}
func (s *service) RegisterFluffyCoreHandler(gwmux *grpc_gateway_runtime.ServeMux, conn *grpc.ClientConn) {
	proto_auth_webauthn.RegisterWebAuthNServiceHandler(context.Background(), gwmux, conn)
}

func AddWebAuthNServiceServer(builder di.ContainerBuilder) {
	proto_auth_webauthn.AddWebAuthNServiceServerWithExternalRegistration(builder,
		stemService.Ctor,
		func() endpoint.IEndpointRegistration {
			return &service{}
		})
}
func (s *service) validateGetCredentialCreateOptionsRequest(request *proto_auth_webauthn.GetCredentialCreateOptionsRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Username) {
		return status.Error(codes.InvalidArgument, "Username is required")
	}
	return nil
}
func (s *service) GetCredentialCreateOptions(ctx context.Context, request *proto_auth_webauthn.GetCredentialCreateOptionsRequest) (*proto_auth_webauthn.GetCredentialCreateOptionsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Interface("request", request).Msg("GetCredentialCreateOptions")
	err := s.validateGetCredentialCreateOptionsRequest(request)
	if err != nil {
		return nil, err
	}
	return &proto_auth_webauthn.GetCredentialCreateOptionsResponse{}, nil
}
