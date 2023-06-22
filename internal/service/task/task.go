package task

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/task"
	"fmt"
)

type Service interface {
	Create(task models.Task) error

	GetByUserID(id int) (models.Task, error)
	GetTasksByUserID(userID int) ([]models.Task, error)
	GetByID(id int) (models.Task, error)

	Edit(task models.Task, userID int) error

	SearchTask(taskName string) ([]models.Task, error)

	MarkValueSet(taskID int, status int) error

	DeleteByTaskID(id int, userID int) error
}

type taskService struct {
	rp task.TaskRepo
}

func NewTaskService(rp task.TaskRepo) Service {
	return taskService{rp: rp}
}

func (s taskService) Create(task models.Task) error {
	if err := s.rp.Create(task); err != nil {
		return fmt.Errorf("create err: %w", err)
	}
	return nil
}

func (s taskService) GetByUserID(id int) (models.Task, error) {
	byID, err := s.rp.GetByUserID(id)
	if err != nil {
		return models.Task{}, fmt.Errorf("users by id err: %w", err)
	}

	return byID, nil
}

func (s taskService) Edit(task models.Task, userID int) error {
	err := s.rp.Update(task, userID)
	if err != nil {
		return fmt.Errorf("update err: %w", err)
	}
	return nil
}

func (s taskService) SearchTask(taskName string) ([]models.Task, error) {
	tasks, err := s.rp.SearchTask(taskName)
	if err != nil {
		return tasks, fmt.Errorf("search task: %w", err)
	}
	return tasks, nil
}

func (s taskService) MarkValueSet(taskID int, status int) error {
	err := s.rp.MarkValueSet(taskID, status)
	if err != nil {
		return fmt.Errorf("mark value set: %w", err)
	}
	return nil
}

func (s taskService) DeleteByTaskID(id int, userID int) error {
	err := s.rp.DeleteByTaskID(id, userID)
	if err != nil {
		return fmt.Errorf("delete by task id err: %w", err)
	}
	return nil
}

func (s taskService) GetTasksByUserID(userID int) ([]models.Task, error) {
	tasks, err := s.rp.GetTasksByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get task by users id err: %w", err)
	}
	return tasks, err
}

func (s taskService) GetByID(id int) (models.Task, error) {
	byID, err := s.rp.GetByID(id)
	if err != nil {
		return models.Task{}, fmt.Errorf("get by id err: %w", err)
	}

	return byID, nil
}
