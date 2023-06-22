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
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("unmarshal and validate, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation, err: %s", validationErrs))
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("user from ctx, err %v", ok))
		return
	}
	newtask.UserID = user.ID

	err = h.srv.Create(newtask)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("create, err %w", err))
		return
	}

	//todo wrong session error

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("user from ctx, err %v", ok))
		return
	}

	var updatedTask models.Task
	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &updatedTask)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("unmarshal, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("validation, err: %s", validationErrs))
		return
	}

	taskID, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("value from url, err %w", err))
		return
	}
	updatedTask.ID = taskID

	err = h.srv.Edit(updatedTask, user.ID)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("edit, err %w", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("value from url, err %w", err))
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusUnauthorized, fmt.Errorf("context, err %v", ok))
		return
	}

	task, err := h.srv.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helper.SendError(w, http.StatusNotFound, fmt.Errorf("task doesnt exist, err: %w", err))
			return
		}
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("sql err: %w", err))
		return
	}

	if user.ID != task.UserID {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("not your task, err: %v", false))
		return
	}

	if err = helper.SendJson(w, task, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("send json, err: %w", err))
		return
	}
}

func (h Handler) DeleteByTaskID(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Errorf("value from url, err: %w", err))
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("user from ctx, err: %v", ok))
		return
	}

	err = h.srv.DeleteByTaskID(id, user.ID)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("delete task by id, err: %w", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) GetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	user, b := ctxpkg.UserFromContext(r.Context())
	if !b {
		helper.SendError(w, http.StatusForbidden, fmt.Errorf("send json, err: %v", b))
	}

	tasks, err := h.srv.GetTasksByUserID(user.ID)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("get task by id, err: %w", err))
		return
	}

	if err = helper.SendJson(w, tasks, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Errorf("send json, err: %w", err))
		return
	}
}
