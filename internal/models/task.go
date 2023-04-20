package models

import (
	"time"
	"unicode/utf8"
)

type Task struct {
	ID          int       `json:"id,omitempty"`
	UserID      int       `json:"user_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Done        bool      `json:"done,omitempty"`
	CreatedDate time.Time `json:"created_date"`
}

func (t Task) Validate() []string {
	errs := make([]string, 0, 2)

	if utf8.RuneCountInString(t.Title) > 30 {
		errs = append(errs, "wrong title len")
	}

	if utf8.RuneCountInString(t.Description) > 300 {
		errs = append(errs, "wrong description len")
	}

	return errs
}
