package helper

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, errMsg string) {
	w.WriteHeader(status)
	io.WriteString(w, errMsg)
}

func SendJson(w http.ResponseWriter, data interface{}, status int) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
		// todo: send error here and delete the return
	}

	w.WriteHeader(status)

	_, err = w.Write(marshal)
	if err != nil {
		return err
	}

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

func GenerateSession() (string, error) {
	buf := make([]byte, 32)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
