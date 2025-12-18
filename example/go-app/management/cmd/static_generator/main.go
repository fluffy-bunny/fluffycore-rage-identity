package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	WASMRuntime "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/WASMRuntime"
	backend_ResourceResolvers "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/backend/ResourceResolvers"
	common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/common"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	logz := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	appContext := logz.WithContext(ctx)
	log := zerolog.Ctx(appContext).With().Logger()

	WASMRuntime.NewWASMApp(appContext, true)

	title := flag.String("title", "MyApp", "Title of the app")
	outputPath := flag.String("output", "./static_output", "Output path for the generated static files")
	baseHRef := flag.String("base-href", "/", "Base HRef for the application")
	flag.Parse()
	basePath := *baseHRef
	fmt.Printf("Base HRef: '%s'\n", basePath)

	fixBasePath := func(basePath string) string {
		// ensure that it looks something like "/oidc-login" or "/"
		if fluffycore_utils.IsEmptyOrNil(basePath) {
			return basePath
		}
		if !strings.HasPrefix(basePath, "/") {
			basePath = "/" + basePath
		}
		basePath = strings.TrimSuffix(basePath, "/")
		return basePath
	}
	basePath = fixBasePath(basePath)
	// Generate cache-busting version based on build info
	version := common.AppVersion
	if version == "" || version == "dev" {
		version = fmt.Sprintf("%d", os.Getpid()) // Fallback if ldflags not set
	}
	log.Info().
		Str("version", version).
		Str("buildTime", common.BuildTime).
		Str("gitCommit", common.GitCommit).
		Msg("Build info")

	resourceResolver := backend_ResourceResolvers.ResourceResolverWithBaseHRefResolverOptions(
		backend_ResourceResolvers.BaseHRefResolverOptions{
			Version: version,
			Prefix:  basePath,
		},
	)

	handler := &app.Handler{
		Name:        *title,
		Description: *title,
		Version:     version,

		// RawHeaders allows us to inject custom HTML into the <head> section
		/*
			<meta name="app-version" content="1764940685">
			<meta name="app-build-time" content="2025-12-05_13:17:39">
			<meta name="app-git-commit" content="5cda623">
			<meta name="app-git-branch" content="initial">

		*/
		RawHeaders: []string{
			fmt.Sprintf(`<meta name="app-version" content="%s">`, version),
			fmt.Sprintf(`<meta name="app-build-time" content="%s">`, common.BuildTime),
			fmt.Sprintf(`<meta name="app-git-commit" content="%s">`, common.GitCommit),
			fmt.Sprintf(`<meta name="app-git-branch" content="%s">`, common.GitBranch),
			`<link rel="icon" type="image/svg+xml" href="web/m_logo.svg">`,
			`<link rel="icon" type="image/png" sizes="16x16" href="web/favicon-16x16.png">`,
			`<link rel="icon" type="image/png" sizes="32x32" href="web/favicon-32x32.png">`,
			`<link rel="icon" type="image/png" sizes="64x64" href="web/favicon-64x64.png">`,
			`<link rel="apple-touch-icon" sizes="180x180" href="web/apple-touch-icon.png">`,
			`<link rel="icon" type="image/png" sizes="192x192" href="web/android-chrome-192x192.png">`,
			`<link rel="icon" type="image/png" sizes="512x512" href="web/android-chrome-512x512.png">`,
			`<meta name="msapplication-TileImage" content="web/mstile-150x150.png">`,
		},
		CacheableResources: []string{
			"/web/app.wasm",
		},
		Styles: []string{
			"/web/styles.css",
		},
		Scripts: []string{
			"/web/build_version.js",
			"/web/common.js",
			"/web/logging-helper.js",
			"/web/webauthn.js",
		},
		Resources: resourceResolver,
		Icon: app.Icon{
			SVG:     "web/m_logo.svg",
			Default: "web/m_logo.svg",
			Large:   "web/m_logo.svg",
		},
		//Image: "web/m_logo.svg",
	}
	if fluffycore_utils.IsNotEmptyOrNil(basePath) {
		// Prepend base href to existing RawHeaders
		handler.RawHeaders = append([]string{
			fmt.Sprintf("<base href=\"%s/\">", basePath),
		}, handler.RawHeaders...)
	}
	outputDir := *outputPath
	log.Info().
		Str("outputDir", outputDir).
		Str("basePath", basePath).
		Msg("Generating static website...")

	err := app.GenerateStaticWebsite(outputDir, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate static website")
	}

	log.Info().
		Str("outputDir", outputDir).
		Str("version", version).
		Str("basePath", basePath).
		Msg("✅ Static website generated successfully! Files are ready to be hosted by another web application.")

}
