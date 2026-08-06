package modules

import (
	"log/slog"
	"net/http"

	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/router"
	service "github.com/YagoSchramm/Golinkr/infrastructure/service/jwt"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
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
			Path:        "/{link_id}",
			Description: "Get the analytics from the link",
			Handler:     m.getAnalyticsByLinkId,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
		{
			Path:        "/{link_id}/hourly-click-averages",
			Description: "List hourly click averages from the link",
			Handler:     m.listHourlyClickAverages,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
		{
			Path:        "/{link_id}/monthly-week-click-averages",
			Description: "List monthly week click averages from the link",
			Handler:     m.listMonthlyWeekClickAverages,
			HttpMethods: []string{http.MethodGet},
			Public:      false,
		},
		{
			Path:        "/{link_id}/weekday-click-averages",
			Description: "List weekday click averages from the link",
			Handler:     m.listWeekdayClickAverages,
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

	linkID, err := uuid.Parse(mux.Vars(r)["link_id"])
	if err != nil {
		router.HandleError(w, derr.InvalidLinkId)
		return
	}

	claims := r.Context().Value("user_claims").(*service.Claims)
	response, err := m.useCase.GetByLinkId(ctx, linkID, claims.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get the analytics by link id", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}
}

func (m analyticsModule) listMonthlyWeekClickAverages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	linkID, err := uuid.Parse(mux.Vars(r)["link_id"])
	if err != nil {
		router.HandleError(w, derr.InvalidLinkId)
		return
	}

	claims := r.Context().Value("user_claims").(*service.Claims)
	response, err := m.useCase.ListMonthlyWeekClickAverages(ctx, linkID, claims.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list monthly week click averages", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}
}

func (m analyticsModule) listWeekdayClickAverages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	linkID, err := uuid.Parse(mux.Vars(r)["link_id"])
	if err != nil {
		router.HandleError(w, derr.InvalidLinkId)
		return
	}

	claims := r.Context().Value("user_claims").(*service.Claims)
	response, err := m.useCase.ListWeekdayClickAverages(ctx, linkID, claims.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list weekday click averages", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}
}

func (m analyticsModule) listHourlyClickAverages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	linkID, err := uuid.Parse(mux.Vars(r)["link_id"])
	if err != nil {
		router.HandleError(w, derr.InvalidLinkId)
		return
	}

	claims := r.Context().Value("user_claims").(*service.Claims)
	response, err := m.useCase.ListHourlyClickAverages(ctx, linkID, claims.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list hourly click averages", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "failed to write response", slog.Any("err", err))
		router.HandleError(w, err)
		return
	}
}
