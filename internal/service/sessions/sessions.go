package sessions

import (
	"ToDoWithKolya/internal/repository/sessions"
	"context"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	GetUserID(ctx context.Context, session string) (string, error)
	CreateOrUpdate(ctx context.Context, userID, session string) error
	LastActiveExpired(ctx context.Context, session string) (bool, error)
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
func (s sessionService) LastActiveExpired(ctx context.Context, session string) (bool, error) {
	const op = "sessionService.LastActiveExpired"

	lastActive, err := s.repo.GetSessionTime(ctx, session)
	if err != nil {
		return true, errors.Wrap(err, op)
	}

	sessionExpireTime := lastActive.Add(170 * time.Hour)
	if sessionExpireTime.Before(time.Now()) {
		return true, nil
	}

	return false, nil
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
