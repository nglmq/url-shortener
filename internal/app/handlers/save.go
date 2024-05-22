package handlers

import (
	"context"
	"fmt"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/auth"
	"github.com/nglmq/url-shortener/internal/app/random"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"github.com/nglmq/url-shortener/internal/app/storage/db"
	"io"
	"net/http"
	"strconv"
)

type URLShortener struct {
	Store       storage.URLStore
	FileStorage *storage.FileStorage
	DBStorage   *db.PostgresStorage
}

func (us *URLShortener) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("userId")
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

	if us.DBStorage != nil {
		userID := auth.GetUserID(token.Value)

		existAlias, err := us.DBStorage.SaveURL(context.Background(), userID, alias, originalURL)
		if err != nil {
			http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
			return
		}
		if existAlias != alias {
			shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + existAlias)
			fmt.Println(shortenedURL)
			contentLength := len(shortenedURL)

			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", strconv.Itoa(contentLength))
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortenedURL))

			return
		}
	}

	err = us.Store.Add(alias, originalURL)
	if err != nil {
		http.Error(w, "Error saving URL", http.StatusBadRequest)
		return
	}
	if us.FileStorage != nil {
		if err := us.FileStorage.WriteURLsToFile(alias, originalURL); err != nil {
			http.Error(w, "Error writing URL to file", http.StatusInternalServerError)
			return
		}
	}

	shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + alias)
	contentLength := len(shortenedURL)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL))
}
