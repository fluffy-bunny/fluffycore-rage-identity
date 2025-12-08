package cache_busting_static_html

import (
	"net/http"
	"os"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

type (
	service struct {
		config           *contracts_config.CacheBustingHTMLConfig
		modifiedContent  string
		staticMiddleware echo.MiddlewareFunc
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder,
	config *contracts_config.CacheBustingHTMLConfig) {

	// load the index.html file and cache bust it

	staticMiddleware := middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   config.StaticPath,
		HTML5:  false, // Disable HTML5 mode so we can handle SPA fallback ourselves with cache-busted content
		Browse: false,
	})

	indexContent, err := os.ReadFile(config.FilePath)
	if err != nil {
		panic(err)

	}
	modifiedContent := string(indexContent)
	for _, kv := range config.ReplaceParams {
		if kv == nil ||
			fluffycore_utils.IsEmptyOrNil(kv.Key) ||
			fluffycore_utils.IsEmptyOrNil(kv.Value) {
			continue
		}
		modifiedContent = strings.ReplaceAll(modifiedContent, kv.Key, kv.Value)
	}
	/*
		// Generate a unique GUID for cache busting
		guid := xid.New().String()
		if pkg_version.Version() != "dev-build" {
			guid = pkg_version.Version()
		}
		// Convert the content to a string
		contentStr := string(indexContent)

		// Replace all instances of {version} with "guid"
		modifiedContent := strings.ReplaceAll(contentStr, "{version}", guid)
	*/
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		func() (*service, error) {
			return &service{
				config:           config,
				modifiedContent:  modifiedContent,
				staticMiddleware: staticMiddleware,
			}, nil
		},
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
		},
		config.EchoPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {

	path := c.Request().URL.Path

	// Serve index.html for root path requests
	if path == s.config.RootPath {
		return c.HTML(http.StatusOK, s.modifiedContent)
	}

	// Strip the RootPath prefix for static file serving
	// e.g., /management/web/app.wasm -> /web/app.wasm
	if strings.HasPrefix(path, s.config.RootPath) {
		strippedPath := strings.TrimPrefix(path, strings.TrimSuffix(s.config.RootPath, "/"))
		c.Request().URL.Path = strippedPath
	}

	// Try to serve as a static file first
	err := s.staticMiddleware(func(c echo.Context) error {
		// return an echo.HTTPError status not found to indicate file not found
		return echo.NewHTTPError(http.StatusNotFound)
	})(c)

	if err == nil {
		// File was successfully served
		return nil
	}

	// Check if this is a 404 error from static middleware
	httpErr, ok := err.(*echo.HTTPError)
	if ok && httpErr.Code == http.StatusNotFound {
		// SPA fallback: serve index.html for 404s to support client-side routing
		// This allows routes like /management/profile to load the WASM app
		return c.HTML(http.StatusOK, s.modifiedContent)
	}

	// Other errors, return as-is
	return err
}
