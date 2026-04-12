package middleware

import (
	"net/http"
	"time"

	"github.com/dimasbaguspm/penster/pkg/observability"
)

// Logging middleware logs incoming requests and records HTTP metrics
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := observability.NewLogger(r.Context(), "http", "middleware")
		log.Info("request started", "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", duration.String())

		observability.HTTPRequestsTotal.Add(r.Context(), 1)
		observability.HTTPRequestDuration.Record(r.Context(), float64(duration.Milliseconds()))
	})
}

// Recovery middleware recovers from panics and records panic metrics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log := observability.NewLogger(r.Context(), "http", "middleware")
				log.Error("panic recovered", "error", err)
				observability.HTTPPanicCount.Add(r.Context(), 1)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
