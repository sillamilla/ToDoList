package service

import (
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service/sessions"
	"ToDoWithKolya/internal/service/tasks"
	"ToDoWithKolya/internal/service/users"
)

type Service struct {
	Tasks    tasks.Service
	Users    users.Service
	Sessions sessions.Service
}

func New(rp *repository.Repo) *Service {
	return &Service{
		Tasks:    tasks.New(rp.Task),
		Users:    users.New(rp.User),
		Sessions: sessions.New(rp.Session),
	}
}
