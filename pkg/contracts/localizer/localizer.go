package localizer

import (
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	ILocalizerBundle interface {
		GetBundle() *i18n.Bundle
	}
	ILocalizer interface {
		// Initialize is only called from the middleware
		Initialize(c echo.Context)
		GetLocalizer() *i18n.Localizer
	}
)
