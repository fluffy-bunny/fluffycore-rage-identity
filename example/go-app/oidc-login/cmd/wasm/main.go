package main

import (
	"context"
	"flag"
	"os"

	WASMRuntime "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/WASMRuntime"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

var appContext context.Context

func main() {

	ctx := context.Background()
	logz := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	appContext = logz.WithContext(ctx)

	//prefix := flag.String("prefix", "", "Prefix for the application")
	flag.Parse()
	WASMRuntime.NewWASMApp(appContext, false)

	app.RunWhenOnBrowser()

}
