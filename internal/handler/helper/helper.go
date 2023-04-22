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
	"strconv"
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

func GetIntFromURL(r *http.Request, key string) (int, error) {
	vars := mux.Vars(r)
	valueStr := vars[key]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	return value, err
}

func SendJson(w http.ResponseWriter, data interface{}, status int) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
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
		return nil, err
	}
	if errs := v.Validate(); len(errs) > 0 {
		return errs, nil
	}

	return nil, nil
}
