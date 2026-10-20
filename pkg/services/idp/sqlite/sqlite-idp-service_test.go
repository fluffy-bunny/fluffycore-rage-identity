package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	proto_types "github.com/fluffy-bunny/fluffycore-rage-identity/proto/types"
	"github.com/stretchr/testify/require"
)

// closer matches the unexported *service's Close method without importing it.
type closer interface {
	Close() error
}

func newTestService(t *testing.T, dbPath string, seed *proto_oidc_models.IDPs) proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer {
	t.Helper()
	config := &contracts_config.Config{
		ConfigFiles: contracts_config.ConfigFiles{IDPsDBPath: dbPath},
	}
	svc, err := stemService.Ctor(config, seed)
	require.NoError(t, err)
	t.Cleanup(func() {
		if c, ok := svc.(closer); ok {
			_ = c.Close()
		}
	})
	return svc
}

// asExample resolves the example CRUD interface off the same underlying
// singleton, via CtorExample — exactly the alias path the real DI
// registration uses.
func asExample(t *testing.T, svc proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer) proto_oidc_idp.IFluffyCoreExampleIDPStoreServer {
	t.Helper()
	writable, err := stemService.CtorExample(svc)
	require.NoError(t, err)
	return writable
}

func TestSeedsFromJSON(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "idps.db")
	seed := &proto_oidc_models.IDPs{
		Idps: []*proto_oidc_models.IDP{
			{Slug: "evilcorp", Enabled: true, Name: "EvilCorp SSO"},
		},
	}

	svc := newTestService(t, dbPath, seed)
	resp, err := svc.ListIDP(context.Background(), &proto_oidc_idp.ListIDPRequest{})
	require.NoError(t, err)
	require.Len(t, resp.IDPs, 1)
	require.Equal(t, "evilcorp", resp.IDPs[0].Slug)
}

func TestSyncSeedOnlyAddsMissingSlugs(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "idps.db")
	seed := &proto_oidc_models.IDPs{
		Idps: []*proto_oidc_models.IDP{
			{Slug: "evilcorp", Enabled: true, Name: "EvilCorp SSO"},
		},
	}
	svc := newTestService(t, dbPath, seed)

	// Simulate an operator/customer edit to the seeded IDP.
	writable := asExample(t, svc)
	_, err := writable.UpdateIDP(context.Background(), &proto_oidc_idp.UpdateIDPRequest{
		Slug: "evilcorp",
		Idp:  &proto_oidc_models.IDP{Slug: "evilcorp", Enabled: true, Name: "EvilCorp SSO (edited by operator)"},
	})
	require.NoError(t, err)

	// "Restart": re-run Ctor against the SAME db file, with a seed list that
	// (a) still contains the original evilcorp entry (now stale relative to
	// the DB) and (b) adds a brand-new slug. The already-present slug must
	// NOT be overwritten by the stale JSON; the new slug must be added.
	secondBootSeed := &proto_oidc_models.IDPs{
		Idps: []*proto_oidc_models.IDP{
			{Slug: "evilcorp", Enabled: true, Name: "EvilCorp SSO"}, // stale — must be ignored
			{Slug: "newco", Enabled: true, Name: "NewCo SSO"},       // new — must be added
		},
	}
	svc2 := newTestService(t, dbPath, secondBootSeed)
	resp, err := svc2.ListIDP(context.Background(), &proto_oidc_idp.ListIDPRequest{})
	require.NoError(t, err)
	require.Len(t, resp.IDPs, 2)

	bySlug := map[string]*proto_oidc_models.IDP{}
	for _, idpEntry := range resp.IDPs {
		bySlug[idpEntry.Slug] = idpEntry
	}
	require.Equal(t, "EvilCorp SSO (edited by operator)", bySlug["evilcorp"].Name, "database edit must survive re-sync from stale JSON")
	require.Equal(t, "NewCo SSO", bySlug["newco"].Name, "a genuinely new slug in JSON must still be picked up")
}

func TestCreateUpdateDeleteRoundTrip(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "idps.db")
	svc := newTestService(t, dbPath, nil)
	writable := asExample(t, svc)

	ctx := context.Background()

	// Create. Slugs are stored exactly as given — this store does no
	// case-normalization, matching pkg/services/idp/inmemory's behavior.
	createResp, err := writable.CreateIDP(ctx, &proto_oidc_idp.CreateIDPRequest{
		Idp: &proto_oidc_models.IDP{Slug: "acmeco", Enabled: true, Name: "Acme"},
	})
	require.NoError(t, err)
	require.Equal(t, "acmeco", createResp.Idp.Slug)

	// Duplicate create (same exact slug) must fail
	_, err = writable.CreateIDP(ctx, &proto_oidc_idp.CreateIDPRequest{
		Idp: &proto_oidc_models.IDP{Slug: "acmeco", Enabled: true},
	})
	require.Error(t, err)

	// Update
	updateResp, err := writable.UpdateIDP(ctx, &proto_oidc_idp.UpdateIDPRequest{
		Slug: "acmeco",
		Idp:  &proto_oidc_models.IDP{Slug: "acmeco", Enabled: false, Name: "Acme Renamed"},
	})
	require.NoError(t, err)
	require.False(t, updateResp.Idp.Enabled)
	require.Equal(t, "Acme Renamed", updateResp.Idp.Name)

	listResp, err := svc.ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{})
	require.NoError(t, err)
	require.Len(t, listResp.IDPs, 1)
	require.Equal(t, "Acme Renamed", listResp.IDPs[0].Name)

	// An UpdateIDP that tries to change the slug itself must be rejected
	// rather than silently creating an orphaned second row.
	_, err = writable.UpdateIDP(ctx, &proto_oidc_idp.UpdateIDPRequest{
		Slug: "acmeco",
		Idp:  &proto_oidc_models.IDP{Slug: "renamed-acmeco"},
	})
	require.Error(t, err)

	// Update of an unknown slug must fail
	_, err = writable.UpdateIDP(ctx, &proto_oidc_idp.UpdateIDPRequest{
		Slug: "does-not-exist",
		Idp:  &proto_oidc_models.IDP{Slug: "does-not-exist"},
	})
	require.Error(t, err)

	// Delete
	_, err = writable.DeleteIDP(ctx, &proto_oidc_idp.DeleteIDPRequest{Slug: "acmeco"})
	require.NoError(t, err)

	listResp2, err := svc.ListIDP(ctx, &proto_oidc_idp.ListIDPRequest{})
	require.NoError(t, err)
	require.Empty(t, listResp2.IDPs)

	// Delete of an already-deleted slug must fail
	_, err = writable.DeleteIDP(ctx, &proto_oidc_idp.DeleteIDPRequest{Slug: "acmeco"})
	require.Error(t, err)
}

func TestListIDPFilterBySlug(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "idps.db")
	seed := &proto_oidc_models.IDPs{
		Idps: []*proto_oidc_models.IDP{
			{Slug: "one", Enabled: true},
			{Slug: "two", Enabled: false},
		},
	}
	svc := newTestService(t, dbPath, seed)

	resp, err := svc.ListIDP(context.Background(), &proto_oidc_idp.ListIDPRequest{
		Filter: &proto_oidc_idp.Filter{
			Slug: &proto_types.StringFilterExpression{Eq: "one"},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp.IDPs, 1)
	require.Equal(t, "one", resp.IDPs[0].Slug)
}
