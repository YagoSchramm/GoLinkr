package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	approuter "github.com/YagoSchramm/Golinkr/infrastructure/router"
	"github.com/gorilla/mux"
)

type healthModule struct {
	db   *sql.DB
	name string
	path string
}

type healthResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Database string `json:"database,omitempty"`
}

func NewHealthModule(db *sql.DB) approuter.Module {
	return healthModule{
		db:   db,
		name: "Health",
		path: "/health",
	}
}

func (m healthModule) Name() string {
	return m.name
}

func (m healthModule) Path() string {
	return m.path
}

func (m healthModule) Routes() []approuter.RouteDefinition {
	return []approuter.RouteDefinition{
		{
			Path:        "",
			Description: "Liveness check",
			Handler:     m.live,
			HttpMethods: []string{http.MethodGet},
			Public:      true,
		},
		{
			Path:        "/db",
			Description: "Database readiness check",
			Handler:     m.ready,
			HttpMethods: []string{http.MethodGet},
			Public:      true,
		},
	}
}

func (m healthModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (m healthModule) live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = approuter.Write(w, healthResponse{
		Code:    "SUCCESS",
		Message: "service is healthy",
	})
}

func (m healthModule) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := m.db.PingContext(ctx); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthResponse{
			Code:     "UNAVAILABLE",
			Message:  "database is not reachable",
			Database: "down",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = approuter.Write(w, healthResponse{
		Code:     "SUCCESS",
		Message:  "service is ready",
		Database: "up",
	})
}
