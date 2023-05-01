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
	case !re.MatchString(u.Email):
		errs = append(errs, "wrong email")
	case utf8.RuneCountInString(u.Email) > 27 || utf8.RuneCountInString(u.Email) < 7:
		errs = append(errs, "wrong email len")

	case utf8.RuneCountInString(u.Login) > 15 || utf8.RuneCountInString(u.Login) < 4:
		errs = append(errs, "wrong login len")

	case utf8.RuneCountInString(u.Password) > 55 || utf8.RuneCountInString(u.Password) < 4:
		errs = append(errs, "wrong password len")

	}

	return errs
}

func (r LoginRequest) Validate() []string {
	errs := make([]string, 0, 2)

	switch {
	case utf8.RuneCountInString(r.Login) > 15 || utf8.RuneCountInString(r.Login) < 4:
		errs = append(errs, "wrong login len")
	case utf8.RuneCountInString(r.Password) > 55 || utf8.RuneCountInString(r.Password) < 4:
		errs = append(errs, "wrong password len")
	}

	return errs
}
