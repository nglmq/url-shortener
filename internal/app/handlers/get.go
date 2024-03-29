package handlers

import (
	"net/http"
	"strings"
)

func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if originalURL, ok := us.Urls[id]; ok {
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Short URL not found", http.StatusBadRequest)
	}
}
