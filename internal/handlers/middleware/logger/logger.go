package logger

import (
	"log/slog"
	"net/http"
)

func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received:", "request method", r.Method)
		next.ServeHTTP(w, r)
	})
}
