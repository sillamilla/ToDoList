package users

import (
	"ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repo interface {
	Create(ctx context.Context, user models.User) error
	IsUsernameExist(ctx context.Context, userName string) (bool, error) //todo check
	GetByUsername(ctx context.Context, username string) (models.User, error)
}

type users struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Repo {
	return users{
		db: database.Collection("users"),
	}
}

func (r users) Create(ctx context.Context, user models.User) error {
	filter := bson.M{
		"id":         user.ID,
		"username":   user.Username,
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

func (r users) GetByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	filter := bson.M{"username": username}
	err := r.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r users) IsUsernameExist(ctx context.Context, userName string) (bool, error) {
	filter := bson.M{"username": userName}

	count, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
