package handlers

import (
	"context"
	"encoding/json"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/auth"
	"net/http"
	"strings"
)

type JSONAllUserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	var (
		originalURL string
		token       *http.Cookie
		err         error
	)

	token, err = r.Cookie("userId")
	if err != nil || token == nil {
		userToken, err := auth.BuildJWTString()
		if err != nil {
			http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "userId",
			Value:    userToken,
			Path:     "/",
			HttpOnly: true,
		})

		token = &http.Cookie{Value: userToken}
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if us.DBStorage != nil {
		url, deleted, err := us.DBStorage.GetURL(context.Background(), id)
		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
		if deleted {
			w.WriteHeader(http.StatusGone)
			return
		}

		originalURL = url
	} else {
		originalURL, err = us.Store.Get(id)
		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (us *URLShortener) GetAllURLsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("userId")
	if err != nil || token == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := us.DBStorage.GetAllUserURLs(context.Background(), auth.GetUserID(token.Value))
	if err != nil {
		http.Error(w, "Error getting URLs", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		http.Error(w, "No URLs found", http.StatusNoContent)
		return
	}

	jsonURLs := make([]JSONAllUserURLs, 0, len(urls))
	for alias, originalURL := range urls {
		jsonURLs = append(jsonURLs, JSONAllUserURLs{
			ShortURL:    config.FlagBaseURL + "/" + alias,
			OriginalURL: originalURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonURLs); err != nil {
		http.Error(w, "Error encoding URLs", http.StatusInternalServerError)
		return
	}
}
