package ResourceResolvers

import (
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type (
	BaseHRefResolverOptions struct {
		Version string
		Prefix  string
	}
)

func ResourceResolverWithBaseHRefResolverOptions(
	options BaseHRefResolverOptions,
) app.ResourceResolver {
	rr := versionedCacheBustingBaseHRefResourceResolver{
		version: options.Version,
		prefix:  options.Prefix,
	}

	return rr
}

type versionedCacheBustingBaseHRefResourceResolver struct {
	app.ResourceResolver
	baseResolver app.ResourceResolver
	version      string
	prefix       string
}

func (r versionedCacheBustingBaseHRefResourceResolver) Resolve(path string) string {
	switch path {

	case
		"/",
		"/web":
		if r.prefix != "" {
			path = r.prefix + path
		}
	case
		"/web/apple-touch-icon.png",
		"/web/app.wasm",
		"/web/styles.css",
		"web/m_logo.svg",
		"/web/build_version.js",
		"/web/common.js",
		"/web/logging-helper.js",
		"/web/webauthn.js",
		"/app.js",
		"/app-worker.js",
		"/wasm_exec.js",
		"/manifest.webmanifest",
		"/app.css":

		path = strings.TrimPrefix(path, "/")
		path = path + "?v=" + r.version

	}

	return path
}
