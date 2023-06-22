package users

import (
	"ToDoWithKolya/internal/handler/helper"
	"context"
	"fmt"
	"net/http"
	"time"
)

func (h Handler) Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if len(key) != 44 {
			helper.SendError(w, http.StatusUnauthorized, fmt.Errorf("key, err: %s", key))
			return
		}

		user, err := h.srv.GetUserBySession(key)
		if err != nil {
			helper.SendError(w, http.StatusUnauthorized, fmt.Errorf("user by session, key: %s \n err: %v", key, err))
			return
		}

		ctx := r.Context()
		ctxWithUser := context.WithValue(ctx, "users", user)

		lastActive, err := h.srv.GetSessionLastActive(key)
		if err != nil {
			helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("last activie, key: %s \n err: %v", key, err))
			return
		}

		sessionExpireTime := lastActive.Add(180 * time.Minute)
		if sessionExpireTime.Before(time.Now()) {
			err = h.srv.Logout(key)
			if err != nil {
				helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("logout, err: %w", err))
			}
			return
		}

		r = r.WithContext(ctxWithUser)
		next(w, r)
	}
}
