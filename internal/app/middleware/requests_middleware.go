package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RequestLogger(h http.HandlerFunc) http.HandlerFunc {
	logger := zap.NewExample()

	return func(w http.ResponseWriter, r *http.Request) {
		sugar := logger.Sugar()

		start := time.Now()

		h(w, r)

		duration := time.Since(start)

		sugar.Infow("Incoming request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("duration", duration.String()),
		)
	}
}
