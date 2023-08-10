package users

import (
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/users"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
)

type Handler struct {
	srv    users.Service
	signUp *template.Template
	signIn *template.Template
}

func New(service users.Service) Handler {
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
	validationErr := chi.URLParam(r, "status")
	err := h.signUp.Execute(w, validationErr)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) SignUpPost(w http.ResponseWriter, r *http.Request) {
	userInput := models.Input{
		Login:    r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	if validatorErr := errs.Validate(userInput); validatorErr != "" {
		link := fmt.Sprintf("/sign-up?status=%s", validatorErr)
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	ok := h.srv.UserCheckExist(r.Context(), userInput.Login)
	if ok {
		http.Redirect(w, r, "/sign-up?status=this user already exist", http.StatusSeeOther)
		return
	}

	user, err := h.srv.SignUp(r.Context(), userInput)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: user.Session,
	}
	http.SetCookie(w, cookie)
	context.WithValue(r.Context(), "user", user)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ok := chi.URLParam(r, "status")
	err := h.signIn.Execute(w, ok)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) SignInPost(w http.ResponseWriter, r *http.Request) {
	userInput := models.Input{
		Login:    r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if validatorErr := errs.Validate(userInput); validatorErr != "" {
		link := fmt.Sprintf("/sign-in?status=%s", validatorErr)
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	user, err := h.srv.SignIn(r.Context(), userInput)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Redirect(w, r, "/sign-in?status=wrong login or password", http.StatusSeeOther)
			return
		}
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: user.Session,
	}
	http.SetCookie(w, cookie)
	context.WithValue(r.Context(), "user", user)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.srv.Logout(r.Context(), session.Value)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	context.WithValue(r.Context(), "user", nil)

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}
