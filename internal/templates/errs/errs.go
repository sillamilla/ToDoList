package errs

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Errors struct {
	errorPage *template.Template
	Status    int
	Err       string
	Text      string
}

type Validator interface {
	Validate() []string
}

func ErrorWrap(w http.ResponseWriter, validErr error, status int) {
	errorPage, tmplErr := template.ParseFiles("./internal/templates/errs/errors.html")
	if tmplErr != nil {
		panic(tmplErr)
	}

	errs := Errors{
		Status: status,
		Err:    fmt.Sprintf("%s", validErr),
		Text:   validErr.Error(),
	}

	err := errorPage.Execute(w, errs)
	if err != nil {
		log.Println(err)
	}

}

func Validate(w http.ResponseWriter, v Validator) bool {
	if errs := v.Validate(); len(errs) != 0 {
		var validatorTextErr string
		for _, text := range errs {
			validatorTextErr += fmt.Sprintf("error: %s \n", text)
		}

		validatorErr := fmt.Errorf(validatorTextErr)
		ErrorWrap(w, validatorErr, http.StatusBadRequest)
		return false
	}
	return true
}
