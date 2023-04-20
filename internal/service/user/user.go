package user

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/user"
	"fmt"
	"log"
	"time"
)

type Service interface {
	Register(user models.User) error
	Login(req models.LoginRequest) (string, error)
	Logout(session string) error

	GetUserBySession(session string) (models.User, error)

	TimeOutSession(session string) error
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
		return fmt.Errorf("register: %w", err)
	}
	return nil
}

func (s userService) TimeOutSession(session string) error {
	time.Sleep(10 * time.Second)

	lastActive, err := s.rp.GetSessionLastActive(session)
	if err != nil {
		log.Println("Не удалось получить время последней активности сессии:", err)
		return err
	}

	if time.Since(lastActive) >= 30*time.Minute {
		err := s.Logout(session)
		if err != nil {
			log.Println("Не удалось выполнить выход из сессии:", err)
		}
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

	go s.TimeOutSession(session)
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

	return bySession, nil
}

func (s userService) GetSessionLastActive(session string) (time.Time, error) {
	sessionTime, err := s.rp.GetSessionLastActive(session)
	if err != nil {
		return time.Time{}, err
	}
	return sessionTime, nil
}
