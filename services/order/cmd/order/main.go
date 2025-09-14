package main

import (
	"fmt"
	"net/http"
	"order-monorepo/services/order/internal/handler"
	"order-monorepo/services/order/internal/store"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	port := os.Getenv("ORDER_HTTP_PORT")
	if port == "" {
		port = "8083"
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
		fmt.Println("Health endpoint checked")
	})

	r.Post("/api/v1/orders", h.CreateOrder)
	r.Get("/api/v1/orders", h.GetOrders)

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Starting catalog service on port", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
