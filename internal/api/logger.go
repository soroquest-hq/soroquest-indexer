package api

import (
	"log"
	"net/http"
	"time"
)

// responseWriterInterceptor captures the HTTP status code
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// RequestLogger logs incoming HTTP requests with method, path, status, and duration
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		interceptor := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(interceptor, r)

		duration := time.Since(start)
		log.Printf("[HTTP] %s %s %d %s", r.Method, r.URL.Path, interceptor.statusCode, duration)
	})
}
