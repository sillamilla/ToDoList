package users

import (
	"ToDoWithKolya/internal/templates/errs"
	"context"
	"fmt"
	"net/http"
	"time"
)

func (h Handler) Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		if session == nil || len(session.Value) != 44 {
			http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
			return
		}

		user, err := h.srv.GetUserBySession(session.Value)
		if err != nil {
			errs.ErrorWrap(w, fmt.Errorf("user by session, err: %w", err), http.StatusInternalServerError)
			return
		}

		lastActive, err := h.srv.GetSessionLastActive(session.Value)
		if err != nil {
			errs.ErrorWrap(w, fmt.Errorf("last activie, err: %w", err), http.StatusInternalServerError)
			return
		}

		sessionExpireTime := lastActive.Add(30 * time.Minute)
		if sessionExpireTime.Before(time.Now()) {
			err = h.srv.Logout(session.Value)
			if err != nil {
				errs.ErrorWrap(w, fmt.Errorf("logout, err: %w", err), http.StatusInternalServerError)
			}
			return
		}

		ctx := r.Context()
		ctxWithUser := context.WithValue(ctx, "users", user)
		r = r.WithContext(ctxWithUser)
		next(w, r)
	}
}
