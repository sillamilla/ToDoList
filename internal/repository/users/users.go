package users

import (
	"ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user models.User) error
	GetByCredentials(ctx context.Context, login, password string) (models.User, error)

	IsUsernameExist(ctx context.Context, userName string) bool //todo check
}

type users struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Users {
	return users{
		db: database.Collection("users"),
	}
}

func (r users) Create(ctx context.Context, user models.User) error {
	filter := bson.M{
		"id":         user.ID,
		"login":      user.Login,
		"password":   user.Password,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	_, err := r.db.InsertOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r users) GetByCredentials(ctx context.Context, login, password string) (models.User, error) {
	var user models.User
	filter := bson.M{"login": login, "password": password}

	err := r.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r users) IsUsernameExist(ctx context.Context, userName string) bool {
	filter := bson.M{"login": userName}

	count, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return true
	}

	return count > 0
}
