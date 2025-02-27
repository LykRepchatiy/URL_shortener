package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"url_shortener/internal/handlers"
	"url_shortener/internal/middleware/logger"
	"url_shortener/internal/middleware/validate"

	"github.com/go-chi/chi/v5"
)

var (
	postgre, cache bool
)

func init() {
	flag.BoolVar(&postgre, "postgre", false, "use postgre for store URLs")
	flag.BoolVar(&cache, "cache", false, "use cache for store URLs")
	flag.Parse()
	if !postgre && !cache {
		postgre = true
	}
	if cache && postgre {
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		logger.Error("Both postgre and cache flags are set. Please choose only one.")
		os.Exit(1)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(logger.MiddlewareLogger)
	if postgre {
		r.With(validate.MiddlewareValidatePost).Post("/post", handlers.Post)
		r.With(validate.MiddlewareValidateGet).Get("/get", handlers.Get)
	} else if cache {
		log.Println("Not implemented yet")
	}
	http.ListenAndServe(":8080", r)
}
