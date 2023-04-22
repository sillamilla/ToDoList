package users

import (
	"context"
	"log"
	"net/http"
	"time"
)

func (h Handler) Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if len(key) != 44 {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("Authorization: Key error, key=%s", key)
			return
		}
		//todo провірити " " провірити на лен key
		user, err := h.srv.GetUserBySession(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("Authorization: GetUserBySession error, key=%s, err=%v", key, err)
			return
		}

		ctx := r.Context()
		ctxWithUser := context.WithValue(ctx, "users", user)

		lastActive, err := h.srv.GetSessionLastActive(key)
		if err != nil {
			log.Println("Не удалось получить время последней активности сессии:", err)
			return
		}

		sessionExpireTime := lastActive.Add(60 * time.Minute)
		if sessionExpireTime.Before(time.Now()) {
			err = h.srv.Logout(key)
			if err != nil {
				log.Println("Не удалось выполнить выход из сессии:", err)
			}
			return
		}

		r = r.WithContext(ctxWithUser)
		next(w, r)
	}
}
