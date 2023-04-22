package repository

import (
	"ToDoWithKolya/internal/models"
	"database/sql"
	"errors"
)

func Err(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNotFound
	}

	return err
}
