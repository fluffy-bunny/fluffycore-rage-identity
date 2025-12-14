package Localizer

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	services_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/services/i18n/LocalizerBundle"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	service struct {
		bundle    *i18n.Bundle
		localizer *i18n.Localizer
	}
)

var stemService = (*service)(nil)

var _ contracts_Localizer.ILocalizer = stemService

func (s *service) Ctor(localizerBundle contracts_Localizer.ILocalizerBundle) (contracts_Localizer.ILocalizer, error) {

	bundle := localizerBundle.GetBundle()
	svc := &service{
		bundle: bundle,
	}
	svc.Initialize()
	return svc, nil
}

func AddScopedILocalizer(cb di.ContainerBuilder) {
	services_LocalizerBundle.AddScopedILocalizerBundle(cb)
	di.AddScoped[contracts_Localizer.ILocalizer](cb, stemService.Ctor)
}
func (s *service) GetLocalizer() *i18n.Localizer {
	if s.localizer == nil {
		s.Initialize()
	}
	return s.localizer
}
func (s *service) Initialize(accept ...string) {
	s.localizer = i18n.NewLocalizer(s.bundle, accept...)
}

func (s *service) GetLocalizedString(localKey contracts_LocalizerBundle.LocaleKey) string {
	d, err := s.TryGetLocalizedString(localKey)
	if err != nil {
		panic(err)
	}
	return d
}

func (s *service) TryGetLocalizedString(localKey contracts_LocalizerBundle.LocaleKey) (string, error) {
	localizer := s.GetLocalizer()
	return localizer.LocalizeMessage(&i18n.Message{
		ID: string(localKey),
	})
}

func (s *service) GetLocalizedStringF(localKey contracts_LocalizerBundle.LocaleKey, templateData map[string]interface{}) string {
	d, err := s.TryGetLocalizedStringF(localKey, templateData)
	if err != nil {
		panic(err)
	}
	return d
}

func (s *service) TryGetLocalizedStringF(localKey contracts_LocalizerBundle.LocaleKey, templateData map[string]interface{}) (string, error) {
	localizer := s.GetLocalizer()
	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    string(localKey),
		TemplateData: templateData,
	})
}
