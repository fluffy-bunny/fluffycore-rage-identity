package cache_busting_html

import (
	"net/http"
	"os"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
)

type (
	service struct {
		config          *contracts_config.CacheBustingHTMLConfig
		modifiedContent string
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder, config *contracts_config.CacheBustingHTMLConfig) {

	// load the index.html file and cache bust it

	indexContent, err := os.ReadFile(config.FilePath)
	if err != nil {
		panic(err)

	}

	// Convert the content to a string
	contentStr := string(indexContent)

	// Replace all instances of {version} with "guid"
	modifiedContent := strings.ReplaceAll(contentStr, "{version}", config.Version)

	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		func() (*service, error) {
			return &service{
				config:          config,
				modifiedContent: modifiedContent,
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
	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		return c.HTML(http.StatusOK, s.modifiedContent)
	}
	// return not found
	return c.NoContent(http.StatusNotFound)

}
