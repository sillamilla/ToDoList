package api

import (
	"ToDoWithKolya/internal/handler/api/auth"
	"ToDoWithKolya/internal/handler/api/task"
	"ToDoWithKolya/internal/handler/api/users"
	"ToDoWithKolya/internal/service"
)

type Handler struct {
	Task task.Handler
	User users.Handler
	Auth auth.Handler
}

func New(srv *service.Service) Handler {
	return Handler{
		Task: task.New(srv.Task),
		User: users.New(srv.User),
		Auth: auth.New(srv.Auth),
	}
}
