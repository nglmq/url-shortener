package handlers

import (
	"net/http"
	"strings"
)

func (us *URLShortener) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var originalURL string
	var err error

	if us.DBStorage != nil {
		url, err := us.DBStorage.GetURL(id)
		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
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

	//url, err := us.Store.Get(id)
	//if err != nil {
	//	http.Error(w, "URL not found", http.StatusBadRequest)
	//	return
	//}
	//originalURL = url

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
