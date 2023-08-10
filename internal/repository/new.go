package repository

import (
	"ToDoWithKolya/internal/repository/tasks"
	"ToDoWithKolya/internal/repository/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	Task tasks.Tasks
	User users.Users
}

func New(db *mongo.Database) *Repo {
	return &Repo{
		Task: tasks.New(db),
		User: users.New(db),
	}
}
