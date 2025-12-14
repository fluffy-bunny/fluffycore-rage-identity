package ResourceResolvers

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type (
	ResourceResolverWithCacheBustingVersionOptions struct {
		Version string
		Prefix  string
	}
)

func ResourceResolverWithCacheBustingVersion(base app.ResourceResolver,
	options ResourceResolverWithCacheBustingVersionOptions,
	paths ...string) app.ResourceResolver {
	pathMap := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if strings.HasPrefix(p, "/web/") {
			pathMap[p] = struct{}{}
		}
	}

	return versionedCacheBustingResourceResolver{
		ResourceResolver: base,
		version:          options.Version,
		prefix:           options.Prefix,
		paths:            pathMap,
	}
}

type versionedCacheBustingResourceResolver struct {
	app.ResourceResolver

	version string
	prefix  string
	paths   map[string]struct{}
}

func (r versionedCacheBustingResourceResolver) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}

	if strings.HasPrefix(path, "/web/") {
		filePath := filepath.FromSlash(strings.TrimPrefix(path, "/"))
		ext := filepath.Ext(filePath)
		if ext != "" {
			switch ext {
			case ".wasm":
				writer.Header().Set("Content-Type", "application/wasm")
			default:
				if contentType := mime.TypeByExtension(ext); contentType != "" {
					writer.Header().Set("Content-Type", contentType)
				}
			}
		}

		http.ServeFile(writer, request, filePath)
		return
	}

	if handler, ok := r.ResourceResolver.(http.Handler); ok {
		handler.ServeHTTP(writer, request)
		return
	}

	http.NotFound(writer, request)
}

func (r versionedCacheBustingResourceResolver) Resolve(path string) string {
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
