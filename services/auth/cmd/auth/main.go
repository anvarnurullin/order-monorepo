package main

import (
	"fmt"
	"net/http"
	"order-monorepo/services/auth/internal/config"
	"order-monorepo/services/auth/internal/handler"
	"order-monorepo/services/auth/internal/logger"
	"order-monorepo/services/auth/internal/store"
	authmiddleware "order-monorepo/services/auth/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.Init()

	cfg := config.Load()

	s, err := store.NewStore(cfg.DBURL)
	if err != nil {
		panic(fmt.Errorf("failed to init store: %w", err))
	}

	h := handler.NewHandler(s, cfg.JWTSecret)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
		logger.Info("Health endpoint checked")
	})
	
	r.Post("/api/v1/auth/register", h.Register)
	r.Post("/api/v1/auth/login", h.Login)
	r.Get("/api/v1/auth/validate", h.ValidateToken)

	r.Route("/api/v1/auth/protected", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/profile", h.GetProfile)
	})

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	logger.Infof("Starting auth service on port %s", cfg.HTTPPort)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Error("failed to start server: ", err)
	}
}
