package main

import (
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"net/http"
)

func main() {
	shortener := &handlers.URLShortener{
		Urls: make(map[string]string),
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", shortener.ShortURLHandler)
	mux.HandleFunc("/{id}", shortener.GetURLHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
