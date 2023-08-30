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

		ok, err := h.auth.LastActiveExpired(r.Context(), session)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if ok {
			h.Logout(w, r)
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
