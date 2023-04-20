package user

import (
	"context"
	"log"
	"net/http"
	"unicode/utf8"
)

func (h Handler) Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if utf8.RuneCountInString(key) != 44 {
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
		ctxWithUser := context.WithValue(ctx, "user", user)

		r = r.WithContext(ctxWithUser)
		next(w, r)
	}
}
