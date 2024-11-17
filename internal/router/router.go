// Package router Description: This package contains the router for ip2country service.
package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"ip2country/internal/config"
	"ip2country/internal/ip2country/handler"
	"ip2country/internal/middleware"
)

func NewRouter(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ErrorHandler)
	r.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimitMiddleware(cfg, next)
	})
	r.HandleFunc("/v1/find-country", handler.FindCountryHandler).Methods("GET")
	return r
}

func StartServer(cfg *config.Config) {
	r := NewRouter(cfg)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}
	_ = httpServer.ListenAndServe()
}
