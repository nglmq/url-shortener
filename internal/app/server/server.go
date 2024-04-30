package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"github.com/nglmq/url-shortener/internal/app/middleware"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"net/http"
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	shortener := &handlers.URLShortener{
		URLs: make(map[string]string),
	}
	err := storage.ReadURLsFromFile(config.FlagInMemoryStorage, shortener.URLs)
	if err != nil {
		storage.CreateFile(config.FlagInMemoryStorage)
	}

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
