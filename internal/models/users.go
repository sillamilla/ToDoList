package models

import (
	"regexp"
	"unicode/utf8"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type LoginRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Session string `json:"session,omitempty"`
}

func (u User) Validate() []string {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	errs := make([]string, 0, 2)

	switch {
	case re.MatchString(u.Email) == false:
		errs = append(errs, "wrong email")

	case utf8.RuneCountInString(u.Email) > 60 || utf8.RuneCountInString(u.Email) < 4:
		errs = append(errs, "wrong description len")

	case utf8.RuneCountInString(u.Login) > 15 || utf8.RuneCountInString(u.Email) < 12:
		errs = append(errs, "wrong title len")
	}

	return errs
}
