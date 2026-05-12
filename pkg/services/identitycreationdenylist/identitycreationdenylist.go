package identitycreationdenylist

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_identitycreationdenylist "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identitycreationdenylist"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	zerolog "github.com/rs/zerolog"
)

type denyListJSON struct {
	Version string   `json:"version"`
	Updated string   `json:"updated"`
	Domains []string `json:"domains"`
}

type service struct {
	config *contracts_config.Config

	mu         sync.RWMutex
	cached     []string  // domains loaded from external source
	lastLoaded time.Time // zero value means never loaded
	loadErr    error     // last load error (non-nil only when load failed)
}

var stemService = (*service)(nil)

var _ contracts_identitycreationdenylist.IIdentityCreationDenyListService = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
) (contracts_identitycreationdenylist.IIdentityCreationDenyListService, error) {
	return &service{
		config: config,
	}, nil
}

func AddSingletonIIdentityCreationDenyListService(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_identitycreationdenylist.IIdentityCreationDenyListService](cb, stemService.Ctor)
}

func (s *service) denyListConfig() *contracts_config.IdentityCreationDenyListConfig {
	if s.config.IdentityCreationDenyListConfig == nil {
		return &contracts_config.IdentityCreationDenyListConfig{
			IgnoreOnLoadError: true,
			CacheTTLSeconds:   3600,
		}
	}
	return s.config.IdentityCreationDenyListConfig
}

// load fetches the deny list from URL or file. Caller must NOT hold mu.
func (s *service) load(ctx context.Context) error {
	cfg := s.denyListConfig()
	log := zerolog.Ctx(ctx).With().Str("service", "IdentityCreationDenyListService").Logger()

	var data []byte
	var loadErr error

	switch {
	case cfg.URL != "":
		log.Info().Str("url", cfg.URL).Msg("loading identity creation deny list from URL")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.URL, nil)
		if err != nil {
			loadErr = err
			break
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			loadErr = err
			break
		}
		defer resp.Body.Close()
		data, loadErr = io.ReadAll(resp.Body)
		if loadErr == nil && resp.StatusCode != http.StatusOK {
			loadErr = &httpStatusError{code: resp.StatusCode, url: cfg.URL}
		}

	case cfg.FilePath != "":
		log.Info().Str("path", cfg.FilePath).Msg("loading identity creation deny list from file")
		data, loadErr = os.ReadFile(cfg.FilePath)

	default:
		// no external source configured — nothing to load
		s.mu.Lock()
		s.cached = nil
		s.lastLoaded = time.Now()
		s.loadErr = nil
		s.mu.Unlock()
		return nil
	}

	if loadErr != nil {
		log.Error().Err(loadErr).Msg("failed to load identity creation deny list")
		s.mu.Lock()
		s.loadErr = loadErr
		s.mu.Unlock()
		return loadErr
	}

	var dl denyListJSON
	if err := json.Unmarshal(data, &dl); err != nil {
		log.Error().Err(err).Msg("failed to parse identity creation deny list JSON")
		s.mu.Lock()
		s.loadErr = err
		s.mu.Unlock()
		return err
	}

	log.Info().
		Str("version", dl.Version).
		Str("updated", dl.Updated).
		Int("domain_count", len(dl.Domains)).
		Msg("identity creation deny list loaded successfully")

	s.mu.Lock()
	s.cached = dl.Domains
	s.lastLoaded = time.Now()
	s.loadErr = nil
	s.mu.Unlock()
	return nil
}

// ensureFresh loads (or reloads after TTL) the deny list. Returns last load error when relevant.
func (s *service) ensureFresh(ctx context.Context) error {
	cfg := s.denyListConfig()
	ttl := time.Duration(cfg.CacheTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = time.Hour
	}

	s.mu.RLock()
	stale := s.lastLoaded.IsZero() || time.Since(s.lastLoaded) > ttl
	s.mu.RUnlock()

	if !stale {
		return nil
	}

	// Upgrade to write lock via a full reload
	return s.load(ctx)
}

func (s *service) IsDeniedDomain(ctx context.Context, domain string) (bool, error) {
	log := zerolog.Ctx(ctx).With().Str("service", "IdentityCreationDenyListService").Logger()
	cfg := s.denyListConfig()

	// Always check the static list first (fast, no I/O).
	if utils.IsDeniedDomain(domain, s.config.DeniedDomains) {
		return true, nil
	}

	// If no external source is configured, we are done.
	if cfg.URL == "" && cfg.FilePath == "" {
		return false, nil
	}

	err := s.ensureFresh(ctx)
	if err != nil {
		if cfg.IgnoreOnLoadError {
			log.Warn().Err(err).
				Str("domain", domain).
				Msg("identity creation deny list unavailable — failing open, allowing domain")
			return false, nil
		}
		log.Error().Err(err).
			Str("domain", domain).
			Msg("identity creation deny list unavailable — failing closed, denying signup")
		return false, err
	}

	s.mu.RLock()
	cached := s.cached
	s.mu.RUnlock()

	d := strings.ToLower(domain)
	for _, denied := range cached {
		if strings.ToLower(denied) == d {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) RefreshNow(ctx context.Context) error {
	s.mu.Lock()
	s.lastLoaded = time.Time{} // force stale
	s.mu.Unlock()
	return s.load(ctx)
}

// httpStatusError is a simple error type for non-200 HTTP responses.
type httpStatusError struct {
	code int
	url  string
}

func (e *httpStatusError) Error() string {
	return "HTTP " + http.StatusText(e.code) + " fetching " + e.url
}
