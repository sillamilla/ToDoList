package helper

import (
	"ToDoWithKolya/internal/models"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, err error) {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		status = http.StatusNotFound
	case errors.Is(err, mongo.ErrNoDocuments): //todo not correct cases
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

func GenerateSession() (string, error) {
	buf := make([]byte, 32)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func HashGenerate(password string) string {
	hashObject := sha256.New()
	hashObject.Write([]byte(password))
	hashedBytes := hashObject.Sum(nil)

	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}

func HashVerify(password, hash string) bool {
	hashedString := HashGenerate(password)
	return hashedString == hash
}

func UserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value("users").(models.User)
	return user, ok
}
