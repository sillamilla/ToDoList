package api

import (
	"ToDoWithKolya/internal/handler/api/task"
	"ToDoWithKolya/internal/handler/api/users"
	"ToDoWithKolya/internal/service"
)

type Handler struct {
	Task task.Handler
	User users.Handler
}

func New(srv *service.Service) Handler {
	return Handler{
		Task: task.New(srv.Tasks),
		User: users.New(srv.Users),
	}
}
