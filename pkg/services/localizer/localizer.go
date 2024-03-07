package localizer

import (
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_localizer "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/localizer"
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type (
	service struct {
		bundle    contracts_localizer.ILocalizerBundle
		localizer *i18n.Localizer
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_localizer.ILocalizer = stemService
}
func (s *service) Ctor(bundle contracts_localizer.ILocalizerBundle) (contracts_localizer.ILocalizer, error) {
	return &service{
		bundle: bundle,
	}, nil
}

func AddScopedILocalizer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_localizer.ILocalizer](cb, stemService.Ctor)
}
func (s *service) GetLocalizer() *i18n.Localizer {
	return s.localizer
}
func (s *service) Initialize(c echo.Context) {
	accept := c.Request().Header.Get("Accept-Language")
	s.localizer = i18n.NewLocalizer(s.bundle.GetBundle(), accept)
}
