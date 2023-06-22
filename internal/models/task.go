package models

import (
	"time"
)

type Task struct {
	ID          int       `json:"id,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	IsDone      int       `json:"is_done,omitempty"`
	CreatedDate time.Time `json:"created_date"`
}
