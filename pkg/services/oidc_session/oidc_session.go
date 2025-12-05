package oidc_session

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
)

type (
	service struct {
		factory contracts_sessions.ISessionFactory
	}
)

var stemService = (*service)(nil)

var _ contracts_oidc_session.IOIDCSession = stemService

func (s *service) Ctor(
	factory contracts_sessions.ISessionFactory,
) (contracts_oidc_session.IOIDCSession, error) {
	return &service{
		factory: factory,
	}, nil
}

func AddScopedIOIDCSession(cb di.ContainerBuilder) {
	di.AddScoped[contracts_oidc_session.IOIDCSession](cb, stemService.Ctor)
}
func (s *service) GetSession() (contracts_sessions.ISession, error) {
	session, err := s.factory.GetCookieSession(models.OIDCSessionName)
	if err != nil {
		return nil, err
	}
	return session, nil
}
