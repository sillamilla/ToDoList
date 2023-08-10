package users

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/users"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Handler struct {
	srv users.Service
}

func New(service users.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userInput models.Input
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &userInput)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("unmarshal and validate, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("validation, err: %s", validationErrs))
		return
	}

	user, err := h.srv.SignUp(r.Context(), userInput)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("register, err: %w", err))
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: user.Session,
	}
	http.SetCookie(w, cookie)
	context.WithValue(r.Context(), "user", user)
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var newUser models.Input
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("need to be registrated, err: %w", err))
		return
	}
	if validationErrs != nil {
		errMsg := fmt.Sprintf("Validation errors: %s", strings.Join(validationErrs, ", "))
		helper.SendError(w, http.StatusBadRequest, errMsg)
		return
	}

	user, err := h.srv.SignIn(r.Context(), newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("login, err: %w", err))
		return
	}

	err = helper.SendJson(w, models.LoginResponse{Session: user.Session}, http.StatusOK)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("send json, err: %w", err))
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: user.Session,
	}
	http.SetCookie(w, cookie)
	context.WithValue(r.Context(), "user", user)
	//todo check context user
	//todo check return res( or writeHead(res))
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.srv.Logout(r.Context(), session.Value)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("logout, err: %w", err))
		return
	}

	context.WithValue(r.Context(), "user", nil)

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
