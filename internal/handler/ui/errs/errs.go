package errs

import (
	"ToDoWithKolya/internal/handler/api/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Stuff struct {
	Status   int
	Err      string
	Redirect string
}

const (
	signUp = "/sign-up"
	signIn = "/sign-in"
	home   = "/"
)

func ErrorWrap(w http.ResponseWriter, validErr error, status int) {
	errorPage, tmplErr := template.ParseFiles("./internal/templates/errs/errors.html")
	if tmplErr != nil {
		panic(tmplErr)
	}

	valErr := fmt.Sprintf("%s", validErr)

	var redirect string
	switch valErr {
	case "wrong login len":
		redirect = signUp
	case "wrong password len":
		redirect = signUp
	case "wrong email":
		redirect = signUp
	case "user with this username already exist":
		redirect = signUp
	default:
		redirect = home
	}

	stuff := Stuff{
		Status:   status,
		Err:      valErr,
		Redirect: redirect,
	}

	err := errorPage.Execute(w, stuff)
	if err != nil {
		log.Println(err)
	}

}

func Validate(v helper.Validator) string {
	if errs := v.Validate(); len(errs) != 0 {
		var validatorTextErr string
		for _, text := range errs {
			validatorTextErr += fmt.Sprintf("%s", text)
		}

		return validatorTextErr
	}
	return ""
}
