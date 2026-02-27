package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ardani17/taskmanager/internal/services"
	"github.com/ardani17/taskmanager/pkg/utils"
)

// Context keys for storing user info
type contextKey string

const (
	UserIDKey    contextKey = "userID"
	UserEmailKey contextKey = "userEmail"
	UserRoleKey  contextKey = "userRole"
)

// AuthMiddleware validates JWT tokens and protects routes
func AuthMiddleware(jwtService *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Extract token
			tokenString, err := services.ExtractToken(authHeader)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			// Validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				if err == services.ErrExpiredToken {
					utils.ErrorResponse(w, http.StatusUnauthorized, "Token has expired")
					return
				}
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add user info to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.DeveloperID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			// Continue to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from request context (returns int)
func GetUserID(r *http.Request) int {
	if userIDStr, ok := r.Context().Value(UserIDKey).(string); ok {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return 0
		}
		return userID
	}
	return 0
}

// GetUserIDString extracts user ID from request context (returns string)
func GetUserIDString(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserEmail extracts user email from request context
func GetUserEmail(r *http.Request) string {
	if email, ok := r.Context().Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

// GetUserRole extracts user role from request context
func GetUserRole(r *http.Request) string {
	if role, ok := r.Context().Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r)
			if role != requiredRole && role != "admin" {
				utils.ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IsAuthenticated checks if request has valid authentication
func IsAuthenticated(r *http.Request) bool {
	return GetUserID(r) > 0
}
