package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/nglmq/url-shortener/config"
	"github.com/nglmq/url-shortener/internal/app/auth"
	"github.com/nglmq/url-shortener/internal/app/random"
)

type JSONRequest struct {
	URL string `json:"url" validate:"required"`
}

type JSONBatchRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required"`
}

type JSONResponse struct {
	Result string `json:"result"`
}

type JSONBatchResponse struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	ShortURL      string `json:"short_url" validate:"required"`
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

		token = &http.Cookie{Value: userToken}
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

	if us.DBStorage != nil {
		userID := auth.GetUserID(token.Value)

		existAlias, err := us.DBStorage.SaveURL(r.Context(), userID, alias, requestJSON.URL)
		if err != nil {
			http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
			return
		}
		if existAlias != alias {
			shortenedURL := fmt.Sprintf(config.FlagBaseURL + "/" + existAlias)
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
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Length", strconv.Itoa(contentLength))
			w.Write(responseData)

			return
		}
	}
	err = us.Store.Add(alias, requestJSON.URL)
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
	//if us.DBStorage != nil {
	//	if _, err := us.DBStorage.SaveURL(alias, requestJSON.URL); err != nil {
	//		http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
	//		return
	//	}
	//}

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

func (us *URLShortener) JSONBatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
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

		token = &http.Cookie{Value: userToken}
	}

	var requestJSON []JSONBatchRequest

	if err := json.NewDecoder(r.Body).Decode(&requestJSON); err != nil {
		http.Error(w, "error decoding request", http.StatusBadRequest)
		return
	}

	if len(requestJSON) == 0 {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	responseJSON := make([]JSONBatchResponse, len(requestJSON))
	//URLs := make(map[string]string)

	for i, req := range requestJSON {
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.Error()
			slog.Info(validateErr)

			http.Error(w, "missed json tag", http.StatusBadRequest)
			return
		}

		alias := random.NewRandomURL()
		//URLs[alias] = req.OriginalURL

		err := us.Store.Add(alias, req.OriginalURL)
		if err != nil {
			http.Error(w, "Error saving URL JSON ", http.StatusBadRequest)
			return
		}
		if us.FileStorage != nil {
			if err := us.FileStorage.WriteURLsToFile(alias, req.OriginalURL); err != nil {
				http.Error(w, "Error writing URL to file", http.StatusInternalServerError)
				return
			}
		}
		//if us.DBStorage != nil {
		//	if err := us.DBStorage.SaveBatch(URLs); err != nil {
		//		http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
		//		return
		//	}
		//}
		if us.DBStorage != nil {
			userID := auth.GetUserID(token.Value)

			if _, err := us.DBStorage.SaveURL(r.Context(), userID, alias, req.OriginalURL); err != nil {
				http.Error(w, "Error saving URL to database", http.StatusInternalServerError)
				return
			}
		}

		responseJSON[i] = JSONBatchResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      fmt.Sprintf(config.FlagBaseURL + "/" + alias),
		}
	}

	responseData, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//w.Header().Set("Content-Length")
	w.Write(responseData)
}
