// Package middleware Description: This package contains the middleware for the ip2country service.
package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"ip2country/internal/config"
)

var (
	tokens        int
	lastTokenTime = time.Now()
	mu            sync.Mutex
	once          sync.Once
)

type errorResponse struct {
	Error string `json:"error"`
}

func RateLimitMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	once.Do(func() {
		tokens = cfg.BurstLimit
		slog.Warn("Tokens initialized")
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		elapsed := now.Sub(lastTokenTime).Seconds()
		tokens += int(elapsed * float64(cfg.RateLimit))
		if tokens > cfg.BurstLimit {
			tokens = cfg.BurstLimit
		}
		lastTokenTime = now

		if tokens > 0 {
			tokens--
			slog.Info(fmt.Sprintf("Rate Limiting tokens remaining: %d", tokens))
			next.ServeHTTP(w, r)
		} else {
			WriteError(w, http.StatusTooManyRequests, "Too Many Requests")
		}
	})
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorResponse{Error: "Internal Server Error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func WriteError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: errMsg})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info(fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))

		// Create a response writer to capture the status code
		rr := &responseRecorder{w, http.StatusOK}
		next.ServeHTTP(rr, r)

		slog.Info(fmt.Sprintf("Completed %s %s in %v with status %d", r.Method, r.URL.Path, time.Since(start), rr.statusCode))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
