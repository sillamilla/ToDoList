package users

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/users"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service interface {
	SignUp(ctx context.Context, user models.Input) (models.User, error)
	SignIn(ctx context.Context, req models.Input) (models.User, error)
	UserCheckExist(ctx context.Context, userName string) bool
	GetByUsername(ctx context.Context, username string) (models.User, error)

	NewSession(ctx context.Context, userID, session string) error
	GetUserID(ctx context.Context, session string) (string, error)
	Upsert(ctx context.Context, userID, session string) error
	LastActiveExpired(ctx context.Context, session string) (bool, error)
	Logout(ctx context.Context, session string) error
}

type userService struct {
	user users.Users
}

func New(u users.Users) Service {
	return userService{user: u}
}

func (s userService) SignUp(ctx context.Context, input models.Input) (models.User, error) {
	ok := s.user.IsUsernameExist(ctx, input.Login)
	if ok {
		return models.User{}, errors.Wrap(nil, "username taken")
	}

	password, err := helper.HashPassword(input.Password)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp.HashPassword")
	}
	input.Password = password

	id := uuid.NewString()
	session, err := helper.GenerateSession()
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp.GenerateSession")
	}

	err = s.Upsert(ctx, id, session)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp.UpsertSessions") //todo to many wrap
	}

	newUser := models.UserFromInput(id, input, session, time.Now())
	err = s.user.Create(ctx, newUser)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp")
	}

	return newUser, nil
}

func (s userService) SignIn(ctx context.Context, input models.Input) (models.User, error) {
	user, err := s.user.GetByCredentials(ctx, input)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, errors.New("User not found")
	} else if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignIn.GetByUsername")
	}

	err = helper.ComparePassword(user.Password, input.Password)
	if err != nil {
		return models.User{}, errors.New("Invalid password")
	}

	input.Password = user.Password

	signUser, err := s.user.GetByCredentials(ctx, input)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignIn")
	}

	err = s.user.UpsertSession(ctx, signUser.ID, signUser.Session)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignIn.UpsertSession")
	}

	return signUser, nil
}

func (s userService) GetByUsername(ctx context.Context, username string) (models.User, error) {
	byUsername, err := s.user.GetByUsername(ctx, username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, errors.Wrap(err, "User not found")
	} else if err != nil {
		return models.User{}, errors.Wrap(err, "service.searchByUsername")
	}

	return byUsername, nil
}

func (s userService) UserCheckExist(ctx context.Context, userName string) bool {
	return s.user.IsUsernameExist(ctx, userName)
}

func (s userService) NewSession(ctx context.Context, userID, session string) error {
	err := s.user.CreateSession(ctx, userID, session)
	if err != nil {
		return errors.Wrap(err, "NewSession_Create err")
	}

	return nil
}

func (s userService) GetUserID(ctx context.Context, session string) (string, error) {
	id, err := s.user.GetUserID(ctx, session)
	if err != nil {
		return "", errors.Wrap(err, "GetUserID_GetUserID")
	}

	return id, nil
}
func (s userService) LastActiveExpired(ctx context.Context, session string) (bool, error) {
	lastActive, err := s.user.GetSessionTime(ctx, session)
	if err != nil {
		return true, errors.Wrap(err, "LastActiveExpired_GetSessionTime")
	}

	sessionExpireTime := lastActive.Add(170 * time.Hour)
	if sessionExpireTime.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

func (s userService) Upsert(ctx context.Context, userID, session string) error {
	err := s.user.UpsertSession(ctx, userID, session)
	if err != nil {
		return errors.Wrap(err, "Upsert_Upsert")
	}

	return nil
}

func (s userService) Logout(ctx context.Context, session string) error {
	err := s.user.Delete(ctx, session)
	if err != nil {
		return errors.Wrap(err, "Logout_Delete")
	}

	return nil
}
