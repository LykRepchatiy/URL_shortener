package handlers

import (
	"context"
	"net/http"
	"url_shortener/internal/handlers/middleware/logger"
	"url_shortener/internal/handlers/middleware/validate"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Router struct {
	PG     *pgx.Conn
	Router chi.Mux
}

func NewRouter() *Router {
	return &Router{
		Router: *chi.NewRouter(),
	}
}

func (r *Router) StartDB(DBConn *pgx.Conn) error {
	r.PG = DBConn
	r.Router.Use(logger.MiddlewareLogger)
	r.Router.With(validate.MiddlewareValidatePost).Post("/post", r.PostDB)
	r.Router.With(validate.MiddlewareValidateGet).Get("/get", r.GetDB)
	return http.ListenAndServe(":8080", &r.Router)
}

func (r *Router) StartCache() error {
	r.Router.Use(logger.MiddlewareLogger)
	// r.Router.With(validate.MiddlewareValidatePost).Post("/post" /*Post*/)
	// r.Router.With(validate.MiddlewareValidateGet).Get("/get" /*Get*/)
	// return http.ListenAndServe(":8080", &r.Router)
	return nil
}

func (r *Router) Finish() {
	if r.PG != nil {
		r.PG.Close(context.Background())
	}
}
