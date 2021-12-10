package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Basic ...
func Basic(checkAuth func(string, string) bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			auth := strings.SplitN(r.Header.Get("Authorization"), "", 2)

			if len(auth) != 2 || auth[0] != "Basic" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)

			if len(pair) != 2 {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			login := pair[0]
			password := pair[1]
			ok := checkAuth(login, password)

			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
