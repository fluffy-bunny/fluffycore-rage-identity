package utils

import (
	"strings"

	services_handlers_shared "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/shared"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func ReplaceStrings(original string, replace map[string]string) string {
	if replace == nil {
		return original
	}
	for k, v := range replace {
		original = strings.ReplaceAll(original, k, v)
	}
	return original
}

func LocalizeReplaceStrings(localizer *i18n.Localizer, id string, replace map[string]string) (string, error) {
	template, err := localizer.LocalizeMessage(&i18n.Message{ID: id})
	if err != nil {
		return id, err
	}

	s := ReplaceStrings(template, replace)
	return s, nil
}

func LocalizeSimple(localizer *i18n.Localizer, id string) string {
	s, err := localizer.LocalizeMessage(&i18n.Message{ID: id})
	if err != nil {
		return id
	}
	return s
}

func LocalizeToError(localizer *i18n.Localizer, id string, replace map[string]string) *services_handlers_shared.Error {
	msg, err := LocalizeReplaceStrings(localizer, id, replace)
	if err == nil {
		services_handlers_shared.NewErrorF(id, msg)
	}
	return services_handlers_shared.NewErrorF(id, id)
}
