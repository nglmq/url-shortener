package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/random"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type JSONRequest struct {
	URL string `json:"url"`
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
		if err == io.EOF {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("{}"))
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validator.New().Struct(&requestJSON); err != nil {
		validateErr := err.Error()

		slog.Error(validateErr)
		http.Error(w, "url tag is required", http.StatusBadRequest)
		return
	}

	alias := random.NewRandomURL()
	us.URLs[alias] = requestJSON.URL
	log.Println(requestJSON.URL)

	shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + alias)
	log.Println(shortenedURL)
	contentLength := len(shortenedURL)

	responseJSON = JSONResponse{
		Result: shortenedURL,
	}
	log.Println(responseJSON.Result)

	responseData, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}
	log.Println(string(responseData))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.Write(responseData)
	log.Println(w.Header().Get("Content-Length"), w.Header().Get("Content-Length"))
}
