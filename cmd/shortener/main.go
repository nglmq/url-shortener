package main

import (
	"github.com/nglmq/url-shortener/internal/app/handlers"
	"net/http"
)

func main() {
	shortener := &handlers.URLShortener{
		URLs: make(map[string]string),
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			shortener.ShortURLHandler(w, r)
		} else if r.Method == http.MethodGet {
			shortener.GetURLHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
