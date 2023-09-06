package auth

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

func (h Handler) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if session == nil {
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		ses, err := h.sessionSrv.GetSession(session.Value)
		if err != nil {
			if errors.Is(err, models.ErrExpired) {
				h.Logout(w, r)
				http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
				return
			}

		}

		//ok, err := h.auth.LastActiveExpired(r.Context(), session.Value)
		//if err != nil {
		//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		//	return
		//}
		//if ok {
		//	h.Logout(w, r)
		//	http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
		//	return
		//}

		context.WithValue(r.Context(), "id", id)

		next.ServeHTTP(w, r)
	})
}
