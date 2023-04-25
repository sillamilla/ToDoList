package ui

import (
	"ToDoWithKolya/internal/handler/ui/tasks"
	"ToDoWithKolya/internal/handler/ui/users"
	"ToDoWithKolya/internal/service"
)

type Handler struct {
	TaskHandler tasks.Handler
	UserHandler users.Handler
}

func New(srv *service.Service) Handler {
	return Handler{
		TaskHandler: tasks.NewHandler(srv.TaskSrv),
		UserHandler: users.NewHandler(srv.UserSrv),
	}
}
