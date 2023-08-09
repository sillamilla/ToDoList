package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Login     string    `json:"login,omitempty"`
	Password  string    `json:"password,omitempty"`
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

type Input struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Session string `json:"session,omitempty"`
}
