package users

import (
	"ToDoWithKolya/internal/service/users"
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
	signUp, err := template.ParseFiles("../../../templates/users/sign-up.html")
	if err != nil {
		panic(err)
	}
	signIn, err := template.ParseFiles("../../../templates/users/sign-in.html")
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
		http.Error(w, "Internal Server Error", 500)
	}
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {

	err := h.signIn.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
