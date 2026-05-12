package identitycreationdenylist

import "context"

type (
	IIdentityCreationDenyListService interface {
		// IsDeniedDomain returns true if the domain is on the deny list.
		// Returns an error only when the list could not be loaded AND IgnoreOnLoadError is false.
		IsDeniedDomain(ctx context.Context, domain string) (bool, error)
		// RefreshNow forces an immediate reload of the external deny list.
		RefreshNow(ctx context.Context) error
	}
)
