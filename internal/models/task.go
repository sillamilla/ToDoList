package models

import (
	"time"
)

type Task struct {
	ID          string    `bson:"id,omitempty"`
	UserID      string    `bson:"user_id,omitempty"`
	Title       string    `bson:"title,omitempty"`
	Description string    `bson:"description,omitempty"`
	IsDone      int       `bson:"is_done,omitempty"`
	CreatedAt   time.Time `bson:"created_at"`
}

type DayTask struct {
	ID        string    `bson:"id,omitempty"`
	UserID    int       `bson:"user_id,omitempty"`
	Title     string    `bson:"title,omitempty"`
	IsDone    int       `bson:"is_done,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
}
