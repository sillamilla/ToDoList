package tasks

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/tasks"
	"context"
	"github.com/pkg/errors"
)

type Service interface {
	NewTask(ctx context.Context, task models.Task) error
	GetTasks(ctx context.Context, userID string) ([]models.Task, error)
	GetByID(ctx context.Context, userID, id string) (models.Task, error)
	Edit(ctx context.Context, task models.Task) error
	SearchTasks(ctx context.Context, taskName, UserID string) ([]models.Task, error)
	MarkValueSet(ctx context.Context, taskID string, status int) error
	Delete(ctx context.Context, userID, id string) error
}

type taskService struct {
	task tasks.Repo
}

func New(t tasks.Repo) Service {
	return taskService{task: t}
}

func (s taskService) NewTask(ctx context.Context, task models.Task) error {
	const op = "taskService.NewTask"

	if err := s.task.Create(ctx, task); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s taskService) Edit(ctx context.Context, task models.Task) error {
	const op = "taskService.Edit"

	err := s.task.Update(ctx, task)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s taskService) SearchTasks(ctx context.Context, taskName, userID string) ([]models.Task, error) {
	const op = "taskService.SearchTasks"
	values, err := s.task.Search(ctx, taskName, userID)
	if err != nil {
		return values, errors.Wrap(err, op)
	}

	return values, nil
}

func (s taskService) MarkValueSet(ctx context.Context, taskID string, status int) error {
	const op = "taskService.MarkValueSet"

	err := s.task.MarkValueSet(ctx, taskID, status)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s taskService) Delete(ctx context.Context, userID, id string) error {
	const op = "taskService.Delete"

	err := s.task.Delete(ctx, userID, id)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s taskService) GetTasks(ctx context.Context, userID string) ([]models.Task, error) {
	const op = "taskService.GetTasks"

	all, err := s.task.GetAll(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return all, err
}

func (s taskService) GetByID(ctx context.Context, userID, id string) (models.Task, error) {
	const op = "taskService.GetByID"

	task, err := s.task.Get(ctx, userID, id)
	if err != nil {
		return models.Task{}, errors.Wrap(err, op)
	}

	return task, nil
}
