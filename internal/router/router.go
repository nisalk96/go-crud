package router

import (
	"net/http"
	"path/filepath"

	"restapi/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(mh *handlers.MovieHandler, uploadDir string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handlers.Health)

	absUpload, err := filepath.Abs(uploadDir)
	if err != nil {
		absUpload = uploadDir
	}
	fileServer := http.FileServer(http.Dir(absUpload))

	r.Route("/api/v1", func(r chi.Router) {
		r.Handle("/files/covers/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.StripPrefix("/api/v1/files/covers", fileServer).ServeHTTP(w, req)
		}))

		r.Route("/movies", func(r chi.Router) {
			r.Get("/", mh.List)
			r.Post("/", mh.Create)
			r.Get("/{id}", mh.Get)
			r.Patch("/{id}", mh.Update)
			r.Delete("/{id}", mh.Delete)
			r.Post("/{id}/cover", mh.UploadCover)
		})
	})

	return r
}
