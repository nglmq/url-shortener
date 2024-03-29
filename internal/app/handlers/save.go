package handlers

import (
	"fmt"
	"github.com/nglmq/url-shortener/internal/app/random"
	"io/ioutil"
	"net/http"
	"strconv"
)

type URLShortener struct {
	Urls map[string]string
}

func (us *URLShortener) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET requests are allowed!", 400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", 400)
		return
	}

	// Получение URL из тела запроса
	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "No URL provided", 400)
		return
	}

	alias := random.NewRandomUrl()
	us.Urls[alias] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", alias)
	contentLength := len(shortenedURL)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL))
}
