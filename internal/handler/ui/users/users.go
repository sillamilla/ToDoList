package users

import (
	"ToDoWithKolya/internal/handler/api/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/users"
	"errors"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
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
	err := h.signUp.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	//todo validate
	session, err := h.srv.Register(user)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: session,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ok := r.URL.Query().Get("status")
	err := h.signIn.Execute(w, ok)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h Handler) SignInPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	req := models.LoginRequest{
		Login:    username,
		Password: password,
	}

	session, err := h.srv.Login(req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Redirect(w, r, "/sign-in?status=false", http.StatusSeeOther)
			return
		}
		log.Println(err)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: session,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
