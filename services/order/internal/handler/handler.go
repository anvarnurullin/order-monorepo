package handler

import (
	"encoding/json"
	"net/http"
	"order-monorepo/services/order/internal/catalog"
	"order-monorepo/services/order/internal/model"
	"order-monorepo/services/order/internal/store"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store 			*store.Store
	catalogClient 	*catalog.Client
}

func NewHandler(s *store.Store, c *catalog.Client) *Handler {
	return &Handler{
		store: 			s,
		catalogClient: 	c,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o model.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.store.CreateOrderWithStockCheck(r.Context(), o, *h.catalogClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o.ID = id
	o.Status = "pending"
	o.CreatedAt = time.Now()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orders, err := h.store.GetOrders(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateOrderStatus(r.Context(), orderID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
