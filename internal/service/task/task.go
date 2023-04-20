package task

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/task"
	"fmt"
)

type Service interface {
	Create(task models.Task) error
	MarkAsDone(id int) error

	GetByUserID(id int) (models.Task, error)
	GetTasksByUserID(userID int) ([]models.Task, error)
	GetByID(id int) (models.Task, error)

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
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (s taskService) GetByUserID(id int) (models.Task, error) {
	byID, err := s.rp.GetByUserID(id)
	if err != nil {
		return models.Task{}, err
	}

	return byID, nil
}

func (s taskService) MarkAsDone(id int) error {
	byID, err := s.GetByUserID(id)
	if err != nil {
		return err
	}

	byID.Done = true
	err = s.rp.Update(byID)
	return err
}

func (s taskService) DeleteByTaskID(id int, userID int) error {
	err := s.rp.DeleteByTaskID(id, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s taskService) GetTasksByUserID(userID int) ([]models.Task, error) {
	tasks, err := s.rp.GetTasksByUserID(userID)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (s taskService) GetByID(id int) (models.Task, error) {
	byID, err := s.rp.GetByID(id)
	if err != nil {
		return models.Task{}, err
	}

	return byID, nil
}
