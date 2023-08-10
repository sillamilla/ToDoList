package users

import (
	"ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Users interface {
	Create(ctx context.Context, user models.User) error
	GetByCredentials(ctx context.Context, input models.Input) (models.User, error)
	IsUsernameExist(ctx context.Context, userName string) bool //todo check
	GetByUsername(ctx context.Context, username string) (models.User, error)

	CreateSession(ctx context.Context, userID, session string) error
	GetUserID(ctx context.Context, session string) (string, error)
	UpsertSession(ctx context.Context, userID, session string) error
	GetSessionTime(ctx context.Context, session string) (time.Time, error)
	Delete(ctx context.Context, session string) error
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

func (r users) GetByCredentials(ctx context.Context, input models.Input) (models.User, error) {
	var user models.User

	filter := bson.M{"username": input.Login, "password": input.Password}
	err := r.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
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

func (r users) IsUsernameExist(ctx context.Context, userName string) bool {
	filter := bson.M{"login": userName}

	count, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return true
	}

	return count > 0
}

func (r users) CreateSession(ctx context.Context, userID, session string) error {
	sessionDoc := bson.M{"user_id": userID, "session": session, "created_at": time.Now()}

	_, err := r.db.InsertOne(ctx, sessionDoc)
	if err != nil {
		return err
	}

	return nil
}

func (r users) UpsertSession(ctx context.Context, userID, session string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"session": session}}
	opts := options.Update().SetUpsert(true)

	_, err := r.db.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r users) GetUserID(ctx context.Context, session string) (string, error) {
	var id string

	err := r.db.FindOne(ctx, bson.M{"session": session}).Decode(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r users) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r users) GetSessionTime(ctx context.Context, session string) (time.Time, error) {
	var last time.Time

	err := r.db.FindOne(ctx, bson.M{"session": session}).Decode(&last)
	if err != nil {
		return time.Time{}, err
	}

	return last, nil
}
