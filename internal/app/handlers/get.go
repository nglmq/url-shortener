package handlers

import (
	"net/http"
	"strings"
)

func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", 400)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", 400)
		return
	}
	//slog.Info("%s", id)

	if originalURL, ok := us.Urls[id]; ok {
		w.Header().Set("Location", "https://"+originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Short URL not found", 400)
	}
}
