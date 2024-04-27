package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/random"
	"io"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Something went wrong",
			http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &requestJSON)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	alias := random.NewRandomURL()
	us.URLs[alias] = requestJSON.URL

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
