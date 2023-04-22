package users

import (
	"ToDoWithKolya/internal/handler/helper"
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
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation err: %s", validationErrs))
		return
	}

	err = h.srv.Register(newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	readAll, err := io.ReadAll(r.Body)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("need to be registrated, err: %e", err))
		return
	}
	defer r.Body.Close()

	var newUser models.LoginRequest
	if err = json.Unmarshal(readAll, &newUser); err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	//todo validate, wrong email or password(sql: no rows in result set)

	session, err := h.srv.Login(newUser)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(models.LoginResponse{Session: session})
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Authorization")
	if key == "" {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("seesion doesnt fiend"))
		return
	}

	h.srv.Logout(key)
	w.WriteHeader(http.StatusNoContent)
}
