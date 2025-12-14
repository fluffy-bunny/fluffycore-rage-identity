package ComposerBase

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	ComposerBase struct {
		app.Compo
		AppContext context.Context
		Container  di.Container

		Bundle    *i18n.Bundle
		Localizer contracts_Localizer.ILocalizer
	}
)
