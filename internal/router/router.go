package router

import (
	"net/http"

	"restapi/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(item *handlers.ItemHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handlers.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/items", func(r chi.Router) {
			r.Get("/", item.List)
			r.Post("/", item.Create)
			r.Get("/{id}", item.Get)
			r.Patch("/{id}", item.Update)
			r.Delete("/{id}", item.Delete)
		})
	})

	return r
}
