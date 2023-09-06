package sessions

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository/sessions"
	"context"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	GetUserID(ctx context.Context, session string) (string, error)
	CreateOrUpdate(ctx context.Context, userID, session string) error
	GetSessionInfo(ctx context.Context, session string) (models.SessionInfo, error)
	Logout(ctx context.Context, session string) error
}

type sessionService struct {
	repo sessions.Repo
}

func New(s sessions.Repo) Service {
	return sessionService{repo: s}
}

func (s sessionService) GetUserID(ctx context.Context, session string) (string, error) {
	const op = "sessionService.GetUserID"

	id, err := s.repo.GetUserID(ctx, session)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return id, nil
}
func (s sessionService) GetSessionInfo(ctx context.Context, session string) (models.SessionInfo, error) {
	const op = "sessionService.GetSessionInfo"

	info, err := s.repo.SessionInfo(ctx, session)
	if err != nil {
		return models.SessionInfo{}, errors.Wrap(err, op)
	}

	return info, nil
}

func (s sessionService) CreateOrUpdate(ctx context.Context, userID, session string) error {
	const op = "sessionService.CreateOrUpdate"

	err := s.repo.Upsert(ctx, userID, session, time.Now())
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s sessionService) Logout(ctx context.Context, session string) error {
	const op = "sessionService.Logout"

	err := s.repo.Delete(ctx, session)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
