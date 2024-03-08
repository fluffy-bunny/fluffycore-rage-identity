package utils

import (
	"fmt"
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
func CurlyBraceReplaceStrings(original string, replace map[string]string) string {
	if replace == nil {
		return original
	}
	for k, v := range replace {
		original = strings.ReplaceAll(original, fmt.Sprintf("{%s}", k), v)
	}
	return original
}

func LocalizeWithInterperlate(localizer *i18n.Localizer, id string, replace map[string]string) string {
	template, err := localizer.LocalizeMessage(&i18n.Message{ID: id})
	if err != nil {
		return id
	}

	s := CurlyBraceReplaceStrings(template, replace)
	return s
}

func LocalizeSimple(localizer *i18n.Localizer, id string) string {
	return LocalizeWithInterperlate(localizer, id, nil)
}
