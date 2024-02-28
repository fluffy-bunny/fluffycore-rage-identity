package inmemory

import (
	"context"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		proto_oidc_client.UnimplementedClientServiceServer
		clientsMap     map[string]*proto_oidc_models.Client
		clients        *proto_oidc_models.Clients
		passwordHasher contracts_identity.IPasswordHasher
	}
)

var stemService = (*service)(nil)

func init() {
	var _ proto_oidc_client.IFluffyCoreClientServiceServer = stemService
}
func (s *service) Ctor(clients *proto_oidc_models.Clients,
	passwordHasher contracts_identity.IPasswordHasher) (proto_oidc_client.IFluffyCoreClientServiceServer, error) {
	clientsMap := make(map[string]*proto_oidc_models.Client)
	for _, client := range clients.Clients {
		clientsMap[client.ClientId] = client
	}
	return &service{
		clientsMap:     clientsMap,
		clients:        clients,
		passwordHasher: passwordHasher,
	}, nil
}

func AddSingletonIFluffyCoreClientServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_client.IFluffyCoreClientServiceServer](cb, stemService.Ctor)
}

func (s *service) validateGetClientRequest(request *proto_oidc_client.GetClientRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ClientId) {
		return status.Error(codes.InvalidArgument, "clientId is required")
	}
	return nil
}

// Get client
func (s *service) GetClient(ctx context.Context, request *proto_oidc_client.GetClientRequest) (*proto_oidc_client.GetClientResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateGetClientRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateGetClientRequest")
		return nil, err
	}
	client, ok := s.clientsMap[request.ClientId]
	if ok {
		return &proto_oidc_client.GetClientResponse{
			Client: client,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "Client not found")
}

// List clients
func (s *service) ListClient(ctx context.Context, request *proto_oidc_client.ListClientRequest) (*proto_oidc_client.ListClientResponse, error) {
	// for now return everything.  There is no reason
	return &proto_oidc_client.ListClientResponse{
		Clients: s.clients.Clients,
	}, nil
}
func (s *service) validateValidateClientSecretRequest(request *proto_oidc_client.ValidateClientSecretRequest) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.ClientId) {
		return status.Error(codes.InvalidArgument, "clientId is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Secret) {
		return status.Error(codes.InvalidArgument, "secret is required")
	}
	return nil
}

// Generate a new client secret
func (s *service) ValidateClientSecret(ctx context.Context, request *proto_oidc_client.ValidateClientSecretRequest) (*proto_oidc_client.ValidateClientSecretResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	err := s.validateValidateClientSecretRequest(request)
	if err != nil {
		log.Warn().Err(err).Msg("validateValidateClientSecretRequest")
		return nil, err
	}
	now := time.Now()
	skewLow := now.Add(-1 * time.Duration(time.Minute*5))
	isExpired := func(expiration time.Time) bool {
		return expiration.Before(skewLow)
	}
	getClientResponse, err := s.GetClient(ctx, &proto_oidc_client.GetClientRequest{
		ClientId: request.ClientId,
	})
	if err != nil {
		log.Warn().Err(err).Msg("GetClient")
		return nil, err
	}

	for _, clientSecrets := range getClientResponse.Client.GetClientSecrets() {
		if isExpired(time.Unix(clientSecrets.ExpirationUnix, 0)) {
			continue
		}
		err = s.passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
			HashedPassword: clientSecrets.Hash,
			Password:       request.Secret,
		})

		if err == nil {
			return &proto_oidc_client.ValidateClientSecretResponse{
				Valid: true,
			}, nil
		}
	}
	log.Warn().Msg("Client secret not found")
	return &proto_oidc_client.ValidateClientSecretResponse{
		Valid: false,
	}, nil
}
