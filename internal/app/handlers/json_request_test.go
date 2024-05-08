package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/nglmq/url-shortener/internal/app/storage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLShortener_JSONHandler(t *testing.T) {
	us := URLShortener{
		Store: storage.NewMemoryURLStore(),
	}

	ts := httptest.NewServer(http.HandlerFunc(us.JSONHandler))
	defer ts.Close()

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid POST Request",
			method:         http.MethodPost,
			body:           `{"url": "https://google.com"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Method",
			method:         http.MethodGet,
			body:           `{"url": "https://google.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing URL in Request Body",
			method:         http.MethodPost,
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid request body",
			method:         http.MethodPost,
			body:           `{"ul": "https://google.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, ts.URL, bytes.NewBufferString(tc.body))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == http.StatusCreated {
				var jsonResponse JSONResponse
				err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
				require.NoError(t, err)
				assert.Equal(t, 9, len(jsonResponse.Result))
			}
		})
	}
}
