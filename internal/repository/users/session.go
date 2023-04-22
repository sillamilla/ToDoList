package users

import (
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/repository"
	"time"
)

func (r userRepo) CreateSession(userID int, session string) error {
	_, err := r.db.Exec("insert into sessions(user_id, session) values (?, ?)", userID, session)
	return repository.Err(err)
}

func (r userRepo) UpsertSession(userID int, session string) error {
	_, err := r.db.Exec("INSERT INTO sessions (user_id, session) VALUES (?, ?) ON CONFLICT (user_id) DO UPDATE SET session = ?", userID, session, session)
	return repository.Err(err)
}

func (r userRepo) GetUserBySession(session string) (models.User, error) {
	row := r.db.QueryRow("select u.id, u.login, u.password, u.email from sessions join users u on sessions.user_id = u.id where session = ?", session)
	if row.Err() != nil {
		return models.User{}, repository.Err(row.Err())
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		return models.User{}, repository.Err(err)
	}

	return user, nil
}

func (r userRepo) DeleteSession(session string) error {
	_, err := r.db.Exec("delete from sessions where session = ?", session)
	return repository.Err(err)
}

func (r userRepo) GetSessionLastActive(session string) (time.Time, error) {
	row := r.db.QueryRow("select created_at from sessions where session = ?", session)
	if row.Err() != nil {
		return time.Time{}, repository.Err(row.Err())
	}

	var sessionTime time.Time
	err := row.Scan(&sessionTime)
	if err != nil {
		return time.Time{}, repository.Err(err)
	}

	return sessionTime, nil
}
