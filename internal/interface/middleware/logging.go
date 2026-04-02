package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware logs incoming requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("--> %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("<-- %s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// Recovery middleware recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
