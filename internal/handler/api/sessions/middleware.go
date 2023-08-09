package sessions

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

		lastActive, err := h.srv.LastActive(r.Context(), session.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sessionExpireTime := lastActive.Add(170 * time.Hour)
		if sessionExpireTime.Before(time.Now()) {
			h.Logout(w, r)
		}

		user, err := h.srv.GetUserID(r.Context(), session.Value)
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
