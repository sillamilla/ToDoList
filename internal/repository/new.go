package repository

import (
	"ToDoWithKolya/internal/repository/task"
	"ToDoWithKolya/internal/repository/users"
	"database/sql"
)

type Repo struct {
	TaskRepo task.TaskRepo
	UserRepo users.UserRepo
}

func New(db *sql.DB) *Repo {
	return &Repo{
		TaskRepo: task.Repo(db),
		UserRepo: users.Repo(db),
	}
}
