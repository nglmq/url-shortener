package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/auth"
)

// JSONAllUserURLs is a JSON response for all user URLs
type JSONAllUserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// JSONStats is a JSON response for GetStats
type JSONStats struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// GetURLHandler is handler to get user URL.
func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
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

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if us.DBStorage != nil {
		url, deleted, err := us.DBStorage.GetURL(r.Context(), id)
		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
		if deleted {
			w.WriteHeader(http.StatusGone)
			return
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		url, err := us.Store.Get(id)
		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// GetAllURLsHandler is handler to get all user URLs. Works only for authenticated users.
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

	urls, err := us.DBStorage.GetAllUserURLs(r.Context(), auth.GetUserID(token.Value))
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

// GetStats is a handler to get stats about urls and users quantity
func (us *URLShortener) GetStats(w http.ResponseWriter, r *http.Request) {
	if config.TrustedSubnet == "" {
		http.Error(w, "Not allowed", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	clientIP := r.Header.Get("X-Real-IP")

	_, subnet, _ := net.ParseCIDR(config.TrustedSubnet)

	ip := net.ParseIP(clientIP)
	if !subnet.Contains(ip) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	uriQuantity, usersQuantity, err := us.DBStorage.GetStats(r.Context())
	if err != nil {
		http.Error(w, "Error getting URLs", http.StatusInternalServerError)
		return
	}

	jsonStats := JSONStats{
		URLs:  uriQuantity,
		Users: usersQuantity,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonStats); err != nil {
		http.Error(w, "Error encoding URLs", http.StatusInternalServerError)
		return
	}
}
