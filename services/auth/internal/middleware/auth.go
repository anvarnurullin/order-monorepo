package middleware

import (
	"context"
	"net/http"
	"order-monorepo/services/auth/internal/jwt"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserContext struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			claims, err := jwt.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userCtx := UserContext{
				UserID: claims.UserID,
				Email:  claims.Email,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(UserContext)
	return &user, ok
}