package task

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/tasks"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Handler struct {
	srv tasks.Service
}

func New(service tasks.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &task)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Sprintf("unmarshal and validate, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("validation, err: %s", validationErrs))
		return
	}

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("user from ctx, err %v", ok))
		return
	}

	task.UserID = user.ID
	err = h.srv.NewTask(r.Context(), task)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("create, err %w", err))
		return
	}

	//todo wrong session error
	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	var updatedTask models.Task

	validationErrs, err := helper.UnmarshalAndValidate(r.Body, &updatedTask)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("unmarshal, err: %w", err))
		return
	}
	if validationErrs != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("validation, err: %s", validationErrs))
		return
	}

	err = h.srv.Edit(r.Context(), updatedTask)
	if err != nil {
		helper.SendError(w, http.StatusBadRequest, fmt.Sprintf("edit, err %w", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		helper.SendError(w, http.StatusUnauthorized, fmt.Sprintf("user not found in ctx"))
		return
	}

	id := chi.URLParam(r, "id")
	task, err := h.srv.GetByID(r.Context(), user.ID, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			helper.SendError(w, http.StatusNotFound, fmt.Sprintf("tasks doesnt exist, err: %w", err))
			return
		}
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("get by id, err: %w", err))
		return
	}

	if err = helper.SendJson(w, task, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("send json, err: %w", err))
		return
	}
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		helper.SendError(w, http.StatusBadRequest, fmt.Sprintf("user not found in ctx"))
		return
	}

	id := chi.URLParam(r, "id")
	err := h.srv.Delete(r.Context(), user.ID, id)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("delete tasks by id, err: %", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		helper.SendError(w, http.StatusUnauthorized, fmt.Sprintf("user not found in ctx"))
		return
	}

	tasks, err := h.srv.GetTasks(r.Context(), user.ID)
	if err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("get tasks by id, err: %w", err))
		return
	}

	if err = helper.SendJson(w, tasks, http.StatusOK); err != nil {
		helper.SendError(w, http.StatusInternalServerError, fmt.Sprintf("send json, err: %w", err))
		return
	}
}
