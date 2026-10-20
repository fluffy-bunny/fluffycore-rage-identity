// Package sqlite implements a persistent, writable IDP store backed by an
// embedded, pure-Go SQLite database (modernc.org/sqlite — no cgo).
//
// This is the registered proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer:
// every existing consumer of ListIDP (login flows, claims minting, etc.) only
// ever sees that minimal, ListIDP-only contract — that's the one interface
// every consumer of this library must implement, however they choose to
// persist IDPs.
//
// The SAME concrete instance also implements
// proto_oidc_idp.IFluffyCoreExampleIDPStoreServer — a full worked
// example of a richer CRUD service, registered separately (see
// AddSingletonIFluffyCoreExampleIDPStoreServer below), for callers
// that need to mutate (the operator portal, and later the customer
// self-service portal) or for anyone using this project as a reference for
// how to build their own backing store. There's exactly one DI registration
// per interface, but only one underlying *service/one DB handle — the
// example-service registration is wired to depend on (not separately
// construct) the singleton registration, so they can never diverge.
//
// idps.json (ConfigFiles.IDPsPath) is synced into this database on every
// boot: each IDP present in the JSON is inserted only if its slug isn't
// already in the database. Once an IDP exists in the database, idps.json no
// longer has any effect on it — the database is the source of truth for
// anything it already holds, so operator/customer-made edits always survive
// a restart even though the JSON file keeps being read.
package sqlite

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	proto_oidc_idp "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/idp"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
	protojson "google.golang.org/protobuf/encoding/protojson"
	_ "modernc.org/sqlite"
)

type (
	service struct {
		proto_oidc_idp.UnimplementedSingletonIDPServiceServer
		proto_oidc_idp.UnimplementedExampleIDPStoreServer

		mu sync.RWMutex
		db *sql.DB
	}
)

var stemService = (*service)(nil)
var _ proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer = stemService
var _ proto_oidc_idp.IFluffyCoreExampleIDPStoreServer = stemService

func (s *service) Ctor(config *contracts_config.Config, seedIdps *proto_oidc_models.IDPs) (proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer, error) {
	dbPath := config.ConfigFiles.IDPsDBPath
	if dbPath == "" {
		dbPath = "./data/idps.db"
	}
	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create IDP store directory %s: %v", dir, err)
		}
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to open IDP store %s: %v", dbPath, err)
	}
	// SQLite only tolerates one writer at a time; this store is small and
	// low-traffic (operator/customer-driven config edits, not request-path
	// writes), so a single shared connection avoids "database is locked"
	// errors under the driver's own concurrent-connection pool.
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS idps (
		slug TEXT PRIMARY KEY,
		data TEXT NOT NULL
	)`); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create idps table: %v", err)
	}

	svc := &service{db: db}
	if err := svc.syncSeed(seedIdps); err != nil {
		return nil, err
	}
	return svc, nil
}

// CtorExample exposes the SAME already-constructed singleton (see Ctor
// above) as IFluffyCoreExampleIDPStoreServer, by depending on the
// IFluffyCoreSingletonIDPServiceServer registration rather than constructing
// a second instance. DI resolves/caches that dependency once regardless of
// how many other registrations depend on it, so this guarantees a single
// underlying *service/one DB handle no matter which interface a caller asks
// the container for.
func (s *service) CtorExample(underlying proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer) (proto_oidc_idp.IFluffyCoreExampleIDPStoreServer, error) {
	typed, ok := underlying.(proto_oidc_idp.IFluffyCoreExampleIDPStoreServer)
	if !ok {
		return nil, status.Error(codes.Internal, "the registered IDP store does not implement the example CRUD extension")
	}
	return typed, nil
}

func AddSingletonIFluffyCoreSingletonIDPServiceServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_idp.IFluffyCoreSingletonIDPServiceServer](cb, stemService.Ctor)
}

// AddSingletonIFluffyCoreExampleIDPStoreServer registers the example
// CRUD interface as an ALIAS of the IFluffyCoreSingletonIDPServiceServer
// singleton — call AddSingletonIFluffyCoreSingletonIDPServiceServer first (or
// at all) in the same container, or resolution will fail.
func AddSingletonIFluffyCoreExampleIDPStoreServer(cb di.ContainerBuilder) {
	di.AddSingleton[proto_oidc_idp.IFluffyCoreExampleIDPStoreServer](cb, stemService.CtorExample)
}

// Close releases the underlying database handle. Not part of either
// registered interface — callers that need it (graceful shutdown, tests)
// type-assert to io.Closer.
func (s *service) Close() error {
	return s.db.Close()
}

// syncSeed imports every IDP from idps.json whose slug isn't already present
// in the database. Runs on every boot, not just once — so adding a new IDP
// to idps.json later still gets picked up, but anything already in the
// database (including anything an operator/customer has since edited) is
// left untouched, even if the JSON version of it has since changed.
func (s *service) syncSeed(seedIdps *proto_oidc_models.IDPs) error {
	if seedIdps == nil || len(seedIdps.Idps) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, idpEntry := range seedIdps.Idps {
		var count int
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM idps WHERE slug = ?`, idpEntry.Slug).Scan(&count); err != nil {
			return status.Errorf(codes.Internal, "failed to check existing IDP %s while seeding: %v", idpEntry.Slug, err)
		}
		if count > 0 {
			// Already present — the database wins, idps.json is not
			// re-applied over it.
			continue
		}
		if err := s.insertOrReplace(idpEntry); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) insertOrReplace(idpEntry *proto_oidc_models.IDP) error {
	b, err := protojson.Marshal(idpEntry)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to marshal IDP %s: %v", idpEntry.Slug, err)
	}
	_, err = s.db.Exec(`INSERT INTO idps (slug, data) VALUES (?, ?)
		ON CONFLICT(slug) DO UPDATE SET data = excluded.data`,
		idpEntry.Slug, string(b))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to persist IDP %s: %v", idpEntry.Slug, err)
	}
	return nil
}

func (s *service) loadAll() ([]*proto_oidc_models.IDP, error) {
	rows, err := s.db.Query(`SELECT data FROM idps`)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query idps: %v", err)
	}
	defer rows.Close()

	var result []*proto_oidc_models.IDP
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan idp row: %v", err)
		}
		idpEntry := &proto_oidc_models.IDP{}
		if err := protojson.Unmarshal([]byte(data), idpEntry); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal idp row: %v", err)
		}
		result = append(result, idpEntry)
	}
	return result, rows.Err()
}

// ListIDP mirrors the filter semantics of pkg/services/idp/inmemory exactly
// (including its slug/claimed-domain lower-casing quirk on the FILTER side
// only — stored slugs are never case-normalized, matching the original).
func (s *service) ListIDP(ctx context.Context, request *proto_oidc_idp.ListIDPRequest) (*proto_oidc_idp.ListIDPResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if request.Filter != nil {
		if request.Filter.Slug != nil {
			request.Filter.Slug.Eq = strings.ToLower(request.Filter.Slug.Eq)
			for i, v := range request.Filter.Slug.In {
				request.Filter.Slug.In[i] = strings.ToLower(v)
			}
		}
		if request.Filter.ClaimedDomains != nil {
			request.Filter.ClaimedDomains.Eq = strings.ToLower(request.Filter.ClaimedDomains.Eq)
			for i, v := range request.Filter.ClaimedDomains.In {
				request.Filter.ClaimedDomains.In[i] = strings.ToLower(v)
			}
		}
	}

	s.mu.RLock()
	all, err := s.loadAll()
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	var result []*proto_oidc_models.IDP
	for _, c := range all {
		if request.Filter != nil {
			if request.Filter.Enabled != nil && request.Filter.Enabled.Eq != c.Enabled {
				continue
			}
			if request.Filter.Hidden != nil && request.Filter.Hidden.Eq != c.Hidden {
				continue
			}
			if request.Filter.Slug != nil {
				if request.Filter.Slug.Eq != c.Slug {
					continue
				}
				if len(request.Filter.Slug.In) > 0 {
					hit := false
					for _, v := range request.Filter.Slug.In {
						if v == c.Slug {
							hit = true
							break
						}
					}
					if !hit {
						continue
					}
				}
			}
			if request.Filter.ClaimedDomains != nil {
				domainSet := make(map[string]bool, len(c.ClaimedDomains))
				for _, v := range c.ClaimedDomains {
					domainSet[v] = true
				}
				if request.Filter.ClaimedDomains.Eq != "" && !domainSet[request.Filter.ClaimedDomains.Eq] {
					continue
				}
				if len(request.Filter.ClaimedDomains.In) > 0 {
					hit := false
					for _, v := range request.Filter.ClaimedDomains.In {
						if domainSet[v] {
							hit = true
							break
						}
					}
					if !hit {
						continue
					}
				}
			}
		}
		result = append(result, c)
	}

	return &proto_oidc_idp.ListIDPResponse{IDPs: result}, nil
}

func (s *service) validateIDP(idpEntry *proto_oidc_models.IDP) error {
	if idpEntry == nil {
		return status.Error(codes.InvalidArgument, "IDP is required")
	}
	if idpEntry.Slug == "" {
		return status.Error(codes.InvalidArgument, "IDP.Slug is required")
	}
	return nil
}

func (s *service) CreateIDP(ctx context.Context, request *proto_oidc_idp.CreateIDPRequest) (*proto_oidc_idp.CreateIDPResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := s.validateIDP(request.Idp); err != nil {
		return nil, err
	}
	slug := request.Idp.Slug

	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM idps WHERE slug = ?`, slug).Scan(&count); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing IDP: %v", err)
	}
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "IDP with slug %s already exists", slug)
	}
	if err := s.insertOrReplace(request.Idp); err != nil {
		return nil, err
	}
	return &proto_oidc_idp.CreateIDPResponse{Idp: request.Idp}, nil
}

func (s *service) UpdateIDP(ctx context.Context, request *proto_oidc_idp.UpdateIDPRequest) (*proto_oidc_idp.UpdateIDPResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if request.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "Slug is required")
	}
	if err := s.validateIDP(request.Idp); err != nil {
		return nil, err
	}
	slug := request.Slug

	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM idps WHERE slug = ?`, slug).Scan(&count); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing IDP: %v", err)
	}
	if count == 0 {
		return nil, status.Errorf(codes.NotFound, "IDP with slug %s not found", slug)
	}
	if request.Idp.Slug != "" && request.Idp.Slug != slug {
		return nil, status.Errorf(codes.InvalidArgument,
			"IDP.Slug (%s) must match the Slug being updated (%s); renaming a slug is not supported", request.Idp.Slug, slug)
	}
	// Force the stored row's slug to match the key being updated, so a
	// caller that left Idp.Slug empty doesn't end up with a row keyed by "".
	request.Idp.Slug = slug
	if err := s.insertOrReplace(request.Idp); err != nil {
		return nil, err
	}
	return &proto_oidc_idp.UpdateIDPResponse{Idp: request.Idp}, nil
}

func (s *service) DeleteIDP(ctx context.Context, request *proto_oidc_idp.DeleteIDPRequest) (*proto_oidc_idp.DeleteIDPResponse, error) {
	if request == nil || request.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "Slug is required")
	}
	slug := request.Slug

	s.mu.Lock()
	defer s.mu.Unlock()

	res, err := s.db.Exec(`DELETE FROM idps WHERE slug = ?`, slug)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete IDP %s: %v", slug, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, status.Errorf(codes.NotFound, "IDP with slug %s not found", slug)
	}
	return &proto_oidc_idp.DeleteIDPResponse{}, nil
}
