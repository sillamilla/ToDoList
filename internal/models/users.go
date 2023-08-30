package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	Session   string    `json:"session,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UserAndTask struct {
	User     User      `json:"user,omitempty"`
	Tasks    []Task    `json:"tasks,omitempty"`
	DayTasks []DayTask `json:"day_tasks,omitempty"`
}

type SearchAndStatus struct {
	Search string `json:"search,omitempty"`
	Status string `json:"status,omitempty"`
}

type UserInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type LoginResponse struct {
	Session string `json:"session,omitempty"`
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
