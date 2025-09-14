package main

import (
	"fmt"
	"net/http"
	"order-monorepo/services/catalog/internal/handler"
	"order-monorepo/services/catalog/internal/store"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
		fmt.Println("Health endpoint checked")
	})
	r.Get("/api/v1/products", h.GetProducts)

	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Starting catalog service on port", port)
	if err := http.ListenAndServe(addr, r); r != nil {
		fmt.Println("failed to start server: ", err)
	}
}
