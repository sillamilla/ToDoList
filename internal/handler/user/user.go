package user

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	srv user.Service
}

func NewHandler(service user.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	readAll, err := io.ReadAll(r.Body)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer r.Body.Close()

	var newuser models.User

	err = json.Unmarshal(readAll, &newuser)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	if errs := newuser.Validate(); len(errs) > 0 {
		helper.SendError(w, http.StatusUnprocessableEntity, nil)
		return
	}

	err = h.srv.Register(newuser)
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

	var req models.LoginRequest
	if err = json.Unmarshal(readAll, &req); err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	//todo validate, wrong email or password(sql: no rows in result set)

	session, err := h.srv.Login(req)
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
