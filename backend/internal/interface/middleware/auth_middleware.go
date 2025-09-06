package middleware

import (
	"context"
	"net/http"
	"strconv"
	"todo-app/internal/usecase"
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
		// Get token from cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Validate JWT token
		claims, err := am.UserInteractor.ValidateJWTToken(cookie.Value)
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
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "username", (*claims)["username"])

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie (optional)
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// No token, continue without auth
			next.ServeHTTP(w, r)
			return
		}

		// Validate JWT token
		claims, err := am.UserInteractor.ValidateJWTToken(cookie.Value)
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
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "username", (*claims)["username"])

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
