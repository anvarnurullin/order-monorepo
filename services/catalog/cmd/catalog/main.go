package main

import (
	"context"
	"fmt"
	"net/http"
	"order-monorepo/services/catalog/internal/config"
	"order-monorepo/services/catalog/internal/handler"
	"order-monorepo/services/catalog/internal/logger"
	"order-monorepo/services/catalog/internal/minio"
	"order-monorepo/services/catalog/internal/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger.Init()

	cfg := config.Load()

	s, err := store.NewStore(cfg.DBURL)
	if err != nil {
		panic(fmt.Errorf("failed to init store: %w", err))
	}

	minioClient, err := minio.NewClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket)
	if err != nil {
		panic(fmt.Errorf("failed to init minio client: %w", err))
	}

	if err := minioClient.InitBucket(context.Background()); err != nil {
		panic(fmt.Errorf("failed to init minio bucket: %w", err))
	}

	h := handler.NewHandler(s, minioClient)

	r := chi.NewRouter()
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
		logger.Info("Health endpoint checked")
	})
	r.Get("/api/v1/products", h.GetProducts)
	r.Get("/api/v1/products/{id}", h.GetProduct)
	r.Post("/api/v1/products/{id}/image", h.UploadProductImage)
	r.Get("/api/v1/images/*", h.GetProductImage)

	r.Patch("/api/v1/products/{id}/decrease", h.DecreaseProductQty)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	logger.Infof("Starting catalog service on port %s", cfg.HTTPPort)
	if err := http.ListenAndServe(addr, r); r != nil {
		logger.Error("failed to start server: ", err)
	}
}
