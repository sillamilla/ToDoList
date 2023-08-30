package errs

import (
	"ToDoWithKolya/internal/handler/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Errors struct {
	Status   int
	Err      string
	Redirect string
}

func HandleError(w http.ResponseWriter, newErr error, status int) {
	errorPage, tmplErr := template.ParseFiles("./internal/templates/errs/errors.html")
	if tmplErr != nil {
		panic(tmplErr)
	}

	stuff := Errors{
		Status: status,
		Err:    newErr.Error(),
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
