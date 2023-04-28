package users

import (
	"ToDoWithKolya/internal/handler/api/helper"
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
			helper.SendError(w, http.StatusUnauthorized, fmt.Errorf("user by session, key: %s \n err: %v", session, err))
			return
		}

		lastActive, err := h.srv.GetSessionLastActive(session.Value)
		if err != nil {
			helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("last activie, key: %s \n err: %v", session, err))
			return
		}

		sessionExpireTime := lastActive.Add(180 * time.Minute)
		if sessionExpireTime.Before(time.Now()) {
			err = h.srv.Logout(session.Value)
			if err != nil {
				helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("logout, err: %w", err))
			}
			return
		}

		ctx := r.Context()
		ctxWithUser := context.WithValue(ctx, "users", user)
		r = r.WithContext(ctxWithUser)
		next(w, r)
	}
}
