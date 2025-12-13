package GoApp

import (
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type (
	IGoAppComponent interface {
		app.Composer
		OnMount(ctx app.Context)
	}
	IBaseComposer interface {
		app.Composer
	}
)
