package users

import (
	"ToDoWithKolya/internal/handler/ui/errs"
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

		lastActive, _ := h.srv.GetSessionLastActive(session.Value)

		sessionExpireTime := lastActive.Add(30 * time.Minute)
		if sessionExpireTime.Before(time.Now()) {
			h.Logout(w, r)
		}

		user, err := h.srv.GetUserBySession(session.Value)
		if err != nil {
			errs.HandleError(w, fmt.Errorf("user by session, err: %w", err), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctxWithUser := context.WithValue(ctx, "users", user)
		r = r.WithContext(ctxWithUser)
		next(w, r)
	}
}
