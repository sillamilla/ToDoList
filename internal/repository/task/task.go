package task

import (
	"ToDoWithKolya/internal/models"
	"database/sql"
)

type TaskRepo interface {
	Create(task models.Task) error
	Update(task models.Task) error

	GetTasksByUserID(userID int) ([]models.Task, error)
	GetByUserID(id int) (models.Task, error)
	GetByID(id int) (models.Task, error)

	DeleteByTaskID(id int, userID int) error
}

type taskRepo struct {
	db *sql.DB
}

func Repo(db *sql.DB) TaskRepo {
	return taskRepo{db: db}
}

func (r taskRepo) Create(task models.Task) error {
	_, err := r.db.Exec(
		"insert into tasks(user_id, title, description) values (?, ?, ?)",
		task.UserID,
		task.Title,
		task.Description,
	)
	return err
}

func (r taskRepo) GetByUserID(id int) (models.Task, error) {
	row := r.db.QueryRow("select * from tasks where user_id = ?", id)
	if row.Err() != nil {
		return models.Task{}, row.Err()
	}

	var task models.Task

	err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.CreatedDate)
	return task, err
}

func (r taskRepo) DeleteByTaskID(id int, userID int) error {
	_, err := r.db.Exec("delete from tasks where id = ? and user_id = ?", id, userID)
	return err
}

func (r taskRepo) Update(task models.Task) error {
	_, err := r.db.Exec(
		"update tasks set title = ?, description = ? where id = ?",
		task.Title,
		task.Description,
		task.ID,
	)
	return err
}

func (r taskRepo) GetTasksByUserID(userID int) ([]models.Task, error) {
	rows, err := r.db.Query("select * from tasks t where t.user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.CreatedDate)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, err
}

func (r taskRepo) GetByID(id int) (models.Task, error) {
	row := r.db.QueryRow("select * from tasks where id = ?", id)
	if row.Err() != nil {
		return models.Task{}, row.Err()
	}

	var task models.Task

	err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.CreatedDate)
	return task, err
}
