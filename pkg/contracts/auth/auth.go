package auth

type (
	IRequiresNoAuth interface {
		GetAuthMap() map[string]bool
	}
)
