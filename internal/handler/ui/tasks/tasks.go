package tasks

import (
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/tasks"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	srv            tasks.Service
	createTemplate *template.Template
	editTemplate   *template.Template
	searchTemplate *template.Template
	homeTemplate   *template.Template
}

func New(service tasks.Service) Handler {
	create, err := template.ParseFiles("./internal/templates/tasks/create.html")
	if err != nil {
		panic(err)
	}
	edit, err := template.ParseFiles("./internal/templates/tasks/edit.html")
	if err != nil {
		panic(err)
	}
	search, err := template.ParseFiles("./internal/templates/tasks/search.html")
	if err != nil {
		panic(err)
	}
	home, err := template.ParseFiles("./internal/templates/home.html")
	if err != nil {
		panic(err)
	}

	return Handler{
		srv:            service,
		createTemplate: create,
		editTemplate:   edit,
		searchTemplate: search,
		homeTemplate:   home,
	}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	err := h.createTemplate.Execute(w, nil)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		errs.HandleError(w, nil, http.StatusInternalServerError)
		return
	}

	newTask := models.Task{ //todo do in service
		ID:          uuid.NewString(),
		UserID:      id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		CreatedAt:   time.Now(),
	}

	if validatorErr := errs.Validate(newTask); validatorErr != "" {
		errs.HandleError(w, errors.New(validatorErr), http.StatusBadRequest)
		return
	}

	err := h.srv.NewTask(r.Context(), newTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		errs.HandleError(w, nil, http.StatusInternalServerError)
		return
	}

	taskName := chi.URLParam(r, "search")
	tasks, err := h.srv.SearchTasks(r.Context(), taskName, id)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.searchTemplate.Execute(w, tasks)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	status, err := strconv.Atoi(chi.URLParam(r, "status"))
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	err = h.srv.MarkValueSet(r.Context(), id, status)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	err := h.srv.Delete(r.Context(), taskID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, err := h.srv.GetByID(r.Context(), taskID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.editTemplate.Execute(w, task)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	updatedTask := models.Task{
		ID:          chi.URLParam(r, "id"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	if validatorErr := errs.Validate(updatedTask); validatorErr != "" {
		errs.HandleError(w, errors.New(validatorErr), http.StatusBadRequest)
		return
	}

	err := h.srv.Edit(r.Context(), updatedTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
