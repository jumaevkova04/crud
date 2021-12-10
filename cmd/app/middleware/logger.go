package middleware

import (
	"log"
	"net/http"
)

// Logger ...
func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("START: %s %s", r.Method, r.URL.Path)

		handler.ServeHTTP(w, r)

		log.Printf("FINISH: %s %s", r.Method, r.URL.Path)
	})
}
