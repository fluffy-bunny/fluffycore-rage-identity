package utils

import (
	"strings"

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
