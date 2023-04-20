package service

import (
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service/task"
	"ToDoWithKolya/internal/service/user"
)

type Service struct {
	TaskSrv task.Service
	UserSrv user.Service
}

func New(rp *repository.Repo) *Service {
	return &Service{
		TaskSrv: task.NewTaskService(rp.TaskRepo),
		UserSrv: user.NewUserService(rp.UserRepo),
	}
}
