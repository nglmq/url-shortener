package handlers

import (
	"log/slog"
	"net/http"
	"strings"
)

func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", 400)
		slog.Info("Ошибка в статусе Only GET")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", 400)
		slog.Info("Ошибка в статусе при ID")
		return
	}
	slog.Info("%s", id)

	if originalURL, ok := us.Urls[id]; ok {
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		slog.Info("Ошибка в статусе TemporaryRedirect")
		slog.Info("%s", originalURL)
	} else {
		http.Error(w, "Short URL not found", 400)
		slog.Info("Ошибка в статусе Short URL not found")
	}
}
