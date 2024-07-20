package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"github.com/nglmq/url-shortener/internal/app/middleware"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"github.com/nglmq/url-shortener/internal/app/storage/db"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	store := storage.NewMemoryURLStore()
	shortener := &handlers.URLShortener{
		Store: store,
	}

	if config.DBConnection != "" {
		dbStorage, err := db.InitDBConnection()
		if err != nil {
			return nil, err
		}
		shortener.DBStorage = dbStorage
	}

	if config.FlagInMemoryStorage != "" && config.DBConnection == "" {
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
		r.Route("/api/shorten/batch", func(r chi.Router) {
			r.Post("/", shortener.JSONBatchHandler)
		})
		r.Get("/{id}", shortener.GetURLHandler)
		r.Get("/ping", shortener.PingDB)
		r.Get("/api/user/urls", shortener.GetAllURLsHandler)
		r.Delete("/api/user/urls", shortener.DeleteHandler)
	})

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))

	return r, nil
}
