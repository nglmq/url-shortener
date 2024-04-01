package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLShortener_ShortURLHandler(t *testing.T) {
	type want struct {
		code          int
		contentType   string
		contentLength string
	}
	tests := []struct {
		name        string
		requestBody string
		request     string
		want        want
	}{
		// TODO: Add test cases.
		{
			name: "simple test 1",
			want: want{
				code:          http.StatusCreated,
				contentType:   "text/plain",
				contentLength: "30",
			},
			requestBody: "https://practicum.yandex.ru/",
			request:     "/",
		},
		{
			name: "No URL provided test 1",
			want: want{
				code:          http.StatusBadRequest,
				contentType:   "text/plain; charset=utf-8",
				contentLength: "",
			},
			requestBody: "",
			request:     "/",
		},
		{
			name: "simple test 2",
			want: want{
				code:          http.StatusCreated,
				contentType:   "text/plain",
				contentLength: "30",
			},
			requestBody: "practicum.yandex.ru",
			request:     "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLShortenerTest := URLShortener{URLs: make(map[string]string)}
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			URLShortenerTest.ShortURLHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.contentLength, result.Header.Get("Content-Length"))

			defer result.Body.Close()
			_, err := io.ReadAll(result.Body)

			require.NoError(t, err)
		})
	}
}
