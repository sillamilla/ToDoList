package repository

import (
	"ToDoWithKolya/internal/repository/task"
	"ToDoWithKolya/internal/repository/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	TaskRepo task.TaskRepo
	UserRepo users.UserRepo
}

func New(db mongo.Database) *Repo {
	return &Repo{
		TaskRepo: task.NewTaskRepo(db),
		UserRepo: users.NewUserRepo(db),
	}
}
