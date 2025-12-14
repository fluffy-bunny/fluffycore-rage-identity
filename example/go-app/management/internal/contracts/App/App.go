package App

import (
	"context"

	contracts_GoApp "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/GoApp"
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/routes"
)

type (
	AppContext context.Context

	IApp interface {
		contracts_GoApp.IGoAppComponent
		SetCurrentPage(page contracts_routes.WellknownRoute)
		GetCurrentPage() contracts_routes.WellknownRoute
		IsAuthenticated() bool
		LoginWithReturnURL(returnURL string)
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
	IVerifyCodeComposer interface {
		contracts_GoApp.IBaseComposer
	}
	IProfileComposer interface {
		contracts_GoApp.IBaseComposer
	}
	IPasswordManagerComposer interface {
		contracts_GoApp.IBaseComposer
	}
	ILinkedAccountsComposer interface {
		contracts_GoApp.IBaseComposer
	}
)
