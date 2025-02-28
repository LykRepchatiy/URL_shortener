package validate

import (
	"net/http"
	"strings"
	"url_shortener/internal/service"
)

func MiddlewareValidatePost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MiddlewareValidateGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
			return
		}
		short_url := r.URL.Query().Get("short_url")
		if len(short_url) != 10 {
			http.Error(w, "Invalid short URL", http.StatusBadRequest)
			return
		}
		for _, char := range short_url {
			if !strings.ContainsRune(service.Alphabet, char) {
				http.Error(w, "Invalid short URL", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
