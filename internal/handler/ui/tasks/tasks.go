package tasks

import (
	"ToDoWithKolya/internal/service/task"
	"net/http"
)

type Handler struct {
	srv task.Service
}

func NewHandler(service task.Service) Handler {
	return Handler{srv: service}
}

func (h Handler) Home(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) Edit(w http.ResponseWriter, r *http.Request) {

}
