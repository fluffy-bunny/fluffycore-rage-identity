package inmemory

import (
	"context"
	"testing"

	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func newTestClientService(t *testing.T, seed *proto_oidc_models.Clients) proto_oidc_client.IFluffyCoreClientServiceServer {
	t.Helper()
	if seed == nil {
		seed = &proto_oidc_models.Clients{}
	}
	svc, err := stemService.Ctor(seed, nil)
	require.NoError(t, err)
	return svc
}

// CreateClient/UpdateClient/DeleteClient were unimplemented stubs (falling
// through to UnimplementedClientServiceServer) before this pass — these
// tests are the first coverage any of the three have had.
func TestCreateClientRoundTrip(t *testing.T) {
	svc := newTestClientService(t, nil)
	ctx := context.Background()

	createResp, err := svc.CreateClient(ctx, &proto_oidc_client.CreateClientRequest{
		Client: &proto_oidc_models.Client{ClientId: "acme", Enabled: true, ClientName: "Acme"},
	})
	require.NoError(t, err)
	require.Equal(t, "acme", createResp.Client.ClientId)

	getResp, err := svc.GetClient(ctx, &proto_oidc_client.GetClientRequest{ClientId: "acme"})
	require.NoError(t, err)
	require.Equal(t, "Acme", getResp.Client.ClientName)

	listResp, err := svc.ListClient(ctx, &proto_oidc_client.ListClientRequest{})
	require.NoError(t, err)
	require.Len(t, listResp.Clients, 1)
}

func TestCreateClientRejectsDuplicateClientId(t *testing.T) {
	svc := newTestClientService(t, nil)
	ctx := context.Background()

	_, err := svc.CreateClient(ctx, &proto_oidc_client.CreateClientRequest{
		Client: &proto_oidc_models.Client{ClientId: "acme"},
	})
	require.NoError(t, err)

	_, err = svc.CreateClient(ctx, &proto_oidc_client.CreateClientRequest{
		Client: &proto_oidc_models.Client{ClientId: "acme"},
	})
	require.Error(t, err)
	require.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestCreateClientRequiresClientId(t *testing.T) {
	svc := newTestClientService(t, nil)
	_, err := svc.CreateClient(context.Background(), &proto_oidc_client.CreateClientRequest{
		Client: &proto_oidc_models.Client{},
	})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestUpdateClientOnlyTouchesSetFields(t *testing.T) {
	svc := newTestClientService(t, &proto_oidc_models.Clients{Clients: []*proto_oidc_models.Client{
		{ClientId: "acme", ClientName: "Acme", Description: "original description", Enabled: true},
	}})
	ctx := context.Background()

	updateResp, err := svc.UpdateClient(ctx, &proto_oidc_client.UpdateClientRequest{
		Client: &proto_oidc_models.ClientUpdate{
			ClientId:   "acme",
			ClientName: &wrapperspb.StringValue{Value: "Acme Renamed"},
			// Description and Enabled deliberately left nil: must survive untouched.
		},
	})
	require.NoError(t, err)
	require.Equal(t, "Acme Renamed", updateResp.Client.ClientName)
	require.Equal(t, "original description", updateResp.Client.Description, "unset ClientUpdate fields must not overwrite existing values")
	require.True(t, updateResp.Client.Enabled, "unset ClientUpdate fields must not overwrite existing values")
}

func TestUpdateClientReplaceStringArray(t *testing.T) {
	svc := newTestClientService(t, &proto_oidc_models.Clients{Clients: []*proto_oidc_models.Client{
		{ClientId: "acme", AllowedScopes: []string{"openid", "profile"}},
	}})
	ctx := context.Background()

	updateResp, err := svc.UpdateClient(ctx, &proto_oidc_client.UpdateClientRequest{
		Client: &proto_oidc_models.ClientUpdate{
			ClientId: "acme",
			AllowedScopes: &proto_types.StringArrayUpdate{
				Update: &proto_types.StringArrayUpdate_Replace{
					Replace: &proto_types.StringValues{Values: []string{"openid", "email"}},
				},
			},
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"openid", "email"}, updateResp.Client.AllowedScopes)
}

func TestUpdateClientGranularStringArray(t *testing.T) {
	svc := newTestClientService(t, &proto_oidc_models.Clients{Clients: []*proto_oidc_models.Client{
		{ClientId: "acme", AllowedGrantTypes: []string{"authorization_code", "refresh_token"}},
	}})
	ctx := context.Background()

	updateResp, err := svc.UpdateClient(ctx, &proto_oidc_client.UpdateClientRequest{
		Client: &proto_oidc_models.ClientUpdate{
			ClientId: "acme",
			AllowedGrantTypes: &proto_types.StringArrayUpdate{
				Update: &proto_types.StringArrayUpdate_Granular_{
					Granular: &proto_types.StringArrayUpdate_Granular{
						Add:    []string{"client_credentials"},
						Remove: []string{"refresh_token"},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"authorization_code", "client_credentials"}, updateResp.Client.AllowedGrantTypes)
}

func TestUpdateClientUnknownClientIdFails(t *testing.T) {
	svc := newTestClientService(t, nil)
	_, err := svc.UpdateClient(context.Background(), &proto_oidc_client.UpdateClientRequest{
		Client: &proto_oidc_models.ClientUpdate{ClientId: "does-not-exist"},
	})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeleteClientRemovesFromListAndMap(t *testing.T) {
	svc := newTestClientService(t, &proto_oidc_models.Clients{Clients: []*proto_oidc_models.Client{
		{ClientId: "acme"},
		{ClientId: "other"},
	}})
	ctx := context.Background()

	_, err := svc.DeleteClient(ctx, &proto_oidc_client.DeleteClientRequest{ClientId: "acme"})
	require.NoError(t, err)

	_, err = svc.GetClient(ctx, &proto_oidc_client.GetClientRequest{ClientId: "acme"})
	require.Error(t, err)

	listResp, err := svc.ListClient(ctx, &proto_oidc_client.ListClientRequest{})
	require.NoError(t, err)
	require.Len(t, listResp.Clients, 1)
	require.Equal(t, "other", listResp.Clients[0].ClientId)

	// Deleting an already-deleted client must fail, not silently no-op.
	_, err = svc.DeleteClient(ctx, &proto_oidc_client.DeleteClientRequest{ClientId: "acme"})
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}
