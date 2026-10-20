package operatorportal

import (
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

// errNotFound is shared by resource adapters (idp_resource.go,
// client_resource.go) whose backing List RPC doubles as Get (via a filter)
// and so needs to synthesize a NotFound itself when the filter comes back
// empty, rather than relying on the backing service to return one.
func errNotFound(kind, id string) error {
	return status.Errorf(codes.NotFound, "%s %q not found", kind, id)
}
