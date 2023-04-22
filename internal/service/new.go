package service

import (
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service/task"
	"ToDoWithKolya/internal/service/users"
)

type Service struct {
	TaskSrv task.Service
	UserSrv users.Service
}

func New(rp *repository.Repo) *Service {
	return &Service{
		TaskSrv: task.NewTaskService(rp.TaskRepo),
		UserSrv: users.NewUserService(rp.UserRepo),
	}
}
