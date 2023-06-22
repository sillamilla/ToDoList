package users

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/users"
	"fmt"
	"time"
)

type Service interface {
	Register(user models.User) (string, error)
	Login(req models.LoginRequest) (string, error)
	Logout(session string) error

	UserCheckExist(userName string) bool
	GetUserBySession(session string) (models.User, error)

	GetSessionLastActive(session string) (time.Time, error)
}

type userService struct {
	rp users.UserRepo
}

func NewUserService(rp users.UserRepo) Service {
	return userService{rp: rp}
}

func (s userService) Register(user models.User) (string, error) {
	old := user.Password
	user.Password = HashGenerate(user.Password)

	if err := s.rp.Create(user); err != nil {
		return "", fmt.Errorf("register err: %w", err)
	}
	session, err := s.Login(models.LoginRequest{
		Login:    user.Login,
		Password: old,
	})
	if err != nil {
		return "", fmt.Errorf("registe:login: error, err: %w", err)
	}

	return session, nil
}

func (s userService) Login(req models.LoginRequest) (string, error) {
	user, err := s.rp.GetByLogin(req.Login, HashGenerate(req.Password))
	if err != nil {
		//if errors.Is(err, models.ErrNotFound) {
		//	return "", fmt.Errorf(" invalid login or password, err: %w", models.ErrUnauthorized)
		//}
		return "", fmt.Errorf("invalid login or password, err: %w", err)
	}

	session, err := GenerateSession()
	if err != nil {
		return "", fmt.Errorf("generate session, err: %w", err)
	}

	err = s.rp.UpsertSession(user.ID, session)
	if err != nil {
		return "", fmt.Errorf("upsert session error, err: %w", err)
	}

	return session, nil
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
		return models.User{}, fmt.Errorf("get users by session err: %w", err)
	}

	return bySession, nil
}

func (s userService) UserCheckExist(userName string) bool {
	ok := s.rp.GetUsernames(userName)

	return ok
}

func (s userService) GetSessionLastActive(session string) (time.Time, error) {
	sessionTime, err := s.rp.GetSessionLastActive(session)
	if err != nil {
		return time.Time{}, fmt.Errorf("last activie err: %w", err)
	}
	return sessionTime, nil
}
