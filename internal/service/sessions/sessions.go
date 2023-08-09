package sessions

import (
	"ToDoWithKolya/internal/repository/sessions"
	"context"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	NewSession(ctx context.Context, userID, session string) error
	GetUserID(ctx context.Context, session string) (string, error)
	Upsert(ctx context.Context, userID, session string) error
	LastActive(ctx context.Context, session string) (time.Time, error)
	Logout(ctx context.Context, session string) error
}

type sessionService struct {
	session sessions.Sessions
}

func New(s sessions.Sessions) Service {
	return sessionService{session: s}
}

func (s sessionService) NewSession(ctx context.Context, userID, session string) error {
	err := s.session.Create(ctx, userID, session)
	if err != nil {
		return errors.Wrap(err, "NewSession_Create err")
	}

	return nil
}

func (s sessionService) GetUserID(ctx context.Context, session string) (string, error) {
	id, err := s.session.GetUserID(ctx, session)
	if err != nil {
		return "", errors.Wrap(err, "GetUserID_GetUserID")
	}

	return id, nil
}
func (s sessionService) LastActive(ctx context.Context, session string) (time.Time, error) {
	sessionTime, err := s.session.GetSessionTime(ctx, session)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "LastActive_GetSessionTime")
	}

	return sessionTime, nil
}

func (s sessionService) Upsert(ctx context.Context, userID, session string) error {
	err := s.session.Upsert(ctx, userID, session)
	if err != nil {
		return errors.Wrap(err, "Upsert_Upsert")
	}

	return nil
}

func (s sessionService) Logout(ctx context.Context, session string) error {
	err := s.session.Delete(ctx, session)
	if err != nil {
		return errors.Wrap(err, "Logout_Delete")
	}

	return nil
}
