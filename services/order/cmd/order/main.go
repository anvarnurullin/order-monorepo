package main

import (
	"fmt"
	"net/http"
	"order-monorepo/services/order/internal/catalog"
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

	catalogURL := os.Getenv("CATALOG_URL")
	if catalogURL == "" {
		catalogURL = "localhost"
	}

	catalogPort := os.Getenv("CATALOG_PORT")
	if catalogPort == "" {
		catalogPort = "8082"
	}

	c := catalog.NewClient("http://" + catalogURL + ":" + catalogPort)

	h := handler.NewHandler(s, c)

	r := chi.NewRouter()

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
		fmt.Println("Health endpoint checked")
	})
	r.Get("/api/v1/orders", h.GetOrders)

	r.Post("/api/v1/orders", h.CreateOrder)

	r.Patch("/api/v1/orders/{id}/status", h.UpdateOrderStatus)

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Starting catalog service on port", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
