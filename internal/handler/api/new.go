package api

import (
	"ToDoWithKolya/internal/handler/api/task"
	"ToDoWithKolya/internal/handler/api/users"
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
