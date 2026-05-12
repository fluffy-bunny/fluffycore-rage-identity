package identitycreationprotection

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_identitycreationprotection "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identitycreationprotection"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/utils"
	zerolog "github.com/rs/zerolog"
)

type service struct {
	config *contracts_config.Config

	mu         sync.RWMutex
	cached     map[string]struct{} // domain set loaded from remote blocklist
	lastLoaded time.Time           // zero means never loaded
	loadErr    error
}

var stemService = (*service)(nil)

var _ contracts_identitycreationprotection.IIdentityCreationProtection = stemService

func (s *service) Ctor(
	config *contracts_config.Config,
) (contracts_identitycreationprotection.IIdentityCreationProtection, error) {
	return &service{
		config: config,
	}, nil
}

func AddSingletonIIdentityCreationProtection(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_identitycreationprotection.IIdentityCreationProtection](cb, stemService.Ctor)
}

func (s *service) protectionConfig() *contracts_config.IdentityCreationProtectionConfig {
	if s.config.IdentityCreationProtectionConfig == nil {
		return &contracts_config.IdentityCreationProtectionConfig{
			IgnoreOnLoadError: true,
			CacheTTLSeconds:   3600,
		}
	}
	return s.config.IdentityCreationProtectionConfig
}

// load fetches and parses the plain-text blocklist. Each non-empty, non-comment line is a domain.
func (s *service) load(ctx context.Context) error {
	cfg := s.protectionConfig()
	log := zerolog.Ctx(ctx).With().Str("service", "IdentityCreationProtection").Logger()

	if cfg.DisposableEmailListURL == "" {
		s.mu.Lock()
		s.cached = nil
		s.lastLoaded = time.Now()
		s.loadErr = nil
		s.mu.Unlock()
		return nil
	}

	log.Info().Str("url", cfg.DisposableEmailListURL).Msg("loading disposable email blocklist")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.DisposableEmailListURL, nil)
	if err != nil {
		return s.setLoadErr(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return s.setLoadErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP %d fetching disposable email blocklist", resp.StatusCode)
		return s.setLoadErr(err)
	}

	domains := parseBlocklist(resp.Body)

	log.Info().Int("domain_count", len(domains)).Msg("disposable email blocklist loaded")

	s.mu.Lock()
	s.cached = domains
	s.lastLoaded = time.Now()
	s.loadErr = nil
	s.mu.Unlock()
	return nil
}

func (s *service) setLoadErr(err error) error {
	s.mu.Lock()
	s.loadErr = err
	s.mu.Unlock()
	return err
}

// parseBlocklist reads a plain-text list of domains (one per line). Lines that are empty
// or begin with '#' are skipped.
func parseBlocklist(r io.Reader) map[string]struct{} {
	domains := make(map[string]struct{})
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		domains[strings.ToLower(line)] = struct{}{}
	}
	return domains
}

func (s *service) ensureFresh(ctx context.Context) error {
	cfg := s.protectionConfig()
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

	return s.load(ctx)
}

func (s *service) IsDisposableEmailDomain(ctx context.Context, domain string) (bool, error) {
	log := zerolog.Ctx(ctx).With().Str("service", "IdentityCreationProtection").Logger()
	cfg := s.protectionConfig()

	// Always check the static deny list first (fast, no I/O).
	if utils.IsDeniedDomain(domain, s.config.DeniedDomains) {
		return true, nil
	}

	// If protection is disabled, skip the remote check entirely.
	if !cfg.Enabled {
		return false, nil
	}

	err := s.ensureFresh(ctx)
	if err != nil {
		if cfg.IgnoreOnLoadError {
			log.Warn().Err(err).Str("domain", domain).
				Msg("disposable email blocklist unavailable — failing open, allowing domain")
			return false, nil
		}
		log.Error().Err(err).Str("domain", domain).
			Msg("disposable email blocklist unavailable — failing closed")
		return false, err
	}

	s.mu.RLock()
	_, denied := s.cached[strings.ToLower(domain)]
	s.mu.RUnlock()

	return denied, nil
}
