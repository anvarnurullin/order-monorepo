package handler

import (
	"encoding/json"
	"net/http"
	"order-monorepo/services/catalog/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	product, err := h.store.GetProductByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
	}

	http.Error(w, "product not found", http.StatusNotFound)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.store.GetProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *Handler) DecreaseProductQty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.DecreaseProductQty(r.Context(), id, req.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
