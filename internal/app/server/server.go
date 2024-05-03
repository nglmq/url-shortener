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

	// Initialize URL storage
	store := storage.NewMemoryURLStore()

	// Optional: Load URLs from a file if a path is provided
	if config.FlagInMemoryStorage != "" {
		err := storage.CreateFile(config.FlagInMemoryStorage)
		if err != nil {
			return nil, err
		}

		if err := storage.ReadURLsFromFile(config.FlagInMemoryStorage, store.URLs); err != nil {
			log.Printf("Error reading URLs from file: %v", err)
			return nil, err
		}
	}

	// Create URLShortener handler with the initialized store
	shortener := &handlers.URLShortener{
		Store: store,
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
