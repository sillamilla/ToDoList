package users

import (
	"net/http"
)

func (h Handler) Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		ok, err := h.srv.LastActiveExpired(r.Context(), session.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if ok {
			h.Logout(w, r)
		}

		next(w, r)
	}
}
