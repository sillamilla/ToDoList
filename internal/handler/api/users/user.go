package users

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/users"
	"fmt"
	"net/http"
)

type Handler struct {
	srv users.Service
}

func New(service users.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("unmarshal and validate, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation, err: %s", validationErrs))
		return
	}

	_, err = h.srv.SignUp(r.Context(), newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("register, err: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var newUser models.Input
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("need to be registrated, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation, err: %s", validationErrs))
		return
	}
	session, err := h.srv.SignIn(r.Context(), newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("login, err: %w", err))
		return
	}

	err = helper.SendJson(w, models.LoginResponse{Session: session}, http.StatusOK)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("send json, err: %w", err))
		return
	}

	//todo check return res( or writeHead(res))
}
