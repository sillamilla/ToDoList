package users

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/users"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	SignUp(ctx context.Context, user models.User) (string, error)
	SignIn(ctx context.Context, req models.Input) (string, error)
	UserCheckExist(ctx context.Context, userName string) bool
}

type userService struct {
	user users.Users
}

func New(u users.Users) Service {
	return userService{user: u}
}

func (s userService) SignUp(ctx context.Context, user models.User) (string, error) {
	Password := user.Password
	newUser := models.User{
		ID:        uuid.New().String(),
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}

	user.Password = helper.HashGenerate(user.Password)

	if err := s.user.Create(ctx, newUser); err != nil {
		return "", fmt.Errorf("register err: %w", err)
	}

	session, err := s.SignIn(ctx, models.Input{
		Login:    user.Login,
		Password: Password,
	})
	if err != nil {
		return "", fmt.Errorf("registe:login: error, err: %w", err)
	}

	return session, nil
}

func (s userService) SignIn(ctx context.Context, input models.Input) (string, error) {
	//user, err := s.GetByUsername(ctx, input.Username)
	//if errors.Is(err, mongo.ErrNoDocuments) {
	//	return model.User{}, errors.New("User not found")
	//} else if err != nil {
	//	return model.User{}, errors.Wrap(err, "service.SignIn.GetByUsername")
	//}
	//
	//err = helper.ComparePassword(user.Password, input.Password)
	//if err != nil {
	//	return model.User{}, errors.New("Invalid password")
	//}
	//
	//input.Password = user.Password
	//
	//signUser, err := s.mo.SignIn(ctx, input)
	//if err != nil {
	//	return model.User{}, errors.Wrap(err, "service.SignIn")
	//}
	//
	//err = s.re.UpsertSession(ctx, signUser.ID, signUser.Session)
	//if err != nil {
	//	return model.User{}, errors.Wrap(err, "service.SignIn.UpsertSession")
	//}
	//
	//context.WithValue(ctx, "user", signUser)
	//
	//return signUser, nil

	return "", nil
}

func (s userService) UserCheckExist(ctx context.Context, userName string) bool {
	return s.user.IsUsernameExist(ctx, userName)
}
