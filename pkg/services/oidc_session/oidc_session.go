package oidc_session

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	contracts_oidc_session "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/oidc_session"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
)

type (
	service struct {
		factory              contracts_sessions.ISessionFactory
		wellknownCookieNames contracts_cookies.IWellknownCookieNames
	}
)

var stemService = (*service)(nil)

var _ contracts_oidc_session.IOIDCSession = stemService

func (s *service) Ctor(
	factory contracts_sessions.ISessionFactory,
	wellknownCookieNames contracts_cookies.IWellknownCookieNames,
) (contracts_oidc_session.IOIDCSession, error) {
	return &service{
		factory:              factory,
		wellknownCookieNames: wellknownCookieNames,
	}, nil
}

func AddScopedIOIDCSession(cb di.ContainerBuilder) {
	di.AddScoped[contracts_oidc_session.IOIDCSession](cb, stemService.Ctor)
}
func (s *service) GetSession() (contracts_sessions.ISession, error) {
	session, err := s.factory.GetCookieSession(s.wellknownCookieNames.GetCookieName(contracts_cookies.CookieName_OIDCSession))
	if err != nil {
		return nil, err
	}
	return session, nil
}
