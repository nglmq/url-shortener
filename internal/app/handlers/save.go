package handlers

import (
	"fmt"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/random"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"io"
	"net/http"
	"strconv"
)

type URLShortener struct {
	Store storage.URLStore
}

func (us *URLShortener) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "No URL provided", http.StatusBadRequest)
		return
	}

	alias := random.NewRandomURL()
	err = us.Store.Add(alias, originalURL)
	if err != nil {
		http.Error(w, "Error saving URL", http.StatusBadRequest)
		return
	}
	storage.WriteURLsToFile(config.FlagInMemoryStorage, alias, originalURL)

	shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + alias)
	contentLength := len(shortenedURL)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL))
}
