package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"order-monorepo/services/catalog/internal/logger"
	"order-monorepo/services/catalog/internal/minio"
	"order-monorepo/services/catalog/internal/store"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store       *store.Store
	minioClient *minio.Client
}

func NewHandler(s *store.Store, mc *minio.Client) *Handler {
	return &Handler{
		store:       s,
		minioClient: mc,
	}
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

func (h *Handler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.store.GetProductByID(r.Context(), id)
	if err != nil || product == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "file must be an image", http.StatusBadRequest)
		return
	}

	objectName := fmt.Sprintf("products/%d/%s", id, header.Filename)
	logger.Infof("Uploading file: %s with content-type: %s\n", objectName, contentType)

	imageURL, err := h.minioClient.UploadFile(r.Context(), objectName, file, header.Size, contentType)
	if err != nil {
		logger.Error("Failed to upload to MinIO: %v\n", err)
		http.Error(w, "failed to upload image", http.StatusInternalServerError)
		return
	}

	logger.Infof("Image uploaded successfully, URL: %s\n", imageURL)

	err = h.store.UpdateProductImage(r.Context(), id, imageURL)
	if err != nil {
		logger.Error("Failed to update database: %v\n", err)
		http.Error(w, "failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"image_url": imageURL})
}

func (h *Handler) GetProductImage(w http.ResponseWriter, r *http.Request) {
	objectName := chi.URLParam(r, "*")

	object, err := h.minioClient.GetObject(r.Context(), objectName)
	if err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		http.Error(w, "failed to get image info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", stat.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=3600")

	_, err = io.Copy(w, object)
	if err != nil {
		http.Error(w, "failed to serve image", http.StatusInternalServerError)
		return
	}
}
