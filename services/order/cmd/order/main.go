package main

import (
	"fmt"
	"net/http"
	"order-monorepo/services/order/internal/catalog"
	"order-monorepo/services/order/internal/config"
	"order-monorepo/services/order/internal/handler"
	"order-monorepo/services/order/internal/kafka"
	"order-monorepo/services/order/internal/logger"
	"order-monorepo/services/order/internal/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger.Init()

	cfg := config.Load()

	s, err := store.NewStore(cfg.DBURL)
	if err != nil {
		panic(fmt.Errorf("failed to init store: %w", err))
	}

	c := catalog.NewClient(cfg.CatalogURL)

	kafkaProducer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	defer kafkaProducer.Close()

	h := handler.NewHandler(s, c, kafkaProducer)

	r := chi.NewRouter()

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
		logger.Info("Health endpoint checked")
	})
	r.Get("/api/v1/orders", h.GetOrders)

	r.Post("/api/v1/orders", h.CreateOrder)

	r.Patch("/api/v1/orders/{id}/status", h.UpdateOrderStatus)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	logger.Info("Starting catalog service on port" + cfg.HTTPPort)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
