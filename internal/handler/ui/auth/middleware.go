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
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		ok, err := h.auth.LastActiveExpired(r.Context(), session.Value) //todo forever new, and not expired
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if ok {
			h.Logout(w, r)
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		sessionInfo, err := h.auth.GetSessionInfo(r.Context(), session.Value)
		if err != nil || sessionInfo.UserID == "" {
			h.Logout(w, r)
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "id", sessionInfo.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
