package handler

import (
	"ToDoWithKolya/internal/handler/task"
	"ToDoWithKolya/internal/handler/users"
	"ToDoWithKolya/internal/service"
)

type Handler struct {
	TaskHandler task.Handler
	UserHandler users.Handler
}

func New(srv *service.Service) Handler {
	return Handler{
		TaskHandler: task.NewHandler(srv.TaskSrv),
		UserHandler: users.NewHandler(srv.UserSrv),
	}
}
