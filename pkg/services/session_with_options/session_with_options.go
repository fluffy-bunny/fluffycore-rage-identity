package session_with_options

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_session_with_options "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/session_with_options"
	contracts_sessions "github.com/fluffy-bunny/fluffycore/echo/contracts/sessions"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
)

type (
	service struct {
		factory contracts_sessions.ISessionFactory
		options *contracts_session_with_options.SessionWithOptions
	}
)

var stemService = (*service)(nil)

var _ contracts_session_with_options.ISessionWithOptions = stemService

func (s *service) Ctor(
	factory contracts_sessions.ISessionFactory,
	options *contracts_session_with_options.SessionWithOptions,
) (contracts_session_with_options.ISessionWithOptions, error) {
	return &service{
		factory: factory,
		options: options,
	}, nil
}

func AddScopedISessionWithOptions(cb di.ContainerBuilder, options *contracts_session_with_options.SessionWithOptions) {
	if options == nil {
		panic("options cannot be nil")
	}
	if fluffycore_utils.IsEmptyOrNil(options.Name) {
		panic("options.Name cannot be empty")
	}
	di.AddInstance[*contracts_session_with_options.SessionWithOptions](cb, options)
	di.AddScoped[contracts_session_with_options.ISessionWithOptions](cb, stemService.Ctor)
}
func (s *service) GetSession() (contracts_sessions.ISession, error) {
	session, err := s.factory.GetCookieSession(s.options.Name)
	if err != nil {
		return nil, err
	}
	return session, nil
}
