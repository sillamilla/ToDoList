package tasks

import (
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/tasks"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	srv    tasks.Service
	create *template.Template
	edit   *template.Template
	home   *template.Template
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
	validationErr := chi.URLParam(r, "status")
	err := h.create.Execute(w, validationErr)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
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

	err := h.srv.NewTask(r.Context(), newTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	params := chi.URLParam(r, "search")
	tasks, err := h.srv.SearchTasks(r.Context(), params, user.ID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	//todo wtf
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
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errs.HandleError(w, fmt.Errorf("user from context"), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	err := h.srv.Delete(r.Context(), id, user.ID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errs.HandleError(w, fmt.Errorf("user from context"), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	task, err := h.srv.GetByID(r.Context(), user.ID, id)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	err = h.edit.Execute(w, task)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	updatedTask := models.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	//todo check
	//if validatorErr := errs.Validate(updatedTask); validatorErr != "" {
	//	link := fmt.Sprintf("/edit/{id}?status=%s", validatorErr)
	//	http.Redirect(w, r, link, http.StatusSeeOther)
	//	return
	//}

	err := h.srv.Edit(r.Context(), updatedTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
