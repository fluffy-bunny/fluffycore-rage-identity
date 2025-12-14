package GoApp

import (
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type (
	IBaseComposer interface {
		app.Composer
	}
	IGoAppComponent interface {
		app.Composer
		OnMount(ctx app.Context)
	}
)
