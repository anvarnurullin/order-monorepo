package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-monorepo/services/order/internal/catalog"
	"order-monorepo/services/order/internal/kafka"
	"order-monorepo/services/order/internal/logger"
	"order-monorepo/services/order/internal/model"
	"order-monorepo/services/order/internal/store"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store 			*store.Store
	catalogClient 	*catalog.Client
	kafkaProducer   *kafka.Producer
}

func NewHandler(s *store.Store, c *catalog.Client, k *kafka.Producer) *Handler {
	return &Handler{
		store: 			s,
		catalogClient: 	c,
		kafkaProducer:  k,
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

	if err := h.kafkaProducer.SendMessage(
		"order",
		fmt.Sprintf(`{"id":%d,"product_id":%d,"quantity":%d}`,
		o.ID, o.ProductID, o.Quantity));
		err != nil {
			logger.Error("failed to send kafka message", err)
	}

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

	order, err := h.store.GetOrderByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
	}

	if !model.IsValidStatusTransition(order.Status, req.Status) {
		http.Error(w, fmt.Sprintf("invalid status transition: %s -> %s", order.Status, req.Status), http.StatusBadRequest)
	}

	if err := h.store.UpdateOrderStatus(r.Context(), orderID, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event := fmt.Sprintf(`{
		"event": "order_status_updated",
		"order_id": %d,
		"old_status": "%s",
		"new_status": "%s",
		"updated_at": "%s"
	}`, order.ID, order.Status, req.Status, time.Now().Format(time.RFC3339))

	if err := h.kafkaProducer.SendMessage("order", event); err != nil {
		logger.Error("failed to send kafka message", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
