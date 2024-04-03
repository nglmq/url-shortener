package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	shortener := &handlers.URLShortener{
		URLs: make(map[string]string),
	}
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", shortener.ShortURLHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", shortener.GetURLHandler)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
