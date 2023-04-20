package handler

import (
	"ToDoWithKolya/internal/handler/task"
	"ToDoWithKolya/internal/handler/user"
	"ToDoWithKolya/internal/service"
)

type Handler struct {
	TaskHandler task.Handler
	UserHandler user.Handler
}

func New(srv *service.Service) Handler {
	return Handler{
		TaskHandler: task.NewHandler(srv.TaskSrv),
		UserHandler: user.NewHandler(srv.UserSrv),
	}
}
