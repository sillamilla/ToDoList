package tasks

import (
	"ToDoWithKolya/internal/ctxpkg"
	"ToDoWithKolya/internal/handler/api/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/task"
	"errors"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	srv  task.Service
	task *template.Template
}

func NewHandler(service task.Service) Handler {
	task, err := template.ParseFiles("./internal/templates/tasks/task.html")
	if err != nil {
		panic(err)
	}
	return Handler{
		srv:  service,
		task: task,
	}
}

func (h Handler) Task(w http.ResponseWriter, r *http.Request) {
	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		return
	}

	id, err := helper.GetIntFromURL(r, "id")
	if err != nil {
		log.Println(err.Error())
		return
	}

	myTask, err := h.srv.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Not Fount", http.StatusNotFound)
		}
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if user.ID != myTask.UserID {
		log.Println(err.Error())
		return
	}

	h.task.Execute(w, myTask)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {

}
