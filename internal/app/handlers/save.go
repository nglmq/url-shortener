package handlers

import (
	"fmt"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/random"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type URLShortener struct {
	URLs map[string]string
	mx   sync.RWMutex
}

func (us *URLShortener) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	us.mx.Lock()
	defer us.mx.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Получение URL из тела запроса
	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "No URL provided", http.StatusBadRequest)
		return
	}

	alias := random.NewRandomURL()
	us.URLs[alias] = originalURL

	shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + alias)
	contentLength := len(shortenedURL)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL))
}
