package auth

import (
	"context"
	"net/http"
	"time"
)

func (h Handler) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		info, err := h.auth.GetSessionInfo(r.Context(), session.Value)
		if err != nil || info.UserID == "" {
			h.Logout(w, r)
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		sessionExpireTime := info.CreatedAt.Add(170 * time.Hour)
		if sessionExpireTime.Before(time.Now()) {
			h.Logout(w, r)
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "id", info.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
