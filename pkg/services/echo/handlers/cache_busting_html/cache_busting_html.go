package cache_busting_html

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	xid "github.com/rs/xid"
)

type (
	Config struct {
		FilePath string `json:"filePath"`
		URIPath  string `json:"uriPath"`
	}
	service struct {
		config          *Config
		modifiedContent string
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder, config *Config) {

	// load the index.html file and cache bust it

	indexContent, err := os.ReadFile(config.FilePath)
	if err != nil {
		panic(err)

	}
	// Generate a unique GUID for cache busting
	guid := xid.New().String()
	// Replace <script src="..."></script> with cache-busted versions
	re := regexp.MustCompile(`<script src="([^"]+)"></script>`)
	modifiedContent := re.ReplaceAllStringFunc(string(indexContent), func(match string) string {
		return strings.Replace(match, `">`, fmt.Sprintf(`?_cacheBust=%s">`, guid), 1)
	})
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		func() (*service, error) {
			return &service{
				config:          config,
				modifiedContent: modifiedContent,
			}, nil
		},
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		config.URIPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	return c.HTML(http.StatusOK, s.modifiedContent)
}
