package localizerbundle

import (
	"encoding/json"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/localizer"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	language "golang.org/x/text/language"
)

type (
	service struct {
		bundle *i18n.Bundle
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_localizer.ILocalizerBundle = stemService
}
func (s *service) Ctor() (contracts_localizer.ILocalizerBundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("resources/en.json")
	//bundle.LoadMessageFile("resources/fr-FR.json")
	return &service{
		bundle: bundle,
	}, nil
}

func AddSingletonILocalizerBundle(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_localizer.ILocalizerBundle](cb, stemService.Ctor)
}
func (s *service) GetBundle() *i18n.Bundle {
	return s.bundle
}
