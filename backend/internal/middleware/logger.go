package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s - Status: %d - Duration: %v",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	}
}
