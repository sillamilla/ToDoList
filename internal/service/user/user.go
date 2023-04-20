package user

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/user"
	"fmt"
)

type Service interface {
	Register(user models.User) error
	Login(req models.LoginRequest) (string, error)
	Logout(session string) error
	GetUserBySession(session string) (models.User, error)
}

type userService struct {
	rp user.UserRepo
}

func NewUserService(rp user.UserRepo) Service {
	return userService{rp: rp}
}

func (s userService) Register(user models.User) error {
	if err := s.rp.Create(user); err != nil {
		return fmt.Errorf("register: %w", err)
	}
	return nil
}

func (s userService) Login(req models.LoginRequest) (string, error) {
	user, err := s.rp.GetByLogin(req.Login, req.Password)
	if err != nil {
		return "Сторінку не знайдено, не вірний логін або пароль", err
	}

	session, err := GenerateSession()
	if err != nil {
		return "", err
	}

	err = s.rp.UpdateSession(user.ID, session)
	if err != nil {
		return "", err
	}

	return session, err
}

func (s userService) Logout(session string) error {
	err := s.rp.DeleteSession(session)
	if err != nil {
		return err
	}

	return nil
}

func (s userService) GetUserBySession(session string) (models.User, error) {
	bySession, err := s.rp.GetUserBySession(session)
	if err != nil {
		return models.User{}, err
	}

	return bySession, err
}
