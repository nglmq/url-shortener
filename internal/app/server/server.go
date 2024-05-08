package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"github.com/nglmq/url-shortener/internal/app/middleware"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"log"
	"net/http"
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	log.Printf("start path: %s", config.FlagInMemoryStorage)

	store := storage.NewMemoryURLStore()
	shortener := &handlers.URLShortener{
		Store: store,
	}

	if config.FlagInMemoryStorage != "" {
		fileStore, err := storage.NewFileStorage(config.FlagInMemoryStorage)
		if err != nil {
			return nil, err
		}
		shortener.FileStorage = fileStore

		if err = fileStore.ReadURLsFromFile(store.URLs); err != nil {
			log.Printf("Error reading URLs from file: %v", err)
			return nil, err
		}
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
