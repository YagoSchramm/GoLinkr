package modules

import (
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/router"
	"github.com/gorilla/mux"
)

func NewAnalyticsModule(analyticUseCase usecase.AnalyticsUseCase) router.Module {
	return analyticsModule{
		useCase: analyticUseCase,
		name:    "Analytics",
		path:    "/analytics",
	}
}

type analyticsModule struct {
	useCase usecase.AnalyticsUseCase
	name    string
	path    string
}

func (m analyticsModule) Name() string {
	return m.name
}

func (m analyticsModule) Path() string {
	return m.path
}

func (m analyticsModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{}
}

func (m analyticsModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}
