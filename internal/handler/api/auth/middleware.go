package auth

import (
	"net/http"
)

func (h Handler) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Header.Get("session")
		if session == "" {
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
