package users

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/users"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Create(ctx context.Context, user models.User) error
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
}

type userService struct {
	repo users.Repo
}

func New(u users.Repo) Service {
	return userService{repo: u}
}

func (s userService) Create(ctx context.Context, user models.User) error {
	const op = "userService.Create"

	if err := s.repo.Create(ctx, user); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s userService) GetByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "userService.GetByUsername"

	byUsername, err := s.repo.GetByUsername(ctx, username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, errors.Wrap(err, "User not found")
	} else if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	return byUsername, nil
}

func (s userService) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	const op = "userService.IsUsernameAvailable"

	isUsernameExist, err := s.repo.IsUsernameExist(ctx, username)
	if err != nil {
		return false, errors.Wrap(err, op)
	}

	return isUsernameExist, nil
}
