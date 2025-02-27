package main

import (
	"flag"
	"log"
	"net/http"

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
		// TO DO check for "Fatal" in SLOG
		log.Fatal("Both postgre and cache flags are set. Please choose only one.")
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
