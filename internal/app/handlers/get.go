package handlers

import (
	"log/slog"
	"net/http"
	"strings"
)

func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	us.mx.RLock()
	defer us.mx.RUnlock()

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	slog.Info(r.URL.String())

	if originalURL, ok := us.URLs[id]; ok {
		slog.Info(originalURL)
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			// Если протокол отсутствует, добавляем http:// по умолчанию
			originalURL = "http://" + originalURL
		}
		slog.Info(originalURL)
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Short URL not found", http.StatusBadRequest)
	}
}
