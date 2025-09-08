package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"todo-app/internal/usecase"
)

// Context key types to avoid collisions
type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UsernameKey contextKey = "username"
)

type AuthMiddleware struct {
	UserInteractor *usecase.UserInteractor
}

func NewAuthMiddleware(userInteractor *usecase.UserInteractor) *AuthMiddleware {
	return &AuthMiddleware{
		UserInteractor: userInteractor,
	}
}

func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Try to get token from cookie first
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// Fallback to Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				log.Printf("No auth token found in cookie or header")
				http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
				return
			}
		} else {
			token = cookie.Value
		}

		// Validate JWT token
		claims, err := am.UserInteractor.ValidateJWTToken(token)
		if err != nil {
			http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Extract user ID from claims
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			http.Error(w, `{"error":"Invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		userID := int(userIDFloat)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, (*claims)["username"])

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Try to get token from cookie first
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// Fallback to Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				// No token, continue without auth
				next.ServeHTTP(w, r)
				return
			}
		} else {
			token = cookie.Value
		}

		// Validate JWT token
		claims, err := am.UserInteractor.ValidateJWTToken(token)
		if err != nil {
			// Invalid token, continue without auth
			next.ServeHTTP(w, r)
			return
		}

		// Extract user ID from claims
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			// Invalid claims, continue without auth
			next.ServeHTTP(w, r)
			return
		}

		userID := int(userIDFloat)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, (*claims)["username"])

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromRequest(r *http.Request) (int, bool) {
	userID := r.Context().Value("userID")
	if userID == nil {
		return 0, false
	}

	switch v := userID.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			return id, true
		}
	}

	return 0, false
}

func GetUsernameFromRequest(r *http.Request) (string, bool) {
	username := r.Context().Value("username")
	if username == nil {
		return "", false
	}

	if usernameStr, ok := username.(string); ok {
		return usernameStr, true
	}

	return "", false
}
