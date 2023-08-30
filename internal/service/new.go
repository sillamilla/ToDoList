package service

import (
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service/auth"
	"ToDoWithKolya/internal/service/sessions"
	"ToDoWithKolya/internal/service/tasks"
	"ToDoWithKolya/internal/service/users"
)

type Service struct {
	Task    tasks.Service
	User    users.Service
	Session sessions.Service
	Auth    auth.Service
}

func New(rp *repository.Repo) *Service {
	userSrv := users.New(rp.User)
	sessionSrv := sessions.New(rp.Session)

	return &Service{
		Task:    tasks.New(rp.Task),
		User:    userSrv,
		Session: sessionSrv,
		Auth:    auth.New(sessionSrv, userSrv),
	}
}
