package logger

import (
	"log"
	"log/slog"
	"net/http"
)

func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("1")
		slog.Info("Request received:", "request method", r.Method)
		next.ServeHTTP(w, r)
	})
}
