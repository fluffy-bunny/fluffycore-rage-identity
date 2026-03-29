package logout

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
		wellknownCookies contracts_cookies.IWellknownCookies
	}
)

var stemService = (*service)(nil)
var _ contracts_handler.IHandler = stemService

func (s *service) Ctor(
	container di.Container,
	wellknownCookies contracts_cookies.IWellknownCookies,
	config *contracts_config.Config,
) (*service, error) {
	return &service{
		BaseHandler:      services_echo_handlers_base.NewBaseHandler(container, config),
		wellknownCookies: wellknownCookies,
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.LogoutPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

type LogoutGetRequest struct {
	RedirectURL string `param:"redirect_url" query:"redirect_url" form:"redirect_url" json:"redirect_url" xml:"redirect_url"`
}

func (s *service) validateLoginGetRequest(model *LogoutGetRequest) error {
	if fluffycore_utils.IsEmptyOrNil(model.RedirectURL) {
		model.RedirectURL = "/"
	}
	return nil
}

func (s *service) DoGet(c *echo.Context) error {
	r := c.Request()
	// is the request get or post?

	ctx := r.Context()
	log := zerolog.Ctx(ctx).With().Logger()
	model := &LogoutGetRequest{}
	if err := c.Bind(model); err != nil {
		log.Error().Err(err).Msg("c.Bind")
		return c.Redirect(http.StatusFound, "/error")
	}
	log.Debug().Interface("model", model).Msg("model")
	err := s.validateLoginGetRequest(model)
	if err != nil {
		log.Error().Err(err).Msg("validateLoginGetRequest")
		return c.Redirect(http.StatusFound, "/error")
	}
	s.wellknownCookies.DeleteAuthCookie(c)

	// Render a styled interstitial page with spinner, then redirect after 1 second
	html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Signing Out...</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#f0f2f5;display:flex;align-items:center;justify-content:center;min-height:100vh;color:#24292f}
.logout-card{background:#fff;border-radius:12px;box-shadow:0 2px 12px rgba(0,0,0,0.1);padding:48px 40px;text-align:center;max-width:400px;width:90%}
.spinner{width:40px;height:40px;border:3px solid #e5e7eb;border-top-color:#1a365d;border-radius:50%;animation:spin 0.8s linear infinite;margin:0 auto 24px}
@keyframes spin{to{transform:rotate(360deg)}}
h1{font-size:20px;font-weight:600;margin-bottom:8px}
p{font-size:14px;color:#57606a}
</style>
</head>
<body>
<div class="logout-card">
<div class="spinner"></div>
<h1>Signing out&hellip;</h1>
<p>You will be redirected shortly.</p>
</div>
<script>setTimeout(function(){window.location.href="` + model.RedirectURL + `"},1000);</script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

// HealthCheck godoc
// @Summary get the home page.
// @Description get the home page.
// @Tags root
// @Accept */*
// @Produce json
// @Param       redirect_url            		query     string  true  "redirect url"
// @Success 200 {object} string
// @Router /logout [get,post]
func (s *service) Do(c *echo.Context) error {

	r := c.Request()
	// is the request get or post?
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)

	}
	// return not found
	return c.NoContent(http.StatusNotFound)
}
