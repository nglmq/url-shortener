package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func ResponseLogger(h http.HandlerFunc) http.HandlerFunc {
	logger := zap.NewExample()

	return func(w http.ResponseWriter, r *http.Request) {
		sugar := logger.Sugar()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		w = &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h(w, r)

		sugar.Infow("Incoming request",
			zap.String("Status Code", strconv.Itoa(responseData.status)),
			zap.String("Content length", strconv.Itoa(responseData.size)),
		)
	}
}
