package handler

import (
	"encoding/json"
	"net/http"
	"order-monorepo/services/auth/internal/jwt"
	"order-monorepo/services/auth/internal/logger"
	"order-monorepo/services/auth/internal/middleware"
	"order-monorepo/services/auth/internal/model"
	"order-monorepo/services/auth/internal/store"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	store     *store.Store
	jwtSecret string
}

func NewHandler(store *store.Store, jwtSecret string) *Handler {
	return &Handler{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.store.CreateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		logger.Error("Failed to create user:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		logger.Error("Failed to generate token:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := model.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		logger.Error("Failed to get user:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.store.ValidatePassword(user.Password, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		logger.Error("Failed to generate token:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := model.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Bearer token required", http.StatusUnauthorized)
		return
	}

	claims, err := jwt.ValidateToken(tokenString, h.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"user_id": claims.UserID,
		"email":   claims.Email,
	})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	dbUser, err := h.store.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		logger.Error("Failed to get user profile:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbUser)
}