package user

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/user"
	"fmt"
	"time"
)

type Service interface {
	Register(user models.User) error
	Update(user models.User) error
	Login(req models.LoginRequest) (string, error)
	Logout(session string) error

	GetUserBySession(session string) (models.User, error)

	GetSessionLastActive(session string) (time.Time, error)
}

type userService struct {
	rp user.UserRepo
}

func NewUserService(rp user.UserRepo) Service {
	return userService{rp: rp}
}

func (s userService) Register(user models.User) error {
	if err := s.rp.Create(user); err != nil {
		return fmt.Errorf("register err: %w", err)
	}
	return nil
}

func (s userService) Update(user models.User) error {
	if err := s.rp.Update(user); err != nil {
		return fmt.Errorf("update err: %w", err)
	}
	return nil
}

func (s userService) Login(req models.LoginRequest) (string, error) {
	user, err := s.rp.GetByLogin(req.Login, req.Password)
	if err != nil {
		return "User not fiend, invalid login or password", err
	}

	session, err := GenerateSession()
	if err != nil {
		return "Generate session err", err
	}

	err = s.rp.UpsertSession(user.ID, session)
	if err != nil {
		return "Upsert session error", err
	}

	return session, err
}

func (s userService) Logout(session string) error {
	err := s.rp.DeleteSession(session)
	if err != nil {
		return fmt.Errorf("logout err: %w", err)
	}

	return nil
}

func (s userService) GetUserBySession(session string) (models.User, error) {
	bySession, err := s.rp.GetUserBySession(session)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by session err: %w", err)
	}

	return bySession, nil
}

func (s userService) GetSessionLastActive(session string) (time.Time, error) {
	sessionTime, err := s.rp.GetSessionLastActive(session)
	if err != nil {
		return time.Time{}, fmt.Errorf("last activie err: %w", err)
	}
	return sessionTime, nil
}
