package identitycreationprotection

import "context"

type (
	IIdentityCreationProtection interface {
		// IsDisposableEmailDomain returns true if the domain is a known disposable/throwaway
		// or otherwise denied email domain.
		// Returns an error only when the blocklist could not be loaded AND IgnoreOnLoadError is false.
		IsDisposableEmailDomain(ctx context.Context, domain string) (bool, error)
	}
)
