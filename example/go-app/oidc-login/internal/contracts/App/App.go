package App

import (
	"context"

	contracts_GoApp "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/GoApp"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
)

type (
	ContextKey string
	AppContext context.Context

	IApp interface {
		contracts_GoApp.IGoAppComponent
		SetCurrentPage(page contracts_routes.WellknownRoute)
		GetCurrentPage() contracts_routes.WellknownRoute
	}

	IHomeComposer interface {
		contracts_GoApp.IBaseComposer
	}
	IPasswordComposer interface {
		contracts_GoApp.IBaseComposer
	}
	ICreateAccountComposer interface {
		contracts_GoApp.IBaseComposer
	}
	IForgotPasswordComposer interface {
		contracts_GoApp.IBaseComposer
	}
	IResetPasswordComposer interface {
		contracts_GoApp.IBaseComposer
	}
	IVerifyCodeComposer interface {
		contracts_GoApp.IBaseComposer
	}
)
