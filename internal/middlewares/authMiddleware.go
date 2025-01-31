package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/mahabub618/minipack/config"
	"github.com/mahabub618/minipack/internal/handlers"
)

func AuthMiddleware(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		cfg, err := config.LoadConfig()
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to load config", http.StatusInternalServerError)
			return
		}

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Role != requiredRole {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
