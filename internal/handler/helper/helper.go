package helper

import (
	"ToDoWithKolya/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func SendError(w http.ResponseWriter, status int, err error) {
	if status != http.StatusInternalServerError {
		fmt.Fprintf(w, "error: %s", err)
	}

	w.WriteHeader(status)
	if err != nil {
		log.Println(err)
	}
}

func UserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value("user").(models.User)
	return user, ok
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
