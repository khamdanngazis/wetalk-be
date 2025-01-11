package middleware

import (
	"chat-be/package/helper"
	"context"
	"log"
	"net/http"
	"strings"
)

// Define a custom type for the context key to avoid collisions
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware checks if the JWT is valid and sets the user context
func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteResponse(w, http.StatusUnauthorized, "Missing Authorization token", nil)
			return
		}

		// Remove the "Bearer " prefix if it exists
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate the token
		claims, err := helper.ValidateToken(tokenString)
		if err != nil {
			log.Println("Invalid token :", err)
			WriteResponse(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		// Set the user information in the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, userContextKey, claims)

		// Pass the request to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GuestMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			ctx = context.WithValue(ctx, userContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := helper.ValidateToken(tokenString)
		if err != nil {
			ctx = context.WithValue(ctx, userContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx = context.WithValue(ctx, userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetUserFromContext(ctx context.Context) (*helper.Claims, bool) {
	user, ok := ctx.Value(userContextKey).(*helper.Claims)
	return user, ok
}
