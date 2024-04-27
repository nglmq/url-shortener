package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
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

func RequestLogger(next http.Handler) http.Handler {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logFn := func(w http.ResponseWriter, r *http.Request) {

		sugar := *logger.Sugar()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		if r.Method == http.MethodPost {
			sugar.Infof("method: %s, url: %s, status: %d, duration: %s, size: %d",
				r.Method, r.RequestURI, responseData.status, duration, responseData.size)
		} else {
			sugar.Infof("method: %s, url: %s, status: %d, size: %d",
				r.Method, r.RequestURI, responseData.status, responseData.size)
		}
	}
	return http.HandlerFunc(logFn)
}
