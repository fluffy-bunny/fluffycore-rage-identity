package session_with_options

import (
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
)

type (
	SessionWithOptions struct {
		Name string
	}
	ISessionWithOptions interface {
		GetSession() (contracts_sessions.ISession, error)
	}
)
