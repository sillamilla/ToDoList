package user

import (
	"ToDoWithKolya/internal/models"
	"database/sql"
	"time"
)

type UserRepo interface {
	Create(user models.User) error
	Update(user models.User) error

	GetByLogin(login, password string) (models.User, error)
	GetUserBySession(session string) (models.User, error)

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

	return err
}

func (r userRepo) Update(user models.User) error {
	_, err := r.db.Exec(
		"update users set login = ?, password = ?, email = ? where id = ?",
		user.Login,
		user.Password,
		user.Email,
		user.ID,
	)
	return err
}

func (r userRepo) GetByLogin(login, password string) (models.User, error) {
	row := r.db.QueryRow("select * from users where login = ? and password = ?", login, password)

	if row.Err() != nil {
		return models.User{}, row.Err()
	}
	var user models.User

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	return user, err
}