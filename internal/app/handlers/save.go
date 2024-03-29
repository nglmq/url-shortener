package handlers

import (
	"fmt"
	"github.com/nglmq/url-shortener/internal/app/random"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
)

type URLShortener struct {
	Urls map[string]string
}

func (us *URLShortener) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", 400)
		slog.Info("статус Only POST")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", 400)
		slog.Info("статус Error reading request body")
		return
	}
	slog.Info("%s", string(body))

	// Получение URL из тела запроса
	originalURL := string(body)
	slog.Info("%s", originalURL)
	if originalURL == "" {
		http.Error(w, "No URL provided", 400)
		slog.Info("статус URL provided")
		return
	}

	alias := random.NewRandomUrl()
	slog.Info("%s", alias)
	us.Urls[alias] = originalURL
	slog.Info("%s", us.Urls)

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", alias)
	slog.Info("%s", shortenedURL)
	contentLength := len(shortenedURL)
	slog.Info("%s", contentLength)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL))
}
