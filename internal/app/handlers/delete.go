package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nglmq/url-shortener/internal/app/auth"
	"io"
	"net/http"
)

func (us *URLShortener) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE requests are allowed!", http.StatusBadRequest)
		return
	}

	token, err := r.Cookie("userID")
	if err != nil || token == nil {
		userToken, err := auth.BuildJWTString()
		if err != nil {
			http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "userID",
			Value:    userToken,
			Path:     "/",
			HttpOnly: true,
		})

		token = &http.Cookie{Value: userToken}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var aliases []string
	if err := json.Unmarshal(body, &aliases); err != nil {
		http.Error(w, "Error unmarshalling body", http.StatusBadRequest)
		return
	}

	userID := auth.GetUserID(token.Value)

	for _, alias := range aliases {
		go func(alias, userID string) {
			err := us.DBStorage.DeleteURL(context.Background(), alias, userID)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(alias + " deleted")
		}(alias, userID)
	}

	w.WriteHeader(http.StatusAccepted)

	fmt.Println(aliases)
}