package users

import (
	"ToDoWithKolya/internal/service/users"
)

type Handler struct {
	srv users.Service
}

func New(service users.Service) Handler {
	return Handler{srv: service}
}
