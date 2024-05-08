package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/random"
	"log/slog"
	"net/http"
	"strconv"
)

type JSONRequest struct {
	URL string `json:"url" validate:"required"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

func (us *URLShortener) JSONHandler(w http.ResponseWriter, r *http.Request) {
	var (
		requestJSON  JSONRequest
		responseJSON JSONResponse
	)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&requestJSON); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.New().Struct(&requestJSON); err != nil {
		validateErr := err.Error()
		slog.Info(validateErr)

		http.Error(w, "url tag is required", http.StatusBadRequest)
		return
	}

	alias := random.NewRandomURL()
	err := us.Store.Add(alias, requestJSON.URL)
	if err != nil {
		http.Error(w, "Error saving URL JSON ", http.StatusBadRequest)
		return
	}
	if us.FileStorage != nil {
		if err := us.FileStorage.WriteURLsToFile(alias, requestJSON.URL); err != nil {
			http.Error(w, "Error writing URL to file", http.StatusInternalServerError)
			return
		}
	}

	shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + alias)
	contentLength := len(shortenedURL)

	responseJSON = JSONResponse{
		Result: shortenedURL,
	}

	responseData, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.Write(responseData)
}
