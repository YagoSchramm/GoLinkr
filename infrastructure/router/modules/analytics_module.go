package modules

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/YagoSchramm/Golinkr/domain/entity"
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
	return []router.RouteDefinition{
		{
			Path:        "/",
			Description: "Get the analytics from the link",
			Handler:     m.getAnalyticsByLinkId,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
	}
}

func (m analyticsModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (m analyticsModule) getAnalyticsByLinkId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", err) // needs logger
		router.HandleError(w, err)
		return
	}
	var analytics entity.GetAnalyticsDTO
	err = json.Unmarshal(body, &analytics)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", err) // needs logger
		router.HandleError(w, err)
		return
	}
	response, err := m.useCase.GetByLinkId(ctx, analytics.LinkID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get the analytics by link id", err)
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", err) // needs logger
		router.HandleError(w, err)
		return
	}
}
