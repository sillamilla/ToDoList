package repository

import (
	"ToDoWithKolya/internal/repository/sessions"
	"ToDoWithKolya/internal/repository/tasks"
	"ToDoWithKolya/internal/repository/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	Task    tasks.Repo
	User    users.Repo
	Session sessions.Repo
}

func New(db *mongo.Database) *Repo {
	return &Repo{
		Task:    tasks.New(db),
		User:    users.New(db),
		Session: sessions.New(db),
	}
}
