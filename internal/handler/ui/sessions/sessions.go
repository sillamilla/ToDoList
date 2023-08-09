package sessions

import (
	"ToDoWithKolya/internal/handler/ui/errs"
	sessions "ToDoWithKolya/internal/service/sessions"
	"net/http"
)

type Handler struct {
	srv sessions.Service
}

func New(service sessions.Service) Handler {
	return Handler{
		srv: service,
	}
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	h.srv.Logout(r.Context(), session.Value)
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}
