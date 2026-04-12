package middleware

import (
	"net/http"
	"time"

	"github.com/dimasbaguspm/penster/pkg/observability"
)

// Logging middleware logs incoming requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := observability.NewLogger(r.Context(), "http", "middleware")
		log.Info("request started", "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(w, r)

		log.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}

// Recovery middleware recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log := observability.NewLogger(r.Context(), "http", "middleware")
				log.Error("panic recovered", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
