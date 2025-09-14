package handler

import (
	"encoding/json"
	"net/http"
	"order-monorepo/services/order/internal/model"
	"order-monorepo/services/order/internal/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) CreateOrder (w http.ResponseWriter, r *http.Request) {
	var o model.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.store.CreateOrder(r.Context(), o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

func (h *Handler) GetOrders (w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orders, err := h.store.GetOrders(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
