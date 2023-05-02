package users

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/users"
	"ToDoWithKolya/internal/templates/errs"
	"errors"
	"html/template"
	"net/http"
)

type Handler struct {
	ers    errs.Errors
	srv    users.Service
	signUp *template.Template
	signIn *template.Template
}

func NewHandler(service users.Service) Handler {
	signUp, err := template.ParseFiles("./internal/templates/users/sign-up.html")
	if err != nil {
		panic(err)
	}
	signIn, err := template.ParseFiles("./internal/templates/users/sign-in.html")
	if err != nil {
		panic(err)
	}
	return Handler{
		srv:    service,
		signUp: signUp,
		signIn: signIn,
	}
}

func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	//session, _ := r.Cookie("session")
	//if session != nil {
	//	errs.ErrorWrap(w, fmt.Errorf("u already logined"), http.StatusForbidden)
	//	return
	//}

	err := h.signUp.Execute(w, nil)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) SignUpPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	user := models.User{
		Login:    username,
		Password: password,
		Email:    email,
	}

	if ok := errs.Validate(w, user); !ok {
		return
	}

	session, err := h.srv.Register(user)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: session,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	//session, _ := r.Cookie("session")
	//if session != nil {
	//	errs.ErrorWrap(w, fmt.Errorf("u already logined"), http.StatusForbidden)
	//	return
	//}

	ok := r.URL.Query().Get("status")
	err := h.signIn.Execute(w, ok)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) SignInPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	req := models.LoginRequest{
		Login:    username,
		Password: password,
	}

	if ok := errs.Validate(w, req); !ok {
		return
	}

	session, err := h.srv.Login(req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Redirect(w, r, "/sign-in?status=false", http.StatusSeeOther)
			return
		}
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: session,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	h.srv.Logout(session.Value)
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}
