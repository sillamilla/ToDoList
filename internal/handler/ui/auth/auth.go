package auth

import (
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/auth"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
)

type Handler struct {
	auth           auth.Service
	signUpTemplate *template.Template
	signInTemplate *template.Template
}

func New(service auth.Service) Handler {
	signUp, err := template.ParseFiles("./internal/templates/users/sign-up.html")
	if err != nil {
		panic(err)
	}
	signIn, err := template.ParseFiles("./internal/templates/users/sign-in.html")
	if err != nil {
		panic(err)
	}
	return Handler{
		auth:           service,
		signUpTemplate: signUp,
		signInTemplate: signIn,
	}
}

func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	validationErr := chi.URLParam(r, "status")
	err := h.signUpTemplate.Execute(w, validationErr)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) SignUpPost(w http.ResponseWriter, r *http.Request) {
	userInput := models.UserInput{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	if validatorErr := errs.Validate(userInput); validatorErr != "" {
		link := fmt.Sprintf("/sign-up?status=%s", validatorErr)
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	user, err := h.auth.SignUp(r.Context(), userInput)
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
	err := h.signInTemplate.Execute(w, ok)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) SignInPost(w http.ResponseWriter, r *http.Request) {
	userInput := models.UserInput{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if validatorErr := errs.Validate(userInput); validatorErr != "" {
		link := fmt.Sprintf("/sign-in?status=%s", validatorErr)
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	user, err := h.auth.SignIn(r.Context(), userInput)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Redirect(w, r, "/sign-in?status=wrong username or password", http.StatusSeeOther)
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
	//todo if dont have it user not logout from base

	err = h.auth.Logout(r.Context(), session.Value)
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
