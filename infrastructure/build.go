package infrastructure

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	usecaseimpl "github.com/YagoSchramm/Golinkr/domain/usecase/impl"
	repoimpl "github.com/YagoSchramm/Golinkr/infrastructure/datastore/repository/impl"
	approuter "github.com/YagoSchramm/Golinkr/infrastructure/router"
	"github.com/YagoSchramm/Golinkr/infrastructure/router/modules"
	routerutil "github.com/YagoSchramm/Golinkr/infrastructure/router/util"
	"github.com/YagoSchramm/Golinkr/infrastructure/service/db"
	"github.com/gorilla/mux"
)

func BuildServer() (*http.Server, func(), error) {
	handler, cleanup, err := buildApp()
	if err != nil {
		return nil, cleanup, err
	}

	port := intFromEnv("SERVER_PORT", 8080)
	srv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	slog.Info("Starting server", slog.Int("port", port))
	return srv, cleanup, nil
}

func buildApp() (http.Handler, func(), error) {
	loadDotEnv()

	dsn := os.Getenv("DATABASE_URL")
	secret := os.Getenv("JWT_SECRET")

	if dsn == "" {
		return nil, func() {}, errors.New("DATABASE_URL is not set")
	}
	if secret == "" {
		return nil, func() {}, errors.New("JWT_SECRET is not set")
	}

	dbConn, err := db.NewPostgresConnection(dsn)
	if err != nil {
		return nil, func() {}, err
	}

	userRepository := repoimpl.NewUserRepository(dbConn)
	linkRepository := repoimpl.NewLinkRepository(dbConn)
	analyticsRepository := repoimpl.NewAnalyticsRepository(dbConn)
	cleanup := func() {
		_ = dbConn.Close()
	}

	authUseCase := usecaseimpl.NewAuthUsecase(userRepository, secret)
	linkUseCase := usecaseimpl.NewLinkUsecase(linkRepository)
	analyticsUseCase := usecaseimpl.NewAnalyticsUseCase(analyticsRepository)

	authModule := modules.NewAuthModule(authUseCase, secret)
	linkModule := modules.NewLinkModule(linkUseCase, analyticsUseCase)
	analytcsModule := modules.NewAnalyticsModule(analyticsUseCase)
	healthModule := modules.NewHealthModule(dbConn)

	router := mux.NewRouter()
	router.Use(routerutil.LoggingMiddleware)
	approuter.Mount(
		router,
		authModule.Middlewares(),
		authModule,
		linkModule,
		analytcsModule,
		healthModule,
	)

	handler := withCORS(
		withRateLimit(router, intFromEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", 60)),
		stringSliceFromEnv("CORS_ORIGINS", []string{"*"}),
	)

	return handler, cleanup, nil
}

func withCORS(next http.Handler, allowedOrigins []string) http.Handler {
	allowedMethods := "GET, POST, DELETE, UPDATE"
	allowedHeaders := "Authorization, Content-Type, Accept"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isOriginAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || strings.EqualFold(origin, allowedOrigin) {
			return true
		}
	}

	return false
}

type rateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	requests map[string]rateLimitEntry
}

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

func withRateLimit(next http.Handler, requestsPerMinute int) http.Handler {
	if requestsPerMinute <= 0 {
		return next
	}

	limiter := &rateLimiter{
		limit:    requestsPerMinute,
		window:   time.Minute,
		requests: make(map[string]rateLimitEntry),
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryAfter, allowed := limiter.allow(clientIP(r), time.Now())
		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code":"RATE_LIMIT_EXCEEDED","message":"too many requests"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (l *rateLimiter) allow(key string, now time.Time) (time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := l.requests[key]
	if entry.resetTime.IsZero() || now.After(entry.resetTime) {
		l.requests[key] = rateLimitEntry{
			count:     1,
			resetTime: now.Add(l.window),
		}
		l.cleanupExpired(now)
		return 0, true
	}

	if entry.count >= l.limit {
		return time.Until(entry.resetTime).Round(time.Second), false
	}

	entry.count++
	l.requests[key] = entry
	return 0, true
}

func (l *rateLimiter) cleanupExpired(now time.Time) {
	for key, entry := range l.requests {
		if now.After(entry.resetTime) {
			delete(l.requests, key)
		}
	}
}

func clientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ip, _, _ := strings.Cut(forwardedFor, ",")
		return strings.TrimSpace(ip)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func stringSliceFromEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}

func intFromEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func loadDotEnv() {
	content, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)

		if key == "" {
			continue
		}

		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
