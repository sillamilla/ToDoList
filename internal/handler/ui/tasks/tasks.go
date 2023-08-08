package tasks

import (
	"ToDoWithKolya/internal/ctxpkg"
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/task"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	srv    task.Service
	create *template.Template
	edit   *template.Template
	home   *template.Template
}

func NewHandler(service task.Service) Handler {
	create, err := template.ParseFiles("./internal/templates/tasks/create.html")
	if err != nil {
		panic(err)
	}
	edit, err := template.ParseFiles("./internal/templates/tasks/edit.html")
	if err != nil {
		panic(err)
	}
	home, err := template.ParseFiles("./internal/templates/home.html")
	if err != nil {
		panic(err)
	}

	return Handler{
		srv:    service,
		create: create,
		edit:   edit,
		home:   home,
	}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	validationErr := r.URL.Query().Get("status")
	err := h.create.Execute(w, validationErr)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		errs.HandleError(w, fmt.Errorf("user from ctx"), http.StatusInternalServerError)
		return
	}

	newTask := models.Task{
		ID:          user.ID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	if validatorErr := errs.Validate(newTask); validatorErr != "" {
		link := fmt.Sprintf("/create?status=%s", validatorErr)
		http.Redirect(w, r, link, http.StatusSeeOther)
		return
	}

	err := h.srv.Create(newTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchValue := vars["search"]

	tasks, err := h.srv.SearchTask(searchValue)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	userAndTask := models.UserAndTask{
		User:  user,
		Tasks: tasks,
	}

	err = h.home.Execute(w, userAndTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := helper.FromURL(r, "taskID")

	status, err := strconv.Atoi(vars["status"])
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.srv.MarkValueSet(taskID, status)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := helper.FromURL(r, "id")

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		errs.HandleError(w, fmt.Errorf("user from context"), http.StatusInternalServerError)
		return
	}

	err := h.srv.DeleteByTaskID(id, user.ID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id := helper.FromURL(r, "id")

	myTask, err := h.srv.GetByID(id)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.edit.Execute(w, myTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	id := helper.FromURL(r, "id")

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		errs.HandleError(w, fmt.Errorf("user from context"), http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	newTask := models.Task{
		ID:          id,
		Title:       title,
		Description: description,
	}

	//if validatorErr := errs.Validate(newTask); validatorErr != "" {
	//	link := fmt.Sprintf("/edit/{id}?status=%s", validatorErr)
	//	http.Redirect(w, r, link, http.StatusSeeOther)
	//	return
	//}

	err := h.srv.Edit(newTask, user.ID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
