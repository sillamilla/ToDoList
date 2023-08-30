package models

import (
	"regexp"
	"unicode/utf8"
)

func (u User) Validate() []string {
	errs := make([]string, 0, 2)

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	email := regexp.MustCompile(pattern)

	switch {
	case u.Email != "" && !email.MatchString(u.Email):
		errs = append(errs, "wrong email")

	case utf8.RuneCountInString(u.Username) < 4 || utf8.RuneCountInString(u.Username) > 15:
		errs = append(errs, "wrong username len")

	case utf8.RuneCountInString(u.Password) < 4 || utf8.RuneCountInString(u.Password) > 60:
		errs = append(errs, "wrong password len")

	}

	return errs
}

func (r UserInput) Validate() []string {
	errs := make([]string, 0, 2)

	switch {
	case utf8.RuneCountInString(r.Username) < 4 || utf8.RuneCountInString(r.Username) > 15:
		errs = append(errs, "wrong username len")
	case utf8.RuneCountInString(r.Password) < 4 || utf8.RuneCountInString(r.Password) > 60:
		errs = append(errs, "wrong password len")
	}

	return errs
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
