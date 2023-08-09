package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	IsDone      int       `json:"is_done,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type DayTask struct {
	ID        string    `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	IsDone    int       `json:"is_done,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
