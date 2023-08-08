package users

import (
	"ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepo interface {
	Create(user models.User) error

	GetByLogin(login, password string) (models.User, error)
	GetUserBySession(session string) (models.User, error)
	GetUsernames(userName string) bool

	CreateSession(userID string, session string) error
	UpsertSession(userID string, session string) error
	GetSessionLastActive(session string) (time.Time, error)
	DeleteSession(session string) error
}

type userRepo struct {
	db *mongo.Collection
}

func NewUserRepo(database mongo.Database) userRepo {
	return userRepo{
		db: database.Collection("users"),
	}
}

func (r userRepo) Create(user models.User) error {
	userDoc := bson.M{"login": user.Login, "password": user.Password, "email": user.Email}
	_, err := r.db.InsertOne(context.TODO(), userDoc)
	return err
}

func (r userRepo) GetByLogin(login, password string) (models.User, error) {
	var user models.User
	err := r.db.FindOne(context.TODO(), bson.M{"login": login, "password": password}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, models.DBErr(models.ErrNotFound)
		}
		return models.User{}, err
	}
	return user, nil
}

func (r userRepo) GetUsernames(userName string) bool {
	filter := bson.M{"login": userName}
	count, err := r.db.CountDocuments(context.TODO(), filter)
	if err != nil {
		return true
	}
	return count > 0
}
