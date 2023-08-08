package helper

import (
	"ToDoWithKolya/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, models.ErrUnauthorized):
		status = http.StatusUnauthorized
	}

	if status != http.StatusInternalServerError {
		fmt.Fprintf(w, "error: %s", err)
	}

	w.WriteHeader(status)
	if err != nil {
		log.Println(err)
	}
}

func FromURL(r *http.Request, key string) string {
	value := mux.Vars(r)[key]

	return value
}

func SendJson(w http.ResponseWriter, data interface{}, status int) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
		//todo send error here and delelete return
	}

	w.Write(marshal)
	w.WriteHeader(status)
	return nil
}

type Validator interface {
	Validate() []string
}

func UnmarshalAndValidate(r io.ReadCloser, v Validator) ([]string, error) {
	readAll, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	err = json.Unmarshal(readAll, v)
	if err != nil {
		return nil, err //badreq
	}

	if errs := v.Validate(); len(errs) > 0 {
		return errs, nil //entity
	}

	return nil, nil
}
