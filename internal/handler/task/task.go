package task

import (
	"ToDoWithKolya/internal/ctxpkg"
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/task"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type Handler struct {
	srv task.Service
}

func NewHandler(service task.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var newtask models.Task
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &newtask)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("unmarshalAndValidate err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation err: %s", validationErrs))
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
	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("userFromContext err"))
		return
	}

	var updatedTask models.Task
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &updatedTask)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation err: %s", validationErrs))
		return
	}

	taskID, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}
	updatedTask.ID = taskID

	h.srv.Edit(updatedTask, user.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusUnauthorized, fmt.Errorf("nil users"))
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

	if user.ID != task.UserID {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("not your task"))
		return
	}

	if err = helper.SendJson(w, task, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h Handler) DeleteByTaskID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
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
	user, b := ctxpkg.UserFromContext(r.Context())
	if !b {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("nil users"))
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
