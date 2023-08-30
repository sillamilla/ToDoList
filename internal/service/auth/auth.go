package auth

import (
	"ToDoWithKolya/internal/handler/helper"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service/sessions"
	"ToDoWithKolya/internal/service/users"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service interface {
	SignUp(ctx context.Context, user models.UserInput) (models.User, error)
	SignIn(ctx context.Context, req models.UserInput) (models.User, error)
	Logout(ctx context.Context, session string) error
	LastActiveExpired(ctx context.Context, session string) (bool, error)
}

type auth struct {
	sessionSrv sessions.Service
	userSrv    users.Service
}

func New(sessionsSrv sessions.Service, userSrv users.Service) Service {
	return auth{
		sessionSrv: sessionsSrv,
		userSrv:    userSrv,
	}
}

func (a auth) SignUp(ctx context.Context, input models.UserInput) (models.User, error) {
	const op = "auth.SignUp"

	ok, err := a.userSrv.IsUsernameAvailable(ctx, input.Username)
	if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	if ok {
		return models.User{}, errors.Wrap(fmt.Errorf("username taken"), "create a new one")
	}

	password, err := helper.HashPassword(input.Password)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp.HashPassword")
	}
	input.Password = password

	id := uuid.NewString()
	newUser := models.UserFromInput(id, input, "", time.Now())

	session, err := helper.GenerateSession()
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp.GenerateSession")
	}

	err = a.sessionSrv.CreateOrUpdate(ctx, id, session)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp.UpsertSessions") //todo to many wrap
	}

	err = a.userSrv.Create(ctx, newUser)
	if err != nil {
		return models.User{}, errors.Wrap(err, "service.SignUp")
	}

	return newUser, nil
}

func (a auth) SignIn(ctx context.Context, input models.UserInput) (models.User, error) {
	const op = "auth.SignIn"

	user, err := a.userSrv.GetByUsername(ctx, input.Username)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return models.User{}, errors.New("User not found")
	} else if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	err = helper.ComparePassword(user.Password, input.Password)
	if err != nil {
		return models.User{}, errors.New("Invalid password")
	}

	if user.Session == "" {
		session, err := helper.GenerateSession()
		if err != nil {
			return models.User{}, errors.Wrap(err, op)
		}
		user.Session = session
	}

	input.Password = user.Password //todo nahuya

	err = a.sessionSrv.CreateOrUpdate(ctx, user.ID, user.Session)
	if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	return user, nil
}

func (a auth) LastActiveExpired(ctx context.Context, session string) (bool, error) {
	const op = "auth.LastActiveExpired"

	lastActive, err := a.sessionSrv.LastActiveExpired(ctx, session)
	if err != nil {
		return true, errors.Wrap(err, op)
	}

	return lastActive, nil
}

func (a auth) Logout(ctx context.Context, session string) error {
	const op = "auth.Logout"

	err := a.sessionSrv.Logout(ctx, session)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
