package users

import (
	"ToDoWithKolya/internal/handler/api/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/users"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	srv users.Service
}

func NewHandler(service users.Service) Handler {
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

	_, err = h.srv.Register(newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("register, err: %w", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	readAll, err := io.ReadAll(r.Body)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("need to be registrated, err: %w", err))
		return
	}
	defer r.Body.Close()

	//todo setting models.Unmarshal and validate under todo

	var newUser models.LoginRequest
	if err = json.Unmarshal(readAll, &newUser); err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	session, err := h.srv.Login(newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("login, err: %w", err))
		return
	}

	helper.SendJson(w, models.LoginResponse{Session: session}, http.StatusOK)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("send json, err: %w", err))
		return
	}

	//todo check return res( or writeHead(res))
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Authorization")
	if len(key) != 44 {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("key, err: %s", key))
		return
	}

	err := h.srv.Logout(key)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("login, err: %w", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
