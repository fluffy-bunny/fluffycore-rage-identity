package Localizer

import (
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	ILocalizerBundle interface {
		GetBundle() *i18n.Bundle
	}

	ILocalizer interface {
		// Initialize is only called from the middleware
		Initialize(accept ...string)
		GetLocalizer() *i18n.Localizer

		GetLocalizedString(localKey contracts_LocalizerBundle.LocaleKey) string
		TryGetLocalizedString(localKey contracts_LocalizerBundle.LocaleKey) (string, error)

		GetLocalizedStringF(localKey contracts_LocalizerBundle.LocaleKey, templateData map[string]interface{}) string
		TryGetLocalizedStringF(localKey contracts_LocalizerBundle.LocaleKey, templateData map[string]interface{}) (string, error)
	}
)
