package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func LogMiddlew(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logResWr := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(logResWr, r)

		stop := time.Since(start)
		logger.Info("Request URL", zap.String("URL", r.RequestURI))
		logger.Info("Request Method", zap.String("Method", r.Method))
		logger.Info("Request Duration", zap.Duration("Duration", stop))
		logger.Info("Response Status", zap.Int("Status", logResWr.statusCode))
		logger.Info("Response Size", zap.Int("Size", logResWr.size))

	})
}

func (logResWr *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := logResWr.ResponseWriter.Write(data)
	logResWr.size += size
	return size, err
}

func (logResWr *loggingResponseWriter) WriteHeader(code int) {
	logResWr.statusCode = code
	logResWr.ResponseWriter.WriteHeader(code)
}
