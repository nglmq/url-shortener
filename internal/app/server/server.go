package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"github.com/nglmq/url-shortener/internal/app/middleware"
	"net/http"
)

func Start() (http.Handler, error) {
	shortener := &handlers.URLShortener{
		URLs: make(map[string]string),
	}

	config.ParseFlags()

	r := chi.NewRouter()

	r.Use(middleware.RequestLogger)
	r.Use(middleware.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", shortener.ShortURLHandler)
		r.Route("/api/shorten", func(r chi.Router) {
			r.Post("/", shortener.JSONHandler)
		})
		r.Get("/{id}", shortener.GetURLHandler)
	})

	return r, nil
}
