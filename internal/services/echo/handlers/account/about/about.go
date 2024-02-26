package about

import (
	"net/http"
	"reflect"
	"strings"

	golinq "github.com/ahmetb/go-linq/v3"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	services_echo_handlers_base "github.com/fluffy-bunny/fluffycore-rage-identity/internal/services/echo/handlers/base"
	wellknown_echo "github.com/fluffy-bunny/fluffycore-rage-identity/internal/wellknown/echo"
	contracts_handler "github.com/fluffy-bunny/fluffycore/echo/contracts/handler"
	echo "github.com/labstack/echo/v4"
	zerolog "github.com/rs/zerolog"
)

type (
	service struct {
		*services_echo_handlers_base.BaseHandler
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_handler.IHandler = stemService
}

func (s *service) Ctor(container di.Container) (*service, error) {
	return &service{
		BaseHandler: services_echo_handlers_base.NewBaseHandler(container),
	}, nil
}

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder di.ContainerBuilder) {
	contracts_handler.AddScopedIHandleWithMetadata[*service](builder,
		stemService.Ctor,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown_echo.AboutPath,
	)

}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

// HealthCheck godoc
// @Summary get the about page.
// @Description get the about page.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /about [get]
func (s *service) Do(c echo.Context) error {

	ctx := c.Request().Context()
	log := zerolog.Ctx(ctx).With().Logger()
	ctn := s.Container
	descriptors := ctn.GetDescriptors()
	log.Info().Msg("about")
	type row struct {
		Verbs string
		Path  string
	}

	var rows []row

	golinq.
		From(descriptors).
		WhereT(func(descriptor *di.Descriptor) bool {
			found := false
			for _, serviceType := range descriptor.ImplementedInterfaceTypes {
				if serviceType == reflect.TypeOf((*contracts_handler.IHandler)(nil)).Elem() {
					found = true
					break
				}
			}
			return found
		}).
		Select(func(c interface{}) interface{} {
			descriptor := c.(*di.Descriptor)
			found := false
			for _, serviceType := range descriptor.ImplementedInterfaceTypes {
				if serviceType == reflect.TypeOf((*contracts_handler.IHandler)(nil)).Elem() {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
			metadata := descriptor.Metadata
			path := metadata["path"].(string)
			httpVerbs, _ := metadata["httpVerbs"].([]contracts_handler.HTTPVERB)
			verbBldr := strings.Builder{}

			for idx, verb := range httpVerbs {
				verbBldr.WriteString(verb.String())
				if idx < len(httpVerbs)-1 {
					verbBldr.WriteString(",")
				}
			}
			return row{
				Verbs: verbBldr.String(),
				Path:  path,
			}

		}).OrderBy(func(i interface{}) interface{} {
		return i.(row).Path
	}).ToSlice(&rows)

	return s.Render(c, http.StatusOK, "account/about/index",
		map[string]interface{}{
			"defs": rows,
		})

}
