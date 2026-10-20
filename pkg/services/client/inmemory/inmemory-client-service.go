package inmemory

import (
	"context"
	"sort"
	"sync"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		proto_oidc_client.UnimplementedClientServiceServer

		mu             sync.RWMutex
		clientsMap     map[string]*proto_oidc_models.Client
		clients        *proto_oidc_models.Clients
		passwordHasher contracts_identity.IPasswordHasher
	}
)

var stemService = (*service)(nil)
var _ proto_oidc_client.IFluffyCoreClientServiceServer = stemService

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
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &proto_oidc_client.ListClientResponse{
		Clients: s.clients.Clients,
	}, nil
}

func (s *service) validateCreateClientRequest(request *proto_oidc_client.CreateClientRequest) error {
	if request == nil || request.Client == nil {
		return status.Error(codes.InvalidArgument, "client is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Client.ClientId) {
		return status.Error(codes.InvalidArgument, "client.client_id is required")
	}
	return nil
}

// Create client
func (s *service) CreateClient(ctx context.Context, request *proto_oidc_client.CreateClientRequest) (*proto_oidc_client.CreateClientResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	if err := s.validateCreateClientRequest(request); err != nil {
		log.Warn().Err(err).Msg("validateCreateClientRequest")
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.clientsMap[request.Client.ClientId]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "client %s already exists", request.Client.ClientId)
	}
	s.clientsMap[request.Client.ClientId] = request.Client
	s.clients.Clients = append(s.clients.Clients, request.Client)
	return &proto_oidc_client.CreateClientResponse{Client: request.Client}, nil
}

// Delete client
func (s *service) DeleteClient(ctx context.Context, request *proto_oidc_client.DeleteClientRequest) (*proto_oidc_client.DeleteClientResponse, error) {
	if request == nil || fluffycore_utils.IsEmptyOrNil(request.ClientId) {
		return nil, status.Error(codes.InvalidArgument, "client_id is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clientsMap[request.ClientId]; !ok {
		return nil, status.Errorf(codes.NotFound, "client %s not found", request.ClientId)
	}
	delete(s.clientsMap, request.ClientId)
	filtered := make([]*proto_oidc_models.Client, 0, len(s.clients.Clients))
	for _, existing := range s.clients.Clients {
		if existing.ClientId != request.ClientId {
			filtered = append(filtered, existing)
		}
	}
	s.clients.Clients = filtered
	return &proto_oidc_client.DeleteClientResponse{}, nil
}

func (s *service) validateUpdateClientRequest(request *proto_oidc_client.UpdateClientRequest) error {
	if request == nil || request.Client == nil {
		return status.Error(codes.InvalidArgument, "client is required")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Client.ClientId) {
		return status.Error(codes.InvalidArgument, "client.client_id is required")
	}
	return nil
}

// Update client. Every scalar/array field on ClientUpdate is a wrapper type
// (google.protobuf.BoolValue, StringArrayUpdate, ...) that's nil when the
// caller didn't intend to change it — a partial update, applied in place
// onto the existing record. Claims and ClientSecrets updates are not yet
// supported by this in-memory store (see docs/OPERATOR_PORTAL_PLAN.md); a
// request that only sets those fields is a silent no-op today.
func (s *service) UpdateClient(ctx context.Context, request *proto_oidc_client.UpdateClientRequest) (*proto_oidc_client.UpdateClientResponse, error) {
	log := zerolog.Ctx(ctx).With().Logger()
	if err := s.validateUpdateClientRequest(request); err != nil {
		log.Warn().Err(err).Msg("validateUpdateClientRequest")
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.clientsMap[request.Client.ClientId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "client %s not found", request.Client.ClientId)
	}
	applyClientUpdate(existing, request.Client)
	return &proto_oidc_client.UpdateClientResponse{Client: existing}, nil
}

// applyClientUpdate merges the non-nil fields of update onto existing.
func applyClientUpdate(existing *proto_oidc_models.Client, update *proto_oidc_models.ClientUpdate) {
	if update.Enabled != nil {
		existing.Enabled = update.Enabled.Value
	}
	if update.ClientName != nil {
		existing.ClientName = update.ClientName.Value
	}
	if update.Description != nil {
		existing.Description = update.Description.Value
	}
	if update.Namespace != nil {
		existing.Namespace = update.Namespace.Value
	}
	if update.RequireRefreshClientSecret != nil {
		existing.RequireRefreshClientSecret = update.RequireRefreshClientSecret.Value
	}
	if update.AllowOfflineAccess != nil {
		existing.AllowOfflineAccess = update.AllowOfflineAccess.Value
	}
	if update.UpdateAccessTokenClaimsOnRefresh != nil {
		existing.UpdateAccessTokenClaimsOnRefresh = update.UpdateAccessTokenClaimsOnRefresh.Value
	}
	if update.IncludeJwtId != nil {
		existing.IncludeJwtId = update.IncludeJwtId.Value
	}
	if update.AlwaysSendClientClaims != nil {
		existing.AlwaysSendClientClaims = update.AlwaysSendClientClaims.Value
	}
	if update.AlwaysIncludeUserClaimsInIdToken != nil {
		existing.AlwaysIncludeUserClaimsInIdToken = update.AlwaysIncludeUserClaimsInIdToken.Value
	}
	if update.RequireClientSecret != nil {
		existing.RequireClientSecret = update.RequireClientSecret.Value
	}
	if update.AccessTokenLifetime != nil {
		existing.AccessTokenLifetime = update.AccessTokenLifetime.Value
	}
	if update.AbsoluteRefreshTokenLifetime != nil {
		existing.AbsoluteRefreshTokenLifetime = update.AbsoluteRefreshTokenLifetime.Value
	}
	if update.SlidingRefreshTokenLifetime != nil {
		existing.SlidingRefreshTokenLifetime = update.SlidingRefreshTokenLifetime.Value
	}
	if update.IdentityTokenLifetime != nil {
		existing.IdentityTokenLifetime = update.IdentityTokenLifetime.Value
	}
	if update.RefreshTokenGraceEnabled != nil {
		existing.RefreshTokenGraceEnabled = update.RefreshTokenGraceEnabled.Value
	}
	if update.RefreshTokenGraceTtl != nil {
		existing.RefreshTokenGraceTtl = update.RefreshTokenGraceTtl.Value
	}
	if update.RefreshTokenGraceMaxAttempts != nil {
		existing.RefreshTokenGraceMaxAttempts = update.RefreshTokenGraceMaxAttempts.Value
	}
	if update.RefreshTokenUsage != nil {
		existing.RefreshTokenUsage = update.RefreshTokenUsage.Value
	}
	if update.RefreshTokenExpiration != nil {
		existing.RefreshTokenExpiration = update.RefreshTokenExpiration.Value
	}
	if update.AccessTokenType != nil {
		existing.AccessTokenType = update.AccessTokenType.Value
	}
	applyStringArrayUpdate(&existing.AllowedScopes, update.AllowedScopes)
	applyStringArrayUpdate(&existing.AllowedGrantTypes, update.AllowedGrantTypes)
	applyStringArrayUpdate(&existing.AllowedPermissions, update.AllowedPermissions)
	applyStringArrayUpdate(&existing.AllowedDynamicClaimTypes, update.AllowedDynamicClaimTypes)
	applyStringArrayUpdate(&existing.AllowedTokenExchangeRequestedTokenTypes, update.AllowedTokenExchangeRequestedTokenTypes)
	applyStringArrayUpdate(&existing.AllowedRedirectUris, update.AllowedRedirectUris)
}

// applyStringArrayUpdate mutates *target in place per update's oneof: nil
// means "no change", Replace overwrites wholesale, Granular adds/removes.
func applyStringArrayUpdate(target *[]string, update *proto_types.StringArrayUpdate) {
	if update == nil {
		return
	}
	switch v := update.Update.(type) {
	case *proto_types.StringArrayUpdate_Replace:
		*target = append([]string{}, v.Replace.GetValues()...)
	case *proto_types.StringArrayUpdate_Granular_:
		set := make(map[string]bool, len(*target))
		for _, existing := range *target {
			set[existing] = true
		}
		for _, removed := range v.Granular.GetRemove() {
			delete(set, removed)
		}
		for _, added := range v.Granular.GetAdd() {
			set[added] = true
		}
		result := make([]string, 0, len(set))
		for s := range set {
			result = append(result, s)
		}
		sort.Strings(result)
		*target = result
	}
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
