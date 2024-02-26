package oidc_session

import (
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
)

type (
	IOIDCSession interface {
		GetSession() (contracts_sessions.ISession, error)
	}
)
