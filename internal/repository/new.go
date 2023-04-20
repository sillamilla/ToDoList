package repository

import (
	"ToDoWithKolya/internal/repository/task"
	"ToDoWithKolya/internal/repository/user"
	"database/sql"
)

type Repo struct {
	TaskRepo task.TaskRepo
	UserRepo user.UserRepo
}

func New(db *sql.DB) *Repo {
	return &Repo{
		TaskRepo: task.Repo(db),
		UserRepo: user.Repo(db),
	}
}
