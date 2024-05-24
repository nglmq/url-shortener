package handlers

import (
	"github.com/nglmq/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLShortener_GetURLHandler(t *testing.T) {
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name    string
		request string
		URLs    map[string]string
		want    want
	}{
		// TODO: Add test cases.
		{
			name: "simple test 1",
			URLs: map[string]string{
				"abcdefgh": "practicum.yandex.ru",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "http://practicum.yandex.ru",
			},
			request: "/abcdefgh",
		},
		{
			name: "Short URL not found",
			URLs: map[string]string{
				"abcdefgh": "practicum.yandex.ru",
			},
			want: want{
				location: "",
				code:     http.StatusBadRequest,
			},
			request: "/abcdefgh35",
		},
		{
			name: "Invalid request Id=0",
			URLs: map[string]string{
				"": "",
			},
			want: want{
				location: "",
				code:     http.StatusBadRequest,
			},
			request: "/",
		},
		{
			name: "simple test 4",
			URLs: map[string]string{
				"abcdefgh1": "https://github.com/stretchr/testify?tab=readme-ov-file",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://github.com/stretchr/testify?tab=readme-ov-file",
			},
			request: "/abcdefgh1",
		},
		{
			name: "simple test 5",
			URLs: map[string]string{
				"abcdefgh": "https://megamarket.ru/personal/order/view/",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://megamarket.ru/personal/order/view/",
			},
			request: "/abcdefgh",
		},
		{
			name: "simple test 6",
			URLs: map[string]string{
				"abcdefghdsghsdheh": "practicum.yandex.ru",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "http://practicum.yandex.ru",
			},
			request: "/abcdefghdsghsdheh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemoryURLStore()
			store.URLs = tt.URLs
			URLShortenerTest := URLShortener{Store: store}

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			URLShortenerTest.GetURLHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			defer result.Body.Close()
			_, err := io.ReadAll(result.Body)

			require.NoError(t, err)
		})
	}
}
