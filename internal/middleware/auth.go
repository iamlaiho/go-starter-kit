package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/iamlaiho/go-starter-kit/internal/handler"
)

const userIDKey contextKey = "userID"

// Authenticate validates a Bearer JWT in the Authorization header.
// On success it stores the subject (user ID) in the request context.
func Authenticate(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				handler.WriteError(w, http.StatusUnauthorized, "missing or invalid authorization header", GetRequestID(r.Context()))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			}, jwt.WithValidMethods([]string{"HS256"}))
			if err != nil || !token.Valid {
				handler.WriteError(w, http.StatusUnauthorized, "invalid or expired token", GetRequestID(r.Context()))
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				handler.WriteError(w, http.StatusUnauthorized, "invalid token claims", GetRequestID(r.Context()))
				return
			}

			sub, err := claims.GetSubject()
			if err != nil || sub == "" {
				handler.WriteError(w, http.StatusUnauthorized, "invalid token subject", GetRequestID(r.Context()))
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
