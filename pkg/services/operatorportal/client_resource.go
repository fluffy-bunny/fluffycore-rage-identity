package operatorportal

import (
	"context"

	resource "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/operatorportal/resource"
	proto_oidc_client "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/client"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// The Client message has ~30 fields (token lifetimes, claims, secrets, ...);
// this resource intentionally manages a working subset — the fields an
// operator actually edits day to day. Extending the subset later is adding
// entries to clientRowFromModel/Fields/applyClientRow, not a new mechanism.

func clientRowFromModel(client *proto_oidc_models.Client) resource.Row {
	return resource.Row{
		"client_id":             client.GetClientId(),
		"client_name":           client.GetClientName(),
		"description":           client.GetDescription(),
		"namespace":             client.GetNamespace(),
		"enabled":               resource.FormatBool(client.GetEnabled()),
		"require_client_secret": resource.FormatBool(client.GetRequireClientSecret()),
		"allow_offline_access":  resource.FormatBool(client.GetAllowOfflineAccess()),
		"allowed_scopes":        resource.JoinList(client.GetAllowedScopes()),
		"allowed_grant_types":   resource.JoinList(client.GetAllowedGrantTypes()),
		"allowed_redirect_uris": resource.JoinList(client.GetAllowedRedirectUris()),
	}
}

func clientModelFromRow(row resource.Row) *proto_oidc_models.Client {
	return &proto_oidc_models.Client{
		ClientId:            row["client_id"],
		ClientName:          row["client_name"],
		Description:         row["description"],
		Namespace:           row["namespace"],
		Enabled:             resource.ParseBool(row["enabled"]),
		RequireClientSecret: resource.ParseBool(row["require_client_secret"]),
		AllowOfflineAccess:  resource.ParseBool(row["allow_offline_access"]),
		AllowedScopes:       resource.SplitList(row["allowed_scopes"]),
		AllowedGrantTypes:   resource.SplitList(row["allowed_grant_types"]),
		AllowedRedirectUris: resource.SplitList(row["allowed_redirect_uris"]),
	}
}

// clientUpdateFromRow builds a partial ClientUpdate covering exactly the
// fields this resource manages — every other ClientUpdate field is left nil
// ("no change"), which is what makes editing only this subset safe even
// though Client has ~30 fields in total.
func clientUpdateFromRow(id string, row resource.Row) *proto_oidc_models.ClientUpdate {
	replace := func(values []string) *proto_types.StringArrayUpdate {
		return &proto_types.StringArrayUpdate{
			Update: &proto_types.StringArrayUpdate_Replace{
				Replace: &proto_types.StringValues{Values: values},
			},
		}
	}
	return &proto_oidc_models.ClientUpdate{
		ClientId:            id,
		ClientName:          &wrapperspb.StringValue{Value: row["client_name"]},
		Description:         &wrapperspb.StringValue{Value: row["description"]},
		Namespace:           &wrapperspb.StringValue{Value: row["namespace"]},
		Enabled:             &wrapperspb.BoolValue{Value: resource.ParseBool(row["enabled"])},
		RequireClientSecret: &wrapperspb.BoolValue{Value: resource.ParseBool(row["require_client_secret"])},
		AllowOfflineAccess:  &wrapperspb.BoolValue{Value: resource.ParseBool(row["allow_offline_access"])},
		AllowedScopes:       replace(resource.SplitList(row["allowed_scopes"])),
		AllowedGrantTypes:   replace(resource.SplitList(row["allowed_grant_types"])),
		AllowedRedirectUris: replace(resource.SplitList(row["allowed_redirect_uris"])),
	}
}

// newClientResourceSpec adapts proto_oidc_client.IFluffyCoreClientServiceServer
// to the generic resource scaffold. Create/Update/Delete were unimplemented
// stubs on the in-memory client service before this pass (only List/Get/
// ValidateClientSecret existed, despite the CRUD-shaped proto) — see
// pkg/services/client/inmemory/inmemory-client-service.go.
func newClientResourceSpec(clientService proto_oidc_client.IFluffyCoreClientServiceServer) *resource.Spec {
	return &resource.Spec{
		Key:          "clients",
		Title:        "OIDC Clients",
		SingularName: "Client",
		IDField:      "client_id",
		Columns: []resource.Field{
			{Key: "client_id", Label: "Client ID"},
			{Key: "client_name", Label: "Name"},
			{Key: "namespace", Label: "Namespace"},
			{Key: "enabled", Label: "Enabled", Kind: resource.Checkbox},
		},
		Fields: []resource.Field{
			{Key: "client_id", Label: "Client ID", Immutable: true,
				Help: "Unique identifier. Cannot be changed after creation."},
			{Key: "client_name", Label: "Name"},
			{Key: "description", Label: "Description", Kind: resource.TextArea},
			{Key: "namespace", Label: "Namespace"},
			{Key: "enabled", Label: "Enabled", Kind: resource.Checkbox},
			{Key: "require_client_secret", Label: "Require client secret", Kind: resource.Checkbox},
			{Key: "allow_offline_access", Label: "Allow offline access (refresh tokens)", Kind: resource.Checkbox},
			{Key: "allowed_scopes", Label: "Allowed Scopes", Kind: resource.StringList},
			{Key: "allowed_grant_types", Label: "Allowed Grant Types", Kind: resource.StringList},
			{Key: "allowed_redirect_uris", Label: "Allowed Redirect URIs", Kind: resource.StringList},
		},
		List: func(ctx context.Context) ([]resource.Row, error) {
			resp, err := clientService.ListClient(ctx, &proto_oidc_client.ListClientRequest{})
			if err != nil {
				return nil, err
			}
			rows := make([]resource.Row, 0, len(resp.Clients))
			for _, client := range resp.Clients {
				rows = append(rows, clientRowFromModel(client))
			}
			return rows, nil
		},
		Get: func(ctx context.Context, id string) (resource.Row, error) {
			resp, err := clientService.GetClient(ctx, &proto_oidc_client.GetClientRequest{ClientId: id})
			if err != nil {
				return nil, err
			}
			return clientRowFromModel(resp.Client), nil
		},
		Create: func(ctx context.Context, row resource.Row) (resource.Row, error) {
			resp, err := clientService.CreateClient(ctx, &proto_oidc_client.CreateClientRequest{
				Client: clientModelFromRow(row),
			})
			if err != nil {
				return nil, err
			}
			return clientRowFromModel(resp.Client), nil
		},
		Update: func(ctx context.Context, id string, row resource.Row) (resource.Row, error) {
			resp, err := clientService.UpdateClient(ctx, &proto_oidc_client.UpdateClientRequest{
				Client: clientUpdateFromRow(id, row),
			})
			if err != nil {
				return nil, err
			}
			return clientRowFromModel(resp.Client), nil
		},
		Delete: func(ctx context.Context, id string) error {
			_, err := clientService.DeleteClient(ctx, &proto_oidc_client.DeleteClientRequest{ClientId: id})
			return err
		},
	}
}
