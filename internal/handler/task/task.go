package task

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/task"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Handler struct {
	srv task.Service
}

func NewHandler(service task.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := helper.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("nil user"))
		return
	}
	readAll, err := io.ReadAll(r.Body)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer r.Body.Close()
	//такой сесії немає
	var newtask models.Task
	newtask.UserID = user.ID

	err = json.Unmarshal(readAll, &newtask)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	if errs := newtask.Validate(); len(errs) > 0 {
		helper.SendError(w, http.StatusUnprocessableEntity, fmt.Errorf("validate error"))
		return
	}

	err = h.srv.Create(newtask)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	readAll, err := io.ReadAll(r.Body)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer r.Body.Close()

	var newtask models.Task
	err = json.Unmarshal(readAll, &newtask)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	if errs := newtask.Validate(); len(errs) > 0 {
		helper.SendError(w, http.StatusUnprocessableEntity, fmt.Errorf("validate error"))
		return
	}

	user := r.Context().Value("user").(models.User)

	taskID, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}
	newtask.ID = taskID

	h.srv.Edit(newtask, user.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	task, err := h.srv.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.SendError(w, http.StatusNotFound, fmt.Errorf("task doesnt exist, err: %s", err))
			return
		}
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("sql err: %s", err))
		return
	}

	user, ok := helper.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusUnauthorized, fmt.Errorf("nil user"))
		return
	}

	if user.ID != task.UserID {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("not your task"))
		return
	}

	if err = helper.SendJson(w, task, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
	//нема таски з таким id
	//нема такой сесії
}

func (h Handler) DeleteByTaskID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := helper.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.srv.DeleteByTaskID(id, user.ID)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) GetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	user, b := helper.UserFromContext(r.Context())
	if !b {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("nil user"))
	}

	tasks, err := h.srv.GetTasksByUserID(user.ID)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if err = helper.SendJson(w, tasks, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
}
