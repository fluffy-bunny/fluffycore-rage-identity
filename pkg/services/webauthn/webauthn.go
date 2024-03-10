package webauthn

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_webauthn "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/webauthn"
	go_webauthn "github.com/go-webauthn/webauthn/webauthn"
)

type (
	service struct {
		config *contracts_webauthn.WebAuthNConfig
		w      *go_webauthn.WebAuthn
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_webauthn.IWebAuthN = stemService
}
func (s *service) Ctor(
	config *contracts_webauthn.WebAuthNConfig,
) (contracts_webauthn.IWebAuthN, error) {
	wConfig := &go_webauthn.Config{
		RPDisplayName: config.RPDisplayName, // Display Name for your site
		RPID:          config.RPID,          // Generally the FQDN for your site
		RPOrigins:     config.RPOrigins,     // The origin URLs allowed for WebAuthn requests
	}
	w, err := go_webauthn.New(wConfig)
	if err != nil {
		return nil, err
	}
	ss := &service{
		w:      w,
		config: config,
	}
	return ss, nil
}

func AddSingletonIWebAuthN(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_webauthn.IWebAuthN](cb, stemService.Ctor)
}

func (s *service) GetWebAuthN() *go_webauthn.WebAuthn {
	return s.w
}
