package prefs

import (
	"net/http"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	components "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/htmx/components"
	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	contracts_cookies "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/cookies"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/wellknown/wellknown_echo"
	proto_oidc_models "github.com/fluffy-bunny/fluffycore-rage-identity/proto/oidc/models"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v5"
	zerolog "github.com/rs/zerolog"
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
			contracts_handler.POST,
		},
		wellknown_echo.HTMXManagementPrefsPath,
	)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c *echo.Context) error {
	// Non-HTMX GET requests (e.g. browser refresh) need the full shell page
	if c.Request().Method == http.MethodGet && !components.IsHTMXRequest(c) {
		return c.Redirect(http.StatusFound, wellknown_echo.HTMXManagementPath+"?redirect="+c.Request().URL.Path)
	}
	r := c.Request()
	switch r.Method {
	case http.MethodGet:
		return s.DoGet(c)
	case http.MethodPost:
		return s.DoPost(c)
	}
	return c.NoContent(http.StatusNotFound)
}

func (s *service) getSubject(c *echo.Context) string {
	memCache := s.ScopedMemoryCache()
	cachedItem, ok := memCache.Get("rootIdentity")
	if !ok {
		return ""
	}
	rootIdentity, ok := cachedItem.(*proto_oidc_models.Identity)
	if !ok || rootIdentity == nil {
		return ""
	}
	return rootIdentity.Subject
}

func (s *service) getPrefsData(c *echo.Context) *components.PreferencesPageData {
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	subject := s.getSubject(c)
	if subject == "" {
		return &components.PreferencesPageData{
			RenderContext: rc,
			Error:         rc.L("mgmt_unexpected_error"),
		}
	}

	data := &components.PreferencesPageData{
		RenderContext: rc,
	}

	getResp, err := s.WellknownCookies().GetKeepSigninPreferencesCookie(c,
		&contracts_cookies.GetKeepSigninPreferencesCookieRequest{
			Subject: subject,
		})
	if err == nil && getResp != nil {
		data.KeepSignedIn = getResp.KeepSigninPreferencesCookie.KeepSignedIn
		data.DontShowAgain = getResp.KeepSigninPreferencesCookie.DoNotAskAgain
	}

	return data
}

func (s *service) DoGet(c *echo.Context) error {
	data := s.getPrefsData(c)
	return components.RenderNode(c, http.StatusOK, components.PreferencesPage(data))
}

func (s *service) DoPost(c *echo.Context) error {
	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	localizer := s.Localizer().GetLocalizer()
	rc := components.NewRenderContext(c, localizer)

	action := c.FormValue("action")

	subject := s.getSubject(c)
	if subject == "" {
		return components.RenderNode(c, http.StatusOK, components.PreferencesPage(&components.PreferencesPageData{
			RenderContext: rc,
			Error:         rc.L("mgmt_unexpected_error"),
		}))
	}

	switch action {
	case "save-prefs":
		keepSignedIn := c.FormValue("keepSignedIn") == "on"
		dontShowAgain := c.FormValue("dontShowAgain") == "on"

		err := s.WellknownCookies().SetKeepSigninPreferencesCookie(c,
			&contracts_cookies.SetKeepSigninPreferencesCookieRequest{
				Subject: subject,
				KeepSigninPreferencesCookie: &contracts_cookies.KeepSigninPreferencesCookie{
					KeepSignedIn:  keepSignedIn,
					DoNotAskAgain: dontShowAgain,
				},
			})
		if err != nil {
			log.Error().Err(err).Msg("SetKeepSigninPreferencesCookie")
			data := s.getPrefsData(c)
			data.Error = rc.L("mgmt_something_went_wrong")
			return components.RenderNode(c, http.StatusOK, components.PreferencesPage(data))
		}

		data := s.getPrefsData(c)
		data.Success = rc.L("mgmt_success")
		return components.RenderNode(c, http.StatusOK, components.PreferencesPage(data))

	case "clear-sso":
		s.WellknownCookies().DeleteSSOCookie(c)

		data := s.getPrefsData(c)
		data.Success = rc.L("mgmt_success")
		return components.RenderNode(c, http.StatusOK, components.PreferencesPage(data))
	}

	data := s.getPrefsData(c)
	return components.RenderNode(c, http.StatusOK, components.PreferencesPage(data))
}
