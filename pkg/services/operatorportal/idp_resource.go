package operatorportal

import (
	"context"

	resource "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/operatorportal/resource"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
)

// idpRowFromModel and idpModelFromRow are the resource.Row<->*IDP mapping —
// the only IDP-specific code in this file besides the field list itself.
// Only the OIDC branch of the IDP.Protocol oneof is exposed here; adding a
// form for the OAuth2/GitHub branches later means adding fields to Fields
// below and a case here, not touching pkg/services/operatorportal/resource.
func idpRowFromModel(idp *proto_oidc_models.IDP) resource.Row {
	oidc := idp.GetProtocol().GetOidc()
	return resource.Row{
		"slug":                        idp.GetSlug(),
		"name":                        idp.GetName(),
		"description":                 idp.GetDescription(),
		"enabled":                     resource.FormatBool(idp.GetEnabled()),
		"hidden":                      resource.FormatBool(idp.GetHidden()),
		"auto_create":                 resource.FormatBool(idp.GetAutoCreate()),
		"multi_factor_required":       resource.FormatBool(idp.GetMultiFactorRequired()),
		"email_verification_required": resource.FormatBool(idp.GetEmailVerificationRequired()),
		"claimed_domains":             resource.JoinList(idp.GetClaimedDomains()),
		"oidc_authority":              oidc.GetAuthority(),
		"oidc_client_id":              oidc.GetClientId(),
		"oidc_client_secret":          oidc.GetClientSecret(),
		"oidc_scope":                  oidc.GetScope(),
	}
}

func idpModelFromRow(row resource.Row) *proto_oidc_models.IDP {
	return &proto_oidc_models.IDP{
		Slug:                      row["slug"],
		Name:                      row["name"],
		Description:               row["description"],
		Enabled:                   resource.ParseBool(row["enabled"]),
		Hidden:                    resource.ParseBool(row["hidden"]),
		AutoCreate:                resource.ParseBool(row["auto_create"]),
		MultiFactorRequired:       resource.ParseBool(row["multi_factor_required"]),
		EmailVerificationRequired: resource.ParseBool(row["email_verification_required"]),
		ClaimedDomains:            resource.SplitList(row["claimed_domains"]),
		Protocol: &proto_oidc_models.Protocol{
			Value: &proto_oidc_models.Protocol_Oidc{
				Oidc: &proto_oidc_models.OIDCProtocol{
					Authority:    row["oidc_authority"],
					ClientId:     row["oidc_client_id"],
					ClientSecret: row["oidc_client_secret"],
					Scope:        row["oidc_scope"],
				},
			},
		},
	}
}

// newIDPResourceSpec adapts proto_oidc_idp.IFluffyCoreExampleIDPStoreServer
// (the sqlite-backed, full-CRUD IDP store — see pkg/services/idp/sqlite) to
// the generic resource scaffold.
func newIDPResourceSpec(idpStore proto_oidc_idp.IFluffyCoreExampleIDPStoreServer) *resource.Spec {
	return &resource.Spec{
		Key:          "idps",
		Title:        "Identity Providers",
		SingularName: "IDP",
		IDField:      "slug",
		Columns: []resource.Field{
			{Key: "slug", Label: "Slug"},
			{Key: "name", Label: "Name"},
			{Key: "enabled", Label: "Enabled", Kind: resource.Checkbox},
			{Key: "hidden", Label: "Hidden", Kind: resource.Checkbox},
			{Key: "claimed_domains", Label: "Claimed Domains"},
		},
		Fields: []resource.Field{
			{Key: "slug", Label: "Slug", Immutable: true,
				Help: "Unique identifier. Cannot be changed after creation."},
			{Key: "name", Label: "Display Name"},
			{Key: "description", Label: "Description", Kind: resource.TextArea},
			{Key: "enabled", Label: "Enabled", Kind: resource.Checkbox},
			{Key: "hidden", Label: "Hidden from login page", Kind: resource.Checkbox},
			{Key: "auto_create", Label: "Auto-create accounts on first login", Kind: resource.Checkbox},
			{Key: "multi_factor_required", Label: "MFA required", Kind: resource.Checkbox},
			{Key: "email_verification_required", Label: "Email verification required", Kind: resource.Checkbox},
			{Key: "claimed_domains", Label: "Claimed Domains", Kind: resource.StringList,
				Help: "Comma-separated email domains routed to this IDP automatically."},
			{Key: "oidc_authority", Label: "OIDC Authority"},
			{Key: "oidc_client_id", Label: "OIDC Client ID"},
			{Key: "oidc_client_secret", Label: "OIDC Client Secret", Kind: resource.Password},
			{Key: "oidc_scope", Label: "OIDC Scope"},
		},
		List: func(ctx context.Context) ([]resource.Row, error) {
			resp, err := idpStore.ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{})
			if err != nil {
				return nil, err
			}
			rows := make([]resource.Row, 0, len(resp.IDPs))
			for _, idp := range resp.IDPs {
				rows = append(rows, idpRowFromModel(idp))
			}
			return rows, nil
		},
		Get: func(ctx context.Context, id string) (resource.Row, error) {
			resp, err := idpStore.ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{
				Filter: &proto_oidc_idp.Filter{Slug: &proto_types.StringFilterExpression{Eq: id}},
			})
			if err != nil {
				return nil, err
			}
			if len(resp.IDPs) == 0 {
				return nil, errNotFound("IDP", id)
			}
			return idpRowFromModel(resp.IDPs[0]), nil
		},
		Create: func(ctx context.Context, row resource.Row) (resource.Row, error) {
			resp, err := idpStore.CreateIDP(ctx, &proto_oidc_idp.CreateIDPRequest{Idp: idpModelFromRow(row)})
			if err != nil {
				return nil, err
			}
			return idpRowFromModel(resp.Idp), nil
		},
		Update: func(ctx context.Context, id string, row resource.Row) (resource.Row, error) {
			resp, err := idpStore.UpdateIDP(ctx, &proto_oidc_idp.UpdateIDPRequest{
				Slug: id,
				Idp:  idpModelFromRow(row),
			})
			if err != nil {
				return nil, err
			}
			return idpRowFromModel(resp.Idp), nil
		},
		Delete: func(ctx context.Context, id string) error {
			_, err := idpStore.DeleteIDP(ctx, &proto_oidc_idp.DeleteIDPRequest{Slug: id})
			return err
		},
	}
}
