package middleware

import (
	"log/slog"
	"net/http"

	"github.com/iamlaiho/go-starter-kit/internal/handler"
)

func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						"error", err,
						"request_id", GetRequestID(r.Context()),
					)
					handler.WriteError(w, http.StatusInternalServerError, "internal server error", GetRequestID(r.Context()))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
