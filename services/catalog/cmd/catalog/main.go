package main

import (
	"fmt"
	"net/http"
	"order-monorepo/services/catalog/internal/handler"
	"order-monorepo/services/catalog/internal/logger"
	"order-monorepo/services/catalog/internal/store"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()

	_ = godotenv.Load("../../.env")
	port := os.Getenv("CATALOG_HTTP_PORT")
	if port == "" {
		port = "8082"
	}

	s, err := store.NewStore()
	if err != nil {
		panic(fmt.Errorf("failed to init store: %w", err))
	}

	h := handler.NewHandler(s)

	r := chi.NewRouter()
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
		logger.Info("Health endpoint checked")
	})
	r.Get("/api/v1/products", h.GetProducts)
	r.Get("/api/v1/products/{id}", h.GetProduct)

	r.Patch("/api/v1/products/{id}/decrease", h.DecreaseProductQty)

	addr := fmt.Sprintf(":%s", port)
	logger.Infof("Starting catalog service on port %s", port)
	if err := http.ListenAndServe(addr, r); r != nil {
		logger.Error("failed to start server: ", err)
	}
}
