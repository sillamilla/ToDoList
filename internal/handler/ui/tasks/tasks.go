package tasks

import (
	"ToDoWithKolya/internal/ctxpkg"
	"ToDoWithKolya/internal/handler/api/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/task"
	"ToDoWithKolya/internal/templates/errs"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	ers    errs.Errors
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
	err := h.create.Execute(w, nil)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		log.Println("user from ctx err")
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	newTask := models.Task{
		UserID:      user.ID,
		Title:       title,
		Description: description,
	}

	if ok = errs.Validate(w, newTask); !ok {
		errs.ErrorWrap(w, fmt.Errorf("validation"), http.StatusBadRequest)
		return
	}

	err := h.srv.Create(newTask)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		log.Println(err.Error())
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		log.Println("user from ctx err")
		return
	}

	err = h.srv.DeleteByTaskID(id, user.ID)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		log.Println(err.Error())
		return
	}
	myTask, err := h.srv.GetByID(id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = h.edit.Execute(w, myTask)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		log.Println(err.Error())
		return
	}

	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		log.Println("user from ctx err")
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	newTask := models.Task{
		ID:          id,
		Title:       title,
		Description: description,
	}

	if ok = errs.Validate(w, newTask); !ok {
		errs.ErrorWrap(w, fmt.Errorf("validation"), http.StatusBadRequest)
		return
	}

	err = h.srv.Edit(newTask, user.ID)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
