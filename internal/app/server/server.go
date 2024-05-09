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
)

func Start() (http.Handler, error) {
	config.ParseFlags()

	store := storage.NewMemoryURLStore()
	//dbStorage, err := db.InitDBConnection()
	//if err != nil {
	//	return nil, err
	//}

	shortener := &handlers.URLShortener{
		Store: store,
		//DBStorage: dbStorage,
	}

	switch config.DBConnection != "" {
	case config.DBConnection != "":
		dbStorage, err := db.InitDBConnection()
		if err != nil {
			return nil, err
		}

		shortener.DBStorage = dbStorage
	//
	//case config.DBConnection != "" && config.FlagInMemoryStorage != "":
	//	dbStorage, err := db.InitDBConnection()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	shortener.DBStorage = dbStorage

	case config.DBConnection == "" && config.FlagInMemoryStorage != "":
		fileStore, err := storage.NewFileStorage(config.FlagInMemoryStorage)
		if err != nil {
			return nil, err
		}
		shortener.FileStorage = fileStore

		if err = fileStore.ReadURLsFromFile(store.URLs); err != nil {
			log.Printf("Error reading URLs from file: %v", err)
			return nil, err
		}

	case config.DBConnection == "" && config.FlagInMemoryStorage == "":
		return nil, nil
	}

	//if config.FlagInMemoryStorage != "" {
	//	fileStore, err := storage.NewFileStorage(config.FlagInMemoryStorage)
	//	if err != nil {
	//		return nil, err
	//	}
	//	shortener.FileStorage = fileStore
	//
	//	if err = fileStore.ReadURLsFromFile(store.URLs); err != nil {
	//		log.Printf("Error reading URLs from file: %v", err)
	//		return nil, err
	//	}
	//}

	r := chi.NewRouter()

	r.Use(middleware.RequestLogger)
	r.Use(middleware.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", shortener.ShortURLHandler)
		r.Route("/api/shorten", func(r chi.Router) {
			r.Post("/", shortener.JSONHandler)
		})
		r.Get("/{id}", shortener.GetURLHandler)
		r.Get("/ping", shortener.PingDB)
	})

	return r, nil
}
