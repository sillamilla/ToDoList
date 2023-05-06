package users

import (
	"ToDoWithKolya/internal/models"
	"database/sql"
	"time"
)

type UserRepo interface {
	Create(user models.User) error

	GetByLogin(login, password string) (models.User, error)
	GetUserBySession(session string) (models.User, error)
	GetUsernames(userName string) bool

	CreateSession(userID int, session string) error
	UpsertSession(userID int, session string) error
	GetSessionLastActive(session string) (time.Time, error)
	DeleteSession(session string) error
}

type userRepo struct {
	db *sql.DB
}

func Repo(db *sql.DB) UserRepo {
	return userRepo{db: db}
}

func (r userRepo) Create(user models.User) error {
	_, err := r.db.Exec("insert into users(login, password, email) values (?, ?, ?)", user.Login, user.Password, user.Email)

	return models.DBErr(err)
}

func (r userRepo) GetByLogin(login, password string) (models.User, error) {
	row := r.db.QueryRow("select * from users where login = ? and password = ?", login, password)
	if models.DBErr(row.Err()) != nil {
		return models.User{}, models.DBErr(row.Err())
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if models.DBErr(err) != nil {
		return models.User{}, models.DBErr(err)
	}

	return user, models.DBErr(err)
}

func (r userRepo) GetUsernames(userName string) bool {
	row := r.db.QueryRow("select * from users where login = ?", userName)
	if models.DBErr(row.Err()) != nil {
		return false
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if models.DBErr(err) != nil {
		return true
	}

	if user.Login == userName {
		return false
	}

	return true
}
