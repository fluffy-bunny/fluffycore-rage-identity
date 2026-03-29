package htmx

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	support_htmx_audits "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/support/htmx/handlers/audits"
	support_htmx_shell "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/support/htmx/handlers/shell"
)

// AddSupportHandlers registers all HTMX-based support portal handlers.
func AddSupportHandlers(builder di.ContainerBuilder) {
	support_htmx_shell.AddScopedIHandler(builder)
	support_htmx_audits.AddScopedIHandler(builder)
}
