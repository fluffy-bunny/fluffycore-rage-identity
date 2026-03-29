package audits

import (
	"bufio"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		config *contracts_config.Config
	}
)

var stemService = (*service)(nil)

var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container, config),
		config:      config,
	}, nil
}

func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.HTMXSupportAuditsPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func isHTMXRequest(c *echo.Context) bool {
	return (*c).Request().Header.Get("HX-Request") == "true"
}

func resolveAuditFilePath() string {
	candidates := []string{
		filepath.Join("tmp", "auditstore.jsonl"),
		filepath.Join("cmd", "server", "tmp", "auditstore.jsonl"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return candidates[0]
}

func readFilteredLines(path string, contains string, limit int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contains = strings.ToLower(strings.TrimSpace(contains))
	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if contains != "" && !strings.Contains(strings.ToLower(line), contains) {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) > limit {
		lines = lines[len(lines)-limit:]
	}
	return lines, nil
}

func (s *service) Do(c *echo.Context) error {
	if !isHTMXRequest(c) {
		return c.Redirect(http.StatusFound, wellknown_echo.HTMXSupportPath+"?redirect="+(*c).Request().URL.Path)
	}
	contains := (*c).QueryParam("contains")
	limit := 150
	path := resolveAuditFilePath()
	lines, err := readFilteredLines(path, contains, limit)
	if err != nil {
		return c.HTML(http.StatusOK, fmt.Sprintf(`<h2>Audit Logs</h2><p style="color:#ef4444;">Failed to read %s: %s</p>`, html.EscapeString(path), html.EscapeString(err.Error())))
	}

	b := strings.Builder{}
	b.WriteString(`<h2>Audit Logs</h2>`)
	b.WriteString(fmt.Sprintf(`<p>Source file: <code>%s</code></p>`, html.EscapeString(path)))
	b.WriteString(fmt.Sprintf(`<form method="get" action="%s" hx-get="%s" hx-target="#content" hx-push-url="true">`, wellknown_echo.HTMXSupportAuditsPath, wellknown_echo.HTMXSupportAuditsPath))
	b.WriteString(`<label>Contains:</label><input name="contains" type="text" style="margin-left:8px;width:280px;" />`)
	b.WriteString(`<button type="submit" style="margin-left:8px;">Filter</button></form>`)
	b.WriteString(`<p style="margin-top:12px;">Showing latest matching entries.</p><pre style="white-space:pre-wrap;max-height:520px;overflow:auto;background:#0b1220;border:1px solid #334155;padding:12px;border-radius:8px;">`)
	for _, line := range lines {
		b.WriteString(html.EscapeString(line))
		b.WriteString("\n")
	}
	b.WriteString(`</pre>`)
	return c.HTML(http.StatusOK, b.String())
}
