package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"github.com/nglmq/url-shortener/internal/app/middleware"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"log"
	"net/http"
	"os"
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	shortener := &handlers.URLShortener{
		URLs: make(map[string]string),
	}

	if config.FlagInMemoryStorage != "" {
		storage.CreateFile(config.FlagInMemoryStorage)
		err := storage.ReadURLsFromFile(config.FlagInMemoryStorage, shortener.URLs)
		if err != nil {
			log.Printf("Error reading URLs from file: %v", err)
		}
		b, err := os.ReadFile(config.FlagInMemoryStorage)
		log.Println("Storage path: ", config.FlagInMemoryStorage, b)
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
