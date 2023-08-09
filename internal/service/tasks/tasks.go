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
	GetByID(ctx context.Context, id string) (models.Task, error)
	Edit(ctx context.Context, task models.Task, userID string) error
	SearchTasks(ctx context.Context, taskName, UserID string) ([]models.Task, error)
	MarkValueSet(ctx context.Context, taskID string, status int) error
	Delete(ctx context.Context, id string, userID string) error
}

type taskService struct {
	task tasks.Tasks
}

func New(t tasks.Tasks) Service {
	return taskService{task: t}
}

func (s taskService) NewTask(ctx context.Context, task models.Task) error {
	if err := s.task.Create(ctx, task); err != nil {
		return errors.Wrap(err, "NewTask_Create err")
	}

	return nil
}

func (s taskService) Edit(ctx context.Context, task models.Task, userID string) error {
	err := s.task.Update(ctx, task, userID)
	if err != nil {
		return errors.Wrap(err, "Edit_Update err")
	}

	return nil
}

func (s taskService) SearchTasks(ctx context.Context, taskName, userID string) ([]models.Task, error) {
	tasks, err := s.task.Search(ctx, taskName, userID)
	if err != nil {
		return tasks, errors.Wrap(err, "GetTasks_Search err")
	}

	return tasks, nil
}

func (s taskService) MarkValueSet(ctx context.Context, taskID string, status int) error {
	err := s.task.MarkValueSet(ctx, taskID, status)
	if err != nil {
		return errors.Wrap(err, "MarkValueSet_MarkValueSet err")
	}

	return nil
}

func (s taskService) Delete(ctx context.Context, id string, userID string) error {
	err := s.task.Delete(ctx, id, userID)
	if err != nil {
		return errors.Wrap(err, "DeleteByTaskID_Delete err")
	}

	return nil
}

func (s taskService) GetTasks(ctx context.Context, userID string) ([]models.Task, error) {
	tasks, err := s.task.GetAll(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "GetTasksByUserID_GetTasks err")
	}

	return tasks, err
}

func (s taskService) GetByID(ctx context.Context, id string) (models.Task, error) {
	task, err := s.task.Get(ctx, id)
	if err != nil {
		return models.Task{}, errors.Wrap(err, "GetByID_Get err")
	}

	return task, nil
}
