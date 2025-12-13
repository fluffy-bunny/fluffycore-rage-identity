package LocalizerBundle

import (
	"github.com/BurntSushi/toml"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/LocalizerBundle"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	language "golang.org/x/text/language"
)

type (
	service struct {
		bundle *i18n.Bundle
	}
)

var stemService = (*service)(nil)

var _ contracts_Localizer.ILocalizerBundle = stemService

func (s *service) Ctor() (contracts_Localizer.ILocalizerBundle, error) {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load all embedded locale files
	for _, localeFile := range contracts_LocalizerBundle.LocaleFiles {
		_, err := bundle.LoadMessageFileFS(contracts_LocalizerBundle.LocalFS, localeFile)
		if err != nil {
			return nil, err
		}
	}

	return &service{
		bundle: bundle,
	}, nil
}

func AddScopedILocalizerBundle(cb di.ContainerBuilder) {
	di.AddScoped[contracts_Localizer.ILocalizerBundle](cb, stemService.Ctor)
}
func (s *service) GetBundle() *i18n.Bundle {
	return s.bundle
}
