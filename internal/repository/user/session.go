package user

import "ToDoWithKolya/internal/models"

func (r userRepo) CreateSession(userID int, session string) error {
	_, err := r.db.Exec("insert into sessions(user_id, session) values (?, ?)", userID, session)
	return err
}

func (r userRepo) UpdateSession(userID int, session string) error {
	_, err := r.db.Exec("update sessions set session = ? where user_id = ?", session, userID)
	return err
}

func (r userRepo) CheckUserInSessions(userID int) (bool, error) {
	var result int
	row := r.db.QueryRow("select exists( select 1 from sessions where user_id = ?)", userID)
	if row.Err() != nil {
		return false, row.Err()
	}

	row.Scan(&result)

	if result == 0 {
		return false, row.Err()
	}

	return true, row.Err()
}

func (r userRepo) GetUserBySession(session string) (models.User, error) {
	row := r.db.QueryRow("select u.id, u.login, u.password, u.email from sessions join users u on sessions.user_id = u.id where session = ?", session)
	if row.Err() != nil {
		return models.User{}, row.Err()
	}

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, err
}
