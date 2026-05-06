package modules

import (
	"context"
	"net/http"
	"strings"

	"github.com/YagoSchramm/Golinkr/domain/entity/derr"
	"github.com/YagoSchramm/Golinkr/domain/usecase"
	"github.com/YagoSchramm/Golinkr/infrastructure/router"
	service "github.com/YagoSchramm/Golinkr/infrastructure/service/jwt"
	"github.com/gorilla/mux"
)

func NewAuthModule(authUseCase usecase.AuthUseCase, secret string) router.Module {
	return &authModule{
		authUseCase: authUseCase,
		name:        "Auth",
		path:        "/auth",
		secret:      secret,
	}
}

type authModule struct {
	authUseCase usecase.AuthUseCase
	name        string
	path        string
	secret      string
}

func (m authModule) Name() string {
	return m.name
}

func (m authModule) Path() string {
	return m.path
}

func (m authModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{}
}

func (m authModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		m.sessionMiddleware(),
	}
}

func (m authModule) sessionMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			scheme, token, ok := strings.Cut(authHeader, " ")
			if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
				w.Header().Set("WWW-Authenticate", "Bearer")
				router.HandleError(w, derr.UnauthorizedError)
				return
			}

			claims, err := service.ValidateToken(token, m.secret)
			if err != nil {
				router.HandleError(w, derr.UnauthorizedError)
				return
			}

			err = m.authUseCase.ValidateSession(r.Context(), claims.UserID, claims.Email)
			if err != nil {
				router.HandleError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), "user_claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
