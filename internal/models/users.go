package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"id,omitempty"`
	Username  string    `bson:"username,omitempty"`
	Password  string    `bson:"password,omitempty"`
	Session   string    `bson:"session,omitempty"`
	Email     string    `bson:"email,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
}

type UserAndTask struct {
	User     User      `bson:"user,omitempty"`
	Tasks    []Task    `bson:"tasks,omitempty"`
	DayTasks []DayTask `bson:"day_tasks,omitempty"`
}

type SearchAndStatus struct {
	Search string `bson:"search,omitempty"`
	Status string `bson:"status,omitempty"`
}

type UserInput struct {
	Username string `bson:"username,omitempty"`
	Password string `bson:"password,omitempty"`
	Email    string `bson:"email,omitempty"`
}

type LoginResponse struct {
	Session string `bson:"session,omitempty"`
}

func UserFromInput(ID string, user UserInput, session string, createdAt time.Time) User {
	return User{
		ID:        ID,
		Username:  user.Username,
		Password:  user.Password,
		Session:   session,
		CreatedAt: createdAt,
		Email:     user.Email,
	}
}
